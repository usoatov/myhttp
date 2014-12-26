package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var usr, pasw, datb string = "root", "", "biotrackonline"
	con, err := sql.Open("mysql", usr+":"+pasw+"@/"+datb)
	if err != nil {
		fmt.Println("error")
	} else {
		fmt.Println("connected")
	}
	defer con.Close()

	id := 31
	row := con.QueryRow("select id, name from company where id=?", id)

	var cid int
	var name string

	err = row.Scan(&cid, &name)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println("Cid=", cid, "name=", name)



	var names []string

	rows, err := con.Query("select name FROM company WHERE pricePlanID=?", "7")

	if err != nil {
	fmt.Println(err)
	return
	}

	for rows.Next() {
	var nm string
	if err := rows.Scan(&nm); err != nil {
	fmt.Println(err)
	return
	}
	names = append(names, nm)
	}

	for i:=range(names) {
		fmt.Println(names[i])
	}	
}
