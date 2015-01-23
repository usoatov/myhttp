package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	logs "github.com/usoatov/my_htt/fl"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/mutil"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white.
//
// Logger prints a request ID if one is provided.
//
// Logger has been designed explicitly to be Good Enough for use in small
// applications and for people just getting started with Goji. It is expected
// that applications will eventually outgrow this middleware and replace it with
// a custom request logger, such as one that produces machine-parseable output,
// outputs logs to a different service (e.g., syslog), or formats lines like
// those printed elsewhere in the application.
func Logger(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqID := GetReqID(*c)

		printStart(reqID, r)

		lw := mutil.WrapWriter(w)

		t1 := time.Now()
		h.ServeHTTP(lw, r)

		if lw.Status() == 0 {
			lw.WriteHeader(http.StatusOK)
		}
		t2 := time.Now()

		printEnd(reqID, lw, t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func printStart(reqID string, r *http.Request) {
	var buf bytes.Buffer

	if reqID != "" {
		cW(&buf, bBlack, "[%s] ", reqID)
	}
	buf.WriteString("Started ")
	cW(&buf, bMagenta, "%s ", r.Method)
	cW(&buf, nBlue, "%q ", r.URL.String())
	buf.WriteString("from ")
	buf.WriteString(r.RemoteAddr)

	log.Print(buf.String())

	file, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file", file, ":", err)
	}
	//multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(file)
	mystr := "[" + reqID + "] Started " + r.Method + " " + r.URL.String() + " from " + r.RemoteAddr

	log.Print(mystr)
	log.SetOutput(os.Stdout)
	logs.All_File(mystr)

}

func printEnd(reqID string, w mutil.WriterProxy, dt time.Duration) {
	var buf bytes.Buffer
	var ss string

	if reqID != "" {
		cW(&buf, bBlack, "[%s] ", reqID)
		ss = "[" + reqID + "]"
	}
	buf.WriteString("Returning ")
	status := w.Status()
	if status < 200 {
		cW(&buf, bBlue, "%03d", status)
	} else if status < 300 {
		cW(&buf, bGreen, "%03d", status)
	} else if status < 400 {
		cW(&buf, bCyan, "%03d", status)
	} else if status < 500 {
		cW(&buf, bYellow, "%03d", status)
	} else {
		cW(&buf, bRed, "%03d", status)
	}
	buf.WriteString(" in ")
	if dt < 500*time.Millisecond {
		cW(&buf, nGreen, "%s", dt)
	} else if dt < 5*time.Second {
		cW(&buf, nYellow, "%s", dt)
	} else {
		cW(&buf, nRed, "%s", dt)
	}
	ss = ss + "Returning "
	ss = ss + fmt.Sprintf("%03d in %s", status, dt)

	log.Print(buf.String())

	file, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file", file, ":", err)
	}
	//multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(file)
	log.Print(ss)
	//log.Print(buf.String())
	log.SetOutput(os.Stdout)

}
