package main

import (
	_ "database/sql"
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func main() {
	//DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")
	//user := "dns"
	//password := "Dns-123"
	//server := "202.173.14.104:13307"
	//database := "dns"

	user := "root"
	password := "Rpstir-123"
	server := "10.1.135.22:13308"
	//database := "rpstir2"
	//maxidleconns := 50
	//maxopenconns := 50

	openSql := user + ":" + password + "@tcp(" + server + ")/" //+ database

	//连接数据库
	engine, err := xorm.NewEngine("mysql", openSql)
	if err != nil {
		fmt.Println(err)

		return
	}
	//连接测试
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ping ok")
	databases := make([]string, 0)

	err = engine.SQL("show databases").Find(&databases)
	if err != nil {
		fmt.Println("databases", err)
		return
	}
	fmt.Println("databases:", jsonutil.MarshalJson(databases))

	//engine.SQL("use mysql")
	for _, database := range databases {
		if database == "sys" || database == "information_schema" ||
			database == "performance_schema" || database == "mysql" {
			fmt.Println("ignore database:", database)
			continue
		}
		tables := make([]string, 0)
		tmp := "use " + database
		_, err = engine.Exec(tmp)
		if err != nil {
			fmt.Println("tables use fail:", err)
			return
		}
		err = engine.SQL("show tables").Find(&tables)
		if err != nil {
			fmt.Println("tables show fail:", err)
			return
		}
		// fmt.Println(database, " tables:", jsonutil.MarshalJson(tables))
		var userTable string
		var userCols []TableDescribe
		for _, table := range tables {
			fmt.Println(table)
			cols := make([]TableDescribe, 0)
			err = engine.SQL("describe " + table).Find(&cols)
			if err != nil {
				fmt.Println("table:", err)
				continue
			}
			// fmt.Println(table, "cols:", jsonutil.MarshalJson(cols), "\r\n")
			for _, col := range cols {

				if col.Field == "password" {
					fmt.Println(jsonutil.MarshalJson(col))
					userTable = table
					userCols = cols
					break
				}
			}
			if userTable != "" {
				break
			}
		}
		fmt.Println("found userTable:", userTable, "  userCols:", jsonutil.MarshalJson(userCols))
		sql := "select * from " + userTable
		results, err := engine.QueryString(sql)
		if err != nil {
			fmt.Println(err)
			return
		}
		var user, password string
		for _, result := range results {
			for k, v := range result {
				fmt.Println(k, v)
				if k == "user" || k == "userName" {
					user = v
				}
				if k == "password" {
					password = v
				}

			}
		}
		if user != "" && password != "" {
			fmt.Println(user, password)
		}
	}
}

type TableDescribe struct {
	Field   string `json:"field" xorm:"field varchar"`
	Type    string `json:"type" xorm:"type varchar"`
	Null    string `json:"null" xorm:"null varchar"`
	Key     string `json:"key" xorm:"key varchar"`
	Default string `json:"default" xorm:"default varchar"`
	Extra   string `json:"extra" xorm:"extra varchar"`
}
