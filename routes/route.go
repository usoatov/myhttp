package route

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	mydb "github.com/usoatov/myhttp/ds"
	logs "github.com/usoatov/myhttp/fl"
	"github.com/zenazn/goji/web"
)

func Getrequest(c web.C, w http.ResponseWriter, r *http.Request) {
	sn := r.URL.Query().Get("SN")
	//info := r.URL.Query().Get("INFO")
	d_id := mydb.Dev_id(sn)
	if d_id != "" {
		lr := mydb.Lastrequesttime(sn)
		if lr {
			logs.All(sn, "all", "Lastrequest")
		}
		Cmds := mydb.Find_cmd(d_id)
		eco := ""
		for i := range Cmds {
			tr := mydb.Transfertime(Cmds[i].Id)
			if tr {
				logs.All(sn, "all", "Command ID="+Cmds[i].Id+" transfered")
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
		logs.All(sn, "coms", eco)

	} else {
		logs.All("", "all", "Device "+d_id+" not found")

	}
}

func Cdata_get(c web.C, w http.ResponseWriter, r *http.Request) {
	var sn = r.URL.Query().Get("SN")
	var op = r.URL.Query().Get("options")
	d_id := mydb.Dev_id(sn)

	if d_id == "" {
		logs.All("", "all", "Device "+d_id+" not found")
	} else {
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

}

func Fdata_post(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//var aa = r.PostFormValue("aa")
	//var aa = r.Form["aa"][0]
	var sn = r.URL.Query().Get("SN")
	//var tbl = r.URL.Query().Get("table")
	var photostamp = r.URL.Query().Get("PhotoStamp")

	d_id := mydb.Dev_id(sn)
	if d_id == "" {
		logs.All("", "all", "Device "+d_id+" not found")
	} else {
		if mydb.Update_photostamp(sn, photostamp) {
			logs.All(sn, "all", "Photostamp changed to "+photostamp)

		}
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
		//logs.Wr_byte(dir+fn, phb)
		logs.Wr_byte("temp/tmp.jpg", phb)

		const shrtFrm = "20060102150405-0700 MST"
		dt := fn[:14]
		dt = dt + "+0500 UZT"
		t2, _ := time.Parse(shrtFrm, dt)
		nt := t2.Local()
		nx := nt.Unix()

		var pin string
		// ichida minus borligini tekshirish
		mb, _ := regexp.MatchString("-", fn)
		if mb {
			p := strings.Split(fn, "-")
			pn := strings.Split(p[1], ".")
			pin = pn[0]

		}
		fn = fmt.Sprintf("%d", nx)
		if pin != "" {
			fn = pin + "_" + fn + ".png"
		} else {
			fn = fn + ".png"
		}

		dataofimagejpg := ImageRead("temp/tmp.jpg", "jpeg")
		if dataofimagejpg != nil {
			if Formatpng(dataofimagejpg, dir+fn) {
				logs.All(sn, "all", "File "+dir+fn+" saved success.")
			} else {
				logs.All(sn, "errors", "Error saving file "+dir+fn)
			}

		} else {
			logs.All(sn, "errors", "Error reading photo tmp.jpg")

		}

		fmt.Fprintf(w, "OK")

	}
}

func ImageRead(ImageFile string, format string) (image image.Image) {
	if format == "jpeg" {
		file, err := os.Open(ImageFile)
		if err != nil {
			log.Println(err)
		}
		img, err := jpeg.Decode(file)
		if err != nil {
			log.Println(err)
		}
		file.Close()

		return img
	}
	return nil
}

func Formatpng(img image.Image, name string) bool {
	out, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true

}

func Cdata_post(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//var aa = r.PostFormValue("aa")
	//var aa = r.Form["aa"][0]
	var sn = r.URL.Query().Get("SN")
	//var tbl = r.URL.Query().Get("table")
	var stamp = r.URL.Query().Get("Stamp")
	var opstamp = r.URL.Query().Get("OpStamp")

	d_id := mydb.Dev_id(sn)
	if d_id == "" {
		logs.All("", "all", "Device "+d_id+" not found")
	} else {
		if stamp != "" {

			if mydb.Update_stamp(sn, stamp) {
				logs.All(sn, "all", "AttlogStamp changed to "+stamp)

			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
			}

			line := strings.Split(string(body), "\n")
			// Сатрларга булиш
			for j := range line {
				if line[j] != "" {
					/*if !strings.Contains(line[j], "2015") {
						continue
					}*/
					logs.All(sn, "all", "Processing line "+line[j])
					if line[j][:5] != "OPLOG" {
						r := mydb.InsertTempinout(sn, line[j])
						if r {
							logs.All(sn, "all", "Inserted inout "+line[j])
						} else {
							logs.All(sn, "all", "Error insering inout "+line[j])
							logs.All_File(sn, "errors", "Error insering inout "+line[j])
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

		if opstamp != "" {
			if mydb.Update_opstamp(sn, opstamp) {
				logs.All(sn, "all", "OpplogStamp changed to "+opstamp)

			}

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

}

func Devicecmd_post(c web.C, w http.ResponseWriter, r *http.Request) {
	var sn = r.URL.Query().Get("SN")
	d_id := mydb.Dev_id(sn)
	if d_id == "" {
		logs.All("", "all", "Device "+d_id+" not found")
	} else {
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

}
