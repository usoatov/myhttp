package logs

import (
	"fmt"
	"log"
	"os"
	"time"
)

func Wr_file(fl, text string) bool {
	f, err := os.OpenFile(fl, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer f.Close()

	_, err = f.WriteString(text + "\n")
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}

func Inout(sn, text string) bool {
	t := time.Now()
	s := fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
	if sn == "" {
		s = "logs/inout/" + s + ".log"
	} else {
		s = "logs/inout/" + sn + "_" + s + ".log"
	}
	f, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file", f, ":", err)
	}

	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}

func Wr_byte(fl string, bb []byte) bool {
	f, err := os.OpenFile(fl, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer f.Close()

	_, err = f.Write(bb)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}

func All(sn, dr, msg string) {
	t := time.Now()
	s := fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
	if sn == "" {
		s = "logs/" + dr + "/" + s + ".log"
	} else {
		s = "logs/" + dr + "/" + sn + "_" + s + ".log"
	}
	file, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file", file, ":", err)
	}
	log.SetOutput(file)
	log.Println(msg)
	log.SetOutput(os.Stdout)
	log.Println(msg)

}

func All_File(sn, dr, msg string) {
	t := time.Now()
	s := fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
	if sn == "" {
		s = "logs/" + dr + "/" + s + ".log"
	} else {
		s = "logs/" + dr + "/" + sn + "_" + s + ".log"
	}
	file, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file", file, ":", err)
	}
	log.SetOutput(file)
	log.Println(msg)
	log.SetOutput(os.Stdout)

}
