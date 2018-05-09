package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//	"html/template"

	//	"net/http"
	"fmt"
	"time"
)

var DB *sql.DB

func main() {
	var err error
	DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	rows, err := DB.Query("select id, source, asn, ski, publicKey from RPKI_BGPSEC")
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var id, temp int
		var date time.Time
		err = rows.Scan(&id, &date, &temp)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(id, date, temp)
	}
}
