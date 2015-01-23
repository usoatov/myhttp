package logs

import (
	"fmt"
	"log"
	"os"
)

func Wr_file(fl, text string) bool {
	f, err := os.OpenFile(fl, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer f.Close()

	_, err = f.WriteString(text + "\t")
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

func All(msg string) {
	file, err := os.OpenFile("logs/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file", file, ":", err)
	}
	//multi := io.MultiWriter(os.Stdout, file)
	log.SetOutput(file)
	log.Println(msg)
	log.SetOutput(os.Stdout)
	log.Println(msg)

}

func All_File(msg string) {
	file, err := os.OpenFile("logs/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file", file, ":", err)
	}
	//multi := io.MultiWriter(os.Stdout, file)
	log.SetOutput(file)
	log.Println(msg)
	log.SetOutput(os.Stdout)

}
