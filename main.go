package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func cdata_get(c web.C, w http.ResponseWriter, r *http.Request) {
	var sn = r.URL.Query().Get("SN")
	var op = r.URL.Query().Get("options")
	fmt.Fprintf(w, "Hello, %s, %s!", sn, op)
}

func cdata_post(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var aa = r.PostFormValue("aa")
	//var aa = r.Form["aa"][0]
	var sn = r.URL.Query().Get("SN")
	fmt.Fprintf(w, "Get, %s!", sn)

	/*fmt.Println(r.Form)
	for i:=range(r.Form) {
		fmt.Println(r.Form[i][0] )
	}*/

	fmt.Fprintf(w, "Hello, %s !", aa)
}

func main() {
	goji.Get("/iclock/cdata", cdata_get)
	goji.Post("/iclock/cdata", cdata_post)
	//goji.bind.default = 1999
	flag.Set("bind", ":1234")
	goji.Serve()

}
