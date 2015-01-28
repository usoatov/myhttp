package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	cfg "github.com/usoatov/my_htt/cfg"
	mydb "github.com/usoatov/my_htt/ds"
	"github.com/usoatov/my_htt/fl"
	route "github.com/usoatov/my_htt/routes"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"github.com/zenazn/goji/web/mutil"
)

var (
	port, db, host, usr, pwd string
	dbcon                    bool
)

func main() {
	port, db, host, usr, pwd = cfg.Read_config()

	dbcon := mydb.Connect(db, host, usr, pwd)
	if dbcon {
		logs.All("", "all", "connected to db")
	}

	goji.Get("/iclock/cdata", route.Cdata_get)
	goji.Get("/iclock/getrequest", route.Getrequest)
	goji.Post("/iclock/cdata", route.Cdata_post)
	// photo ni olish
	goji.Post("/iclock/fdata", route.Fdata_post)
	goji.Post("/iclock/devicecmd", route.Devicecmd_post)
	flag.Set("bind", ":"+port)
	goji.Use(PlainText)
	goji.Use(MyLogger)
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
// Implementation of goji logger to process all request
func MyLogger(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(*c)

		var sn = r.URL.Query().Get("SN")
		printStart(sn, reqID, r)

		lw := mutil.WrapWriter(w)

		t1 := time.Now()
		h.ServeHTTP(lw, r)

		if lw.Status() == 0 {
			lw.WriteHeader(http.StatusOK)
		}
		t2 := time.Now()

		printEnd(sn, reqID, lw, t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func printStart(sn, reqID string, r *http.Request) {

	mystr := "[" + reqID + "] Started " + r.Method + " " + r.URL.String() + " from " + r.RemoteAddr

	logs.All_File(sn, "all", mystr)
	logs.All_File(sn, "reqs", mystr)

}

func printEnd(sn, reqID string, w mutil.WriterProxy, dt time.Duration) {
	var ss string

	status := w.Status()
	ss = ss + "Returning "
	ss = ss + fmt.Sprintf("%03d in %s", status, dt)

	logs.All_File(sn, "all", ss)
	logs.All_File(sn, "reqs", ss)

}
