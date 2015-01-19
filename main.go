package main

import (
	"flag"
	"fmt"
	"net/http"

	cfg "github.com/usoatov/my_htt/cfg"
	mydb "github.com/usoatov/my_htt/ds"
	route "github.com/usoatov/my_htt/routes"
	"github.com/zenazn/goji"
)

var (
	port, db, host, usr, pwd string
	dbcon                    bool
)

func main() {
	port, db, host, usr, pwd = cfg.Read_config()

	dbcon := mydb.Connect(db, host, usr, pwd)
	if dbcon {
		fmt.Println("connect")
	}

	goji.Get("/iclock/cdata", route.Cdata_get)
	goji.Get("/iclock/getrequest", route.Getrequest)
	goji.Post("/iclock/cdata", route.Cdata_post)
	goji.Post("/iclock/devicecmd", route.Devicecmd_post)
	flag.Set("bind", ":"+port)
	goji.Use(PlainText)
	goji.Serve()

}

// PlainText sets the content-type of responses to text/plain.
func PlainText(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		/*now := time.Now()
		yr := now.Year()
		mon := now.Month()
		day := now.Day()
		soat := now.Hour()
		min := now.Minute()
		sec := now.Second()
		t := time.Date(yr, mon, day, soat, min, sec, 0, time.Local)
		// layout shows by example how the reference time should be represented.
		const layout = "Mon, 02 Jan 2006 15:04:05"
		ts := t.UTC().Format(layout)
		fmt.Println(ts)
		w.Header().Set("Date", ts)*/
		w.Header().Set("Content-Type", "text/plain")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
