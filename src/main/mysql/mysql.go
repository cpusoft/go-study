package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//	"html/template"
	"log"
	//	"net/http"
	"fmt"
)

var (
	serial_num        int
	asn               int
	prefix_length     int
	prefix_max_length int
)

var Db *sql.DB

func main() {
	var minNum sql.NullInt64
	minNum.Valid = false
	fmt.Println(minNum)
	Db, _ = sql.Open("mysql", "rpstir:Rpstir-123@tcp(192.168.138.135:3306)/rpstir")

	db, err := sql.Open("mysql", "rpstir:Rpstir-123@tcp(192.168.138.135:3306)/rpstir")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select serial_num,asn, prefix_length,prefix_max_length from slurm_target_2 where asn = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&serial_num, &asn, &prefix_length, &prefix_max_length)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(serial_num, asn, prefix_length, prefix_max_length)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
