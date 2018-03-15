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
	DB, err = sql.Open("mysql", "devops:knet@tcp(202.173.9.53:3306)/devops?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	rows, err := DB.Query("select id, date, temp from temp where temp=?", 10)
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
