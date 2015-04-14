package mydb

import (
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/coopernurse/gorp"
	logs "github.com/usoatov/myhttp/fl"

	_ "github.com/go-sql-driver/mysql"
)

type Opts struct {
	Stamp      string
	Opstamp    string
	Photostamp string
	Errdel     string
	Delay      string
	Transtime  string
	Transint   string
	Realtime   string
	Encrypt    string
	Timezone   string
}

type Oplog struct {
	Opcode, Adminid, DtTime, Obj1, Obj2, Obj3, Obj4 string
}

type Cmd struct {
	Id, Cmdbody string
}

var db *sql.DB

func Connect(mdb, host, usr, pwd string) bool {
	var err error
	//fmt.Println(mdb, host, usr, pwd)
	s := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", usr, pwd, host, mdb)
	db, err = sql.Open("mysql", s)
	if err != nil {
		log.Print(err)
		return false
	}
	err = db.Ping()
	return true

}

func Dev_id(sn string) string {
	var d_id string
	rows, err := db.Query("select id from device where serialnumber=?", sn)
	//	rows, err := db.Query("select attLogStamp as Stamp from device where serialnumber=?", sn)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&d_id)
		if err != nil {
			log.Println("ERROR in Scan")
			log.Print(err)
		}
	}

	return d_id
}

func Billing(sn string) bool {
	cmp_id := Comp_id(sn)
	n := time.Now()
	var st [4]int
	var i int
	for i = 0; i <= 3; i++ {
		var fq time.Duration
		fq = time.Duration(-24*i) * time.Hour
		d3 := n.Add(fq)
		//s := fmt.Sprintf("%04d", t.Year())
		//s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", d3.Year(), d3.Month(), d3.Day(), d3.Hour(), d3.Minute(), d3.Second())
		//fmt.Println(s)

		rows, err := db.Query("select status from billing_status where companyID = ? and ? BETWEEN f_time AND l_time", cmp_id, d3)

		if err != nil {
			log.Print(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&st[i])
			if err != nil {
				log.Println("ERROR in Scan")
				log.Print(err)
			}
		}
	}

	var j int
	for j = 0; j <= 3; j++ {
		if st[j] == 1 {
			return true
		}
	}
	return false

}

func Comp_id(sn string) string {
	var c_id string
	rows, err := db.Query("select companyID from device where serialnumber=?", sn)

	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&c_id)
		if err != nil {
			log.Println("ERROR in Scan")
			log.Print(err)
		}
	}

	return c_id
}

func Options(sn string) Opts {
	var opres Opts
	rows, err := db.Query("select attLogStamp, operLogStamp, photoStamp, errorDelay, delay, transTimes, transInterval, realtime, encrypt, timeZoneAdj from device where serialnumber=?", sn)
	//	rows, err := db.Query("select attLogStamp as Stamp from device where serialnumber=?", sn)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&opres.Stamp, &opres.Opstamp, &opres.Photostamp, &opres.Errdel, &opres.Delay, &opres.Transtime, &opres.Transint, &opres.Realtime, &opres.Encrypt, &opres.Timezone)
		if err != nil {
			log.Println("ERROR in Scan")
			log.Print(err)
		}
	}

	return opres
}

func Lastrequesttime(sn string) bool {
	stmt, err := db.Prepare("update device set lastRequestTime=NOW() where serialnumber=?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(sn)
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}
	return true
}

func Transfertime(id string) bool {
	stmt, err := db.Prepare("update devicecmds set commandTransferTime=NOW() where id=?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}
	return true
}

func GetLastreq(sn string) string {
	var gl string

	err := db.QueryRow("select lastRequestTime from device where serialNumber = ?", sn).Scan(&gl)
	if err != nil {
		log.Print(err)
		return ""
	}
	return gl
}

func InsertTempinout(sn, line string) bool {
	ls := strings.Split(line, "\t")
	pin := ls[0]
	dt := ls[1]
	eventcode := ls[2]
	verify := ls[3]
	s := fmt.Sprintf("%s\t%s\t%s\t%s\n", pin, dt, eventcode, verify)
	logs.Inout(sn, s)
	Update_stamp(sn, dt)
	//fmt.Println("pin=", pin, "dt=", dt, eventcode, verify)

	//Inout eskiligini tekshirish tInout inout time tLasreq last request time
	/*lastreq := GetLastreq(sn)
	const layout = "2006-01-02 15:04:05"
	tInout, _ := time.Parse(layout, dt)
	tLastreq, _ := time.Parse(layout, lastreq)
	dur := -360 * time.Hour //15 kun orqaga qaytarish
	tLastreq = tLastreq.Add(dur)

	// Bu yerda eski danniylar comtroli
	/*if tLastreq.Unix() > tInout.Unix() {
		return false
	}*/
	stmt, err := db.Prepare("INSERT INTO `temp_inout` (deviceSN, pin, time, status, verify) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(sn, pin, dt, eventcode, verify)
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}

	return true
}

func Device_Pin(d_id, pincode string) string {
	var emp string
	rows, err := db.Query("select employeeID from devicepin where deviceID=? and pinCode=?", d_id, pincode)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&emp)
		if err != nil {
			log.Print(err)
			emp = ""
		}
		return emp
	}
	return emp

}

func DeleteAllFP(emp_id string) bool {
	stmt, err := db.Prepare("delete from fingerprint where employeeID=?")
	if err != nil {
		log.Print(err)
		return false
	}
	_, err = stmt.Exec(emp_id)
	if err != nil {
		log.Print(err)
		return false
	}
	return true

}

func Update_pwd(emp_id, pwd string) bool {

	var s sql.NullString

	if pwd != "" {
		s.String = pwd
	}

	stmt, err := db.Prepare("update employee set passwd=? where ID=?")
	if err != nil {
		log.Print(err)
	}

	_, err = stmt.Exec(s, emp_id)
	if err != nil {
		//log.Print(err)
		return false
	}
	return true
}

func PinFromCompany(cmp_id, pin string) string {
	var emp string
	rows, err := db.Query("select id from employee where companyID=? and pinCode=? and state='active'", cmp_id, pin)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&emp)
		if err != nil {
			log.Print(err)
			emp = ""
		}
		return emp
	}
	return emp

}

func Device_Group(d_id string) int {
	var d_gr int
	rows, err := db.Query("select devicegroupid from devicetype where id in (select devicetypeid from device where id=?)", d_id)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&d_gr)
		if err != nil {
			log.Print(err)
		}
	}
	return d_gr
}

func Isactive(emp_id string) bool {
	var id string
	rows, err := db.Query("select id from employee where ID=? and state='active'", emp_id)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Print(err)
		}
	}
	if id == "" {
		logs.All("", "all", "Employee ID="+emp_id+" is not Active")
		return false
	} else {
		logs.All("", "all", "Employee ID="+emp_id+" is Active")
		return true
	}

}

func Isfinger(emp_id string) bool {
	var isf, pid string
	//rows, err := db.Query("select isFingerprint from policy where ID in (select policyID from employee where ID=?)", emp_id)
	rows, err := db.Query("select policyid from employee where id=?", emp_id)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&pid)
		if err != nil {
			log.Print(err)
		}
	}

	rows, err = db.Query("select isFingerprint from policy where ID=?", pid)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&isf)
		if err != nil {
			log.Print(err)
		}
	}

	if isf == "1" {
		return true
	} else {
		return false
	}

}

func Add_FP_Base(emp_id, fid int, fp []byte, st, devgr int) (bool, string) {
	//fmt.Println("Add FP Base EMP=", emp_id, "FID=", fid, fp, st, devgr)
	//sqlstr := "INSERT INTO fingerprint (employeeid, finger, fingerprint, state, devicegroupid) VALUES(" + emp_id + ", " + fid + ", '" + fp + "', " + st + ", " + devgr + ")"

	/*	stmt, err := db.Prepare("INSERT INTO `fingerprint` (employeeid, finger, fingerprint, state, devicegroupid) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			log.Print(err)
		}
		res, err := stmt.Exec(emp_id, fid, fp, st, devgr)
		fmt.Println(stmt)*/
	/*res, err := db.Exec(sqlstr, emp_id, 0, "")
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}*/

	type Finger struct {
		Id            int    `db:"ID"`
		Employeeid    int    `db:"employeeID"`
		Finger        int    `db:"finger"`
		Fingerprint   []byte `db:"fingerprint"`
		State         int    `db:"state"`
		Devicegroupid int    `db:"devicegroupID"`
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	//table := dbmap.AddTable(Finger{})
	dbmap.AddTableWithName(Finger{}, "fingerprint").SetKeys(true, "ID")
	ff := Finger{Employeeid: emp_id,
		Finger:        fid,
		Fingerprint:   fp,
		State:         st,
		Devicegroupid: devgr}
	err := dbmap.Insert(&ff)
	if err != nil {
		fmt.Println(err)
		logs.All_File("", "errors", err.Error())
		return false, err.Error()
	}

	return true, ""

}

func Get_location(d_id string) string {
	var loc_id string

	err := db.QueryRow("select locationid from device where id = ?", d_id).Scan(&loc_id)
	if err != nil {
		log.Print(err)
		return ""
	}
	return loc_id
}

func Get_loc_devices(loc_id string, d_gr int) []string {
	var devs []string

	rows, err := db.Query("select id from device where locationid=? and devicetypeid in (select id from devicetype where devicegroupid=?)", loc_id, d_gr)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmp_d string
		err := rows.Scan(&tmp_d)
		devs = append(devs, tmp_d)
		if err != nil {
			log.Print(err)
		}
	}

	return devs
}

func Update_photostamp(sn, phs string) bool {
	stmt, err := db.Prepare("update device set photoStamp=? where serialnumber=?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(phs, sn)
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}
	return true

}

func Update_stamp(sn, sta string) bool {
	const shrtFrm = "2006-01-02 15:04:05-0700 MST"
	dt := sta + "+0500 UZT"
	t2, _ := time.Parse(shrtFrm, dt)
	nt := t2.Local()
	nx := nt.Unix()

	stmt, err := db.Prepare("update device set attLogStamp=? where serialnumber=?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(nx, sn)
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}
	return true

}

func Update_opstamp(sn, op string) bool {
	stmt, err := db.Prepare("update device set operLogStamp=? where serialnumber=?")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(op, sn)
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}
	return true

}

func Get_server_id(d_id string) string {
	var s_id sql.NullString
	var str string

	err := db.QueryRow("select serverid from device where id = ?", d_id).Scan(&s_id)
	if err != nil {
		log.Print(err)
		return str
	}

	if s_id.Valid {
		str = s_id.String
	} else {
		str = ""
	}

	return str
}

func Add_Gprs_command(cmnd, d_id, pin, fid, fp string) bool {
	var cmd_content string
	if cmnd == "dataFp" {
		cmd_content = "DATA FP PIN=" + pin + "\tFID=" + fid + "\tValid=1\tTMP=" + fp
	}

	cmd_status := 1

	stmt, err := db.Prepare("INSERT INTO devicecmds (deviceID, CmdContent, CmdCommitTime, cmdStatus, pinCode, command) VALUES(?, ?, NOW(), ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(d_id, cmd_content, cmd_status, pin, cmnd)
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}
	return true

}

func Add_Server_command(sid, emp, sn, pin, fid, fpt string) bool {
	const LAN_COMMAND_ADD_FINGERPRINT_CODE = 5
	const LAN_COMMAND_MODE_ADD_CODE = 0
	code := "5"
	mode := "0"

	stmt, err := db.Prepare("INSERT INTO command (code, mode, serverID, serialNumber, employeeID, pinCode, finger, fingerprint, status) VALUES(?, ?, ?, ?, ?, ?, ?, ?, 0)")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(code, mode, sid, sn, emp, pin, fid, fpt)
	if err != nil {
		log.Print(err)
		log.Print(res)
		return false
	}

	return true
}

func InsertOplogData(sn, line string) bool {
	logs.All_File(sn, "oplog", line)
	res := false
	companyid := Comp_id(sn)
	d_id := Dev_id(sn)
	delim := strings.Split(line, " ")
	tip := delim[0]

	/*for i := range delim {
		fmt.Println("delim [", i, delim[i])
	}*/

	if tip == "FP" {
		keyvalue := strings.Split(delim[1], "\t")
		var pin, fid, tmp string
		for i := range keyvalue {
			//fmt.Println("keyvalue[", i, "]=", keyvalue[i])
			data := strings.Split(keyvalue[i], "=")
			if data[0] == "PIN" {
				pin = data[1]
			}
			if data[0] == "FID" {
				fid = data[1]
			}
			if data[0] == "TMP" {
				s := keyvalue[i]
				tmp = s[4:]
			}
		}
		emp_id := PinFromCompany(companyid, pin)
		//fmt.Println("My pin=", pin, "fid=", fid, "tmp=", tmp, emp_id)
		//fingerprint = strtoupper(bin2hex(base64_decode($data[4])))

		fp, err := base64.StdEncoding.DecodeString(tmp)
		if err != nil {
			fmt.Println("error while decoding base64:", err)
		}
		//fmt.Println("--- base64  ---", fp)
		fpt := strings.ToUpper(hex.EncodeToString(fp))
		//fmt.Println("--- HEX  ---", fpt)
		dev_group := Device_Group(d_id)
		var ext_name string
		/*const ZKT      = 1;
		  const TERMINAL   = 2;
		  const ANVIZ      = 3;*/
		if dev_group == 1 {
			// ZK Device
			ext_name = "_FPDbfile.fpt"
		}
		if dev_group == 3 {
			// Anviz device
			ext_name = "_FPDbfile.avz"
		}
		fname := "logs/fps/" + emp_id + "_" + fid + ext_name
		logs.Wr_file(fname, fpt)

		emp_id = Device_Pin(d_id, pin)

		if emp_id != "" {
			// FP ni bazaga saqlash
			e, _ := strconv.Atoi(emp_id)
			fi, _ := strconv.Atoi(fid)
			adding, errm := Add_FP_Base(e, fi, []byte(fpt), 1, dev_group)
			if adding {
				logs.All(sn, "all", "FP Added to base for employee="+emp_id+" FID="+fid)
				res = true
				if Isactive(emp_id) && Isfinger(emp_id) {
					loc_id := Get_location(d_id)
					if loc_id != "" {
						devs := Get_loc_devices(loc_id, dev_group)
						//fmt.Println("*- Devs **", devs)
						for i := range devs {
							devid := devs[i]
							sid := Get_server_id(devid)
							if sid == "" {
								if Add_Gprs_command("dataFp", devs[i], pin, fid, tmp) {
									res = true
									logs.All(sn, "all", "GPRS Command added for "+devs[i])
								} else {
									logs.All(sn, "all", "Error adding Command for "+devs[i])
									logs.All_File(sn, "errors", "Error adding Command for "+devs[i])
								}
							} else {
								sc := Add_Server_command(sid, emp_id, sn, pin, fid, fpt)
								if sc {
									logs.All(sn, "all", "Server command ID="+sid+" added successfully")
								} else {
									logs.All(sn, "all", "Error adding Server command ID="+sid)
									logs.All_File(sn, "errors", "Error adding Server command ID="+sid)
								}
							}

						}

					} else {
						logs.All(sn, "all", "Error location not found")
						logs.All_File(sn, "errors", "Error location not found")
					}

				} else {
					logs.All(sn, "all", "Error Employee "+emp_id+" is not active or is not finger")
					logs.All_File(sn, "errors", "Error Employee "+emp_id+" is not active or is not finger")
				}

			} else {
				logs.All(sn, "all", "Error Adding FP to base for employee="+emp_id+" FID="+fid+" "+errm)
				logs.All_File(sn, "errors", "Error Adding FP to base for employee="+emp_id+" FID="+fid+" "+errm)
			}

		}

	}
	if tip == "USER" {
		keyvalue := strings.Split(delim[1], "\t")
		var pin, passw string
		for i := range keyvalue {
			value := strings.Split(keyvalue[i], "=")
			if value[0] == "PIN" {
				pin = value[1]

			}
			if value[0] == "Passwd" {
				passw = value[1]
			}
			emp_id := PinFromCompany(companyid, pin)
			if passw != "" {
				if Update_pwd(emp_id, passw) {
					logs.All(sn, "all", "Password changed for PIN="+pin)
				}

			}

		}

	}
	if tip == "OPLOG" {
		delim[1] = delim[1] + " " + delim[2]
		var opl Oplog
		oplogs := strings.Split(delim[1], "\t")
		opl.Opcode = oplogs[0]
		opl.Adminid = oplogs[1]
		opl.DtTime = oplogs[2]
		opl.Obj1 = oplogs[3]
		opl.Obj2 = oplogs[4]
		opl.Obj3 = oplogs[5]
		opl.Obj4 = oplogs[6]

		emp_id := Device_Pin(d_id, opl.Obj1)

		if opl.Opcode == "10" {
			if emp_id != "" {
				if DeleteAllFP(emp_id) {
					logs.All(sn, "all", "All Fingerprints successfully deleted")
				} else {
					logs.All(sn, "all", "Error on deleting all Fingerprints")
					logs.All_File(sn, "errors", "Error on deleting all Fingerprints")
				}

			}

		}
		if opl.Opcode == "11" {
			Update_pwd(emp_id, "")
		}

	}

	return res
}

func Find_cmd(id string) []Cmd {
	var Cmds []Cmd
	rows, err := db.Query("select id, CmdContent from devicecmds where deviceID = ? AND commandTransferTime IS NULL AND cmdStatus = 1 LIMIT 1", id)
	//	rows, err := db.Query("select attLogStamp as Stamp from device where serialnumber=?", sn)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cmd_one Cmd
		err := rows.Scan(&cmd_one.Id, &cmd_one.Cmdbody)
		Cmds = append(Cmds, cmd_one)
		if err != nil {
			log.Print(err)
		}
	}

	return Cmds
}

func Update_Cmdstatus(sn, cmdid, rvalue, cmd string) {
	d_id := Dev_id(sn)
	if d_id != "" {
		var cmdstatus int
		if rvalue == "0" {
			cmdstatus = 0
		} else {
			cmdstatus = -1
		}
		//fmt.Println(cmdstatus)

		stmt, err := db.Prepare("update devicecmds set commandOverTime=NOW(), cmdStatus=?, Rvalue=?, returnCMD=? where id=?")
		if err != nil {
			log.Print(err)
		}
		res, err := stmt.Exec(cmdstatus, rvalue, cmd, cmdid)
		if err != nil {
			log.Print(err)
			log.Print(res)
			logs.All(sn, "all", "Error updating status Command ID="+cmdid)
		}
		logs.All(sn, "all", "Command ID="+cmdid+" status updated")

	}

}

func Companyname(sn string) string {
	var cmp_n string

	err := db.QueryRow("select name from company where id in (select companyid from device where serialnumber=?)", sn).Scan(&cmp_n)
	if err != nil {
		log.Print(err)
	}

	return cmp_n

}
