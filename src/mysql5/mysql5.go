package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	//	"html/template"
	//	"net/http"
)

var DB *sql.DB

func main() {
	var err error
	DB, err = sql.Open("mysql", "rpstir2:Rpstir-123@tcp(127.0.0.1:13306)/rpstir2")
	//DB, err = sql.Open("mysql", "rpstir2:Rpstir-123@tcp(172.17.0.2:3306)/rpstir2")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	DB.SetMaxOpenConns(3)
	DB.SetMaxIdleConns(0)

	err = DB.Ping()
	if err != nil {
		log.Printf("数据库连接出错:%s\n", err.Error())
		return
	}

	var value int
	err = DB.QueryRow("select 1").Scan(&value)
	if err != nil {
		log.Println("query failed:", err)
		return
	}

	log.Println("value:", value)

}
