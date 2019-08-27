package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

type RtrFull struct {
	SerialNumber       int    `json:"serialNumber" xorm:"serial_num INT"`
	Asn             int    `json:"asn" xorm:"asn INT"`
	Prefix          []byte `json:prefix" xorm:"prefix VARBINARY"`
	PrefixLength    int    `json:prefix_length" xorm:"prefix_length TINYINT"`
	PrefixMaxLength int    `json:prefix_max_length" xorm:"prefix_max_length TINYINT"`
}

type RtrFullCopy struct {
	SerialNumber       int    `json:"serialNumber" xorm:"serial_num INT"`
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
	/*
	   f, err := os.Create("sql.log")
	   	if err != nil {
	   		log.Fatalf("Fail to create log file: %v\n", err)
	   		return
	   	}
	   	logg := xorm.NewSimpleLogger(f)
	*/

	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(5)
	//设置最大打开连接数
	engine.SetMaxOpenConns(5)
	/*
		http://blog.xorm.io/2016/1/4/1-about-mapper.html
		SnakeMapper
		SnakeMapper是默认的映射机制，他支持数据库表采用匈牙利命名法，而程序中采用驼峰式命名法。下面是一些常见的映射：
		表中名称		程序名称
		user_info	UserInfo
		id			Id

		SameMapper
		SameMapper就是数据库中的命名法和程序中是相同的。那么鉴于在Go中，基本上要求首字母必须大写。所以一般都是表中和程序中均采用驼峰式命名。下面是一些常见的映射：
		表中名称	程序名称
		UserInfo	UserInfo
		Id	Id


		GonicMapper
		GonicMapper是在SnakeMapper的基础上增加了特例，对于常见的缩写不新增下划线处理。这个同时也符合golint的规则。下面是一些常见的映射：
		表中名称	程序名称
		user_info	UserInfo
		id	ID
		url	URL

	*/
	engine.SetTableMapper(core.SnakeMapper{})

	fmt.Println(engine)

	gsql := "SELECT `serial_num` FROM `rtr_update` order by  create_time desc limit 1"
	fmt.Println(gsql)
	gres, gerr := engine.Query(gsql)
	if gerr != nil {
		panic(gerr)
	}
	maxSerialNum := -1
	for _, v := range gres {
		maxSerialNum, _ = strconv.Atoi(string(v["serial_num"]))
		fmt.Println(maxSerialNum)
	}
	if maxSerialNum < 0 {
		return
	}

	var rtrFull []RtrFull
	err = engine.Where("serial_num=?", maxSerialNum).OrderBy("asn").Find(&rtrFull)
	if err != nil {
		fmt.Println(err)
	}
	//j, _ := json.Marshal(rtrFull)
	//fmt.Println(string(j))

	// batch insert
	session := engine.NewSession()
	defer session.Close()

	if err = session.Begin(); err != nil {
		return
	}

	fmt.Println("batch insert :", len(rtrFull))
	for index, rtrOne := range rtrFull {
		rtrFullCopy := RtrFullCopy{}
		j, _ := json.Marshal(rtrOne)
		json.Unmarshal(j, &rtrFullCopy)
		fmt.Println("insert :", index, "  rtrFullCopy:", string(j))
		affected, err := engine.Insert(rtrFullCopy)
		if affected < 0 || err != nil {
			fmt.Println(err)
			session.Rollback()
			return
		}
	}

	err = session.Commit()
	if err != nil {

		return
	}
}
