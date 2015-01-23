package cfg

import (
	"log"

	"github.com/alyu/configparser"
)

func Read_config() (string, string, string, string, string) {
	var p, db, host, usr, pwd string
	config, err := configparser.Read("conf.ini")
	if err != nil {
		log.Println(err)
	}

	p = "1900"

	section, err := config.Section("Server")
	if err != nil {
		log.Fatal(err)
	} else {
		p = section.ValueOf("port")

	}
	section2, err := config.Section("Database")
	if err != nil {
		log.Fatal(err)
	} else {
		db = section2.ValueOf("db")
		host = section2.ValueOf("host")
		usr = section2.ValueOf("user")
		pwd = section2.ValueOf("pwd")

	}
	return p, db, host, usr, pwd
}
