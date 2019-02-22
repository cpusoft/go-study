package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

type RtrFull struct {
	SerialNum       int    `json:"serialNum" xorm:"serial_num INT"`
	Asn             int    `json:"asn" xorm:"asn INT"`
	Prefix          []byte `json:prefix" xorm:"prefix VARBINARY"`
	PrefixLength    int    `json:prefix_length" xorm:"prefix_length TINYINT"`
	PrefixMaxLength int    `json:prefix_max_length" xorm:"prefix_max_length TINYINT"`
}

func main() {
	//DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")

	server := "202.173.14.102:13306"
	user := "rpstir"
	password := "Rpstir-123"
	database := "rpstir"
	openSql := user + ":" + password + "@tcp(" + server + ")/" + database
	//params := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", "root", "root", "localhost:3306", "go")
	//连接数据库
	fmt.Println(openSql)
	engine, err := xorm.NewEngine("mysql", openSql)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer engine.Close()
	//连接测试
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()

	//日志打印SQL
	engine.ShowSQL(true)

	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(5)
	//设置最大打开连接数
	engine.SetMaxOpenConns(5)
	engine.SetTableMapper(core.SnakeMapper{})

	fmt.Println(engine)

	gsql := "SELECT `serial_num` FROM `rtr_update` order by  create_time desc limit 1"
	fmt.Println(gsql)
	gres, gerr := engine.Query(gsql)
	if gerr != nil {
		panic(gerr)
	}

	for _, v := range gres {
		maxSerialNum := (string(v["serial_num"]))
		fmt.Println(maxSerialNum)
	}

	var rtrFull []RtrFull
	err = engine.Where("serial_num=?", 8).OrderBy("asn").Find(&rtrFull)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(rtrFull)
	fmt.Println(string(j))

}
