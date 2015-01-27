package route

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	mydb "github.com/usoatov/my_htt/ds"
	logs "github.com/usoatov/my_htt/fl"
	"github.com/zenazn/goji/web"
)

func Getrequest(c web.C, w http.ResponseWriter, r *http.Request) {
	sn := r.URL.Query().Get("SN")
	//info := r.URL.Query().Get("INFO")
	d_id := mydb.Dev_id(sn)
	lr := mydb.Lastrequesttime(sn)
	if lr {
		logs.All_File(sn, "all", "Lastrequest")
	}
	Cmds := mydb.Find_cmd(d_id)
	eco := ""
	for i := range Cmds {
		tr := mydb.Transfertime(Cmds[i].Id)
		if tr {
			logs.All_File(sn, "all", "Command ID="+Cmds[i].Id+" transfered")
		}
		eco += "C:" + Cmds[i].Id + ":" + Cmds[i].Cmdbody + "\n"

		//fmt.Println("Cmd[id]=", Cmds[i].Id)
		//fmt.Println("Cmd[Cmdbody]=", Cmds[i].Cmdbody)
	}
	if eco == "" {
		eco = "OK"
		logs.All(sn, "all", "No Commands")
	}
	fmt.Fprintf(w, eco)
	logs.All_File(sn, "coms", eco)
}

func Cdata_get(c web.C, w http.ResponseWriter, r *http.Request) {
	var sn = r.URL.Query().Get("SN")
	var op = r.URL.Query().Get("options")
	lr := mydb.Lastrequesttime(sn)
	if lr {
		logs.All(sn, "all", "Lastrequest")
	}

	if op == "all" {
		// mysql connect
		opt := mydb.Options(sn)
		transflag := "TransData AttLog\tOpLog\tEnrollUser\tChgUser\tChgFP\tAttPhoto\tEnrollFP"

		resp := fmt.Sprintf("GET OPTION FROM: %s\n", sn)
		resp += fmt.Sprintf("Stamp=%s\n", opt.Stamp)
		resp += fmt.Sprintf("OpStamp=%s\n", opt.Opstamp)
		resp += fmt.Sprintf("PhotoStamp=%s\n", opt.Photostamp)
		resp += fmt.Sprintf("ErrorDelay=%s\n", opt.Errdel)
		resp += fmt.Sprintf("Delay=%s\n", opt.Delay)
		resp += fmt.Sprintf("TransTimes=%s\n", opt.Transtime)
		resp += fmt.Sprintf("TransInterval=%s\n", opt.Transint)
		resp += fmt.Sprintf("TransFlag=%s\n", transflag)
		resp += fmt.Sprintf("TimeZone=%s\n", opt.Timezone)
		resp += fmt.Sprintf("Realtime=%s\n", opt.Realtime)
		resp += fmt.Sprintf("Encrypt=%s\n", opt.Encrypt)

		fmt.Fprintf(w, resp)
		logs.All_File(sn, "all", resp)
	}
}

func Fdata_post(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//var aa = r.PostFormValue("aa")
	//var aa = r.Form["aa"][0]
	var sn = r.URL.Query().Get("SN")
	//var tbl = r.URL.Query().Get("table")
	//var photostamp = r.URL.Query().Get("PhotoStamp")
	//fmt.Println("====== PHOTO =======")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	line := strings.SplitN(string(body), "\n", 4)
	/*fmt.Println(sn, tbl, photostamp, "PH Line")
	for i := range line {
		fmt.Println("line[", i, "]=", line[i])
	}
	l := len(line)*/

	ls := strings.Split(line[0], "=")
	fn := ls[1]
	ph := strings.SplitN(line[3], "\u0000", 2)
	phb := []byte(ph[1])

	dir := "../../Apache/home/images/"
	t := time.Now()
	yr := fmt.Sprintf("%04d", t.Year())
	cmp_name := mydb.Companyname(sn)
	dir = dir + yr + "/" + cmp_name + "/" + sn + "/"

	logs.All(sn, "all", "Saving photo to file: "+dir+fn)
	logs.Wr_byte(dir+fn, phb)
	fmt.Fprintf(w, "OK")

}

func Cdata_post(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//var aa = r.PostFormValue("aa")
	//var aa = r.Form["aa"][0]
	var sn = r.URL.Query().Get("SN")
	var tbl = r.URL.Query().Get("table")
	var stamp = r.URL.Query().Get("Stamp")
	var opstamp = r.URL.Query().Get("OpStamp")

	if tbl != "" && stamp != "" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}

		line := strings.Split(string(body), "\n")

		// Сатрларга булиш
		for j := range line {
			if line[j] != "" {
				logs.All(sn, "all", "Processing line "+line[j])
				if line[j][:5] != "OPLOG" {
					r := mydb.InsertTempinout(sn, line[j])
					if r {
						logs.All(sn, "all", "Inserted inout "+line[j])
					} else {
						logs.All(sn, "all", "Error insering inout "+line[j])
						logs.All(sn, "errors", "Error insering inout "+line[j])
					}
				} else {
					r := mydb.InsertOplogData(sn, line[j])
					if r {
						logs.All(sn, "all", "Inserted OPLOG"+line[j])
					}
				}
			}
		}

	}

	if tbl != "" && opstamp != "" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logs.All(sn, "errors", "Error while trying receiving post Body!")
		}
		line := strings.Split(string(body), "\n")

		for j := range line {
			if line[j] != "" {
				logs.All(sn, "all", "Processing line "+line[j])
				r := mydb.InsertOplogData(sn, line[j])
				if r {
					fmt.Println("Inserted")
				}
			}

		}

	}

	fmt.Fprintf(w, "OK")
}

func Devicecmd_post(c web.C, w http.ResponseWriter, r *http.Request) {
	var sn = r.URL.Query().Get("SN")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Body")
	//fmt.Printf("%# v", body)

	line := strings.Split(string(body), "\n")
	for j := range line {
		if line[j] != "" {
			logs.All(sn, "all", "Receiving command status "+line[j])
			p := strings.Split(line[j], "&")
			var id, ret, cmd string
			for i := range p {
				v := strings.Split(p[i], "=")
				if v[0] == "ID" {
					id = v[1]
				}
				if v[0] == "Return" {
					ret = v[1]
				}
				if v[0] == "CMD" {
					cmd = v[1]
				}

			}
			mydb.Update_Cmdstatus(sn, id, ret, cmd)
		}
	}
	fmt.Fprintf(w, "OK")

}
