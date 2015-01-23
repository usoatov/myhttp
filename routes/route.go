package route

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	mydb "github.com/usoatov/my_htt/ds"
	logs "github.com/usoatov/my_htt/fl"
	"github.com/zenazn/goji/web"
)

func Getrequest(c web.C, w http.ResponseWriter, r *http.Request) {
	sn := r.URL.Query().Get("SN")
	info := r.URL.Query().Get("INFO")
	d_id := mydb.Dev_id(sn)
	fmt.Println("Dev_id=", d_id, "Info", info)
	lr := mydb.Lastrequesttime(sn)
	if lr {
		fmt.Println("Lastrequest bajarildi")
	}
	Cmds := mydb.Find_cmd(d_id)
	eco := ""
	for i := range Cmds {
		tr := mydb.Transfertime(Cmds[i].Id)
		if tr {
			fmt.Println("Trans bajarildi")
		}
		eco += "C:" + Cmds[i].Id + ":" + Cmds[i].Cmdbody + "\n"

		//fmt.Println("Cmd[id]=", Cmds[i].Id)
		//fmt.Println("Cmd[Cmdbody]=", Cmds[i].Cmdbody)
	}
	if eco == "" {
		eco = "OK"
		fmt.Println("Bosh")
	}
	fmt.Fprintf(w, eco)
}

func Cdata_get(c web.C, w http.ResponseWriter, r *http.Request) {
	var sn = r.URL.Query().Get("SN")
	var op = r.URL.Query().Get("options")
	lr := mydb.Lastrequesttime(sn)
	if lr {
		fmt.Println("Lastreq Bajarildi")
	}
	//fmt.Fprintf(w, "Hello, %s, %s!", sn, op)
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
	}
}

func Fdata_post(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//var aa = r.PostFormValue("aa")
	//var aa = r.Form["aa"][0]
	var sn = r.URL.Query().Get("SN")
	var tbl = r.URL.Query().Get("table")
	var photostamp = r.URL.Query().Get("PhotoStamp")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	line := strings.SplitN(string(body), "\n", 4)
	fmt.Println(sn, tbl, photostamp, "PH Line")
	/*for i := range line {
		fmt.Println("line[", i, "]=", line[i])
	}*/
	//l := len(line)
	ph := strings.SplitN(line[3], "\u0000", 2)
	for i := range ph {
		fmt.Println("ph[", i, ph[i])
	}
	phb := []byte(ph[1])
	logs.Wr_byte("photo.jpg", phb)

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
				if line[j][:5] != "OPLOG" {
					r := mydb.InsertTempinout(sn, line[j])
					if r {
						fmt.Println("Inserting ...")
					}
				} else {
					r := mydb.InsertOplogData(sn, line[j])
					if r {
						fmt.Println("Ins OPLOG")
					}
				}
			}
		}

	}

	if tbl != "" && opstamp != "" {
		fmt.Println("Opstamp=", opstamp)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		line := strings.Split(string(body), "\n")

		for j := range line {
			if line[j] != "" {
				r := mydb.InsertOplogData(sn, line[j])
				if r {
					fmt.Println("")
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
			fmt.Println("line=", j, line[j])
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
			fmt.Println("id=", id, "ret=", ret, "cmd=", cmd)
			mydb.Update_Cmdstatus(sn, id, ret, cmd)
		}
	}
	fmt.Fprintf(w, "OK")

}
