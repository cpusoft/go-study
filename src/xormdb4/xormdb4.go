package main

import (
	_ "database/sql"
	"fmt"
	"time"

	"github.com/cpusoft/goutil/dnsutil"
	"github.com/cpusoft/goutil/jsonutil"
	"github.com/guregu/null"

	"xorm.io/xorm"
)

// for mysql, zonefile
type RrModel struct {
	Id       uint64 `json:"id" xorm:"id int"`
	OriginId uint64 `json:"originId" xorm:"originId int"`

	// not have "." in the end
	Origin string `json:"origin" xorm:"origin varchar"` // lower
	// is host/subdomain, not have "." int the end
	// if no subdomain, is "", not "@"
	RrName string `json:"rrName" xorm:"rrName varchar"` // lower
	// == rrName+.+Origin
	RrFullDomain string `json:"rrFullDomain" xorm:"rrFullDomain varchar"` // lower: rrName+"."+Origin[-"."]

	RrType  string `json:"rrType" xorm:"rrType varchar"`   // upper
	RrClass string `json:"rrClass" xorm:"rrClass varchar"` // upper
	// null.NewInt(0, false) or null.NewInt(i64, true)
	RrTtl  null.Int `json:"rrTtl" xorm:"rrTtl int"`
	RrData string   `json:"rrData" xorm:"rrData varchar"`

	UpdateTime time.Time `json:"updateTime" xorm:"updateTime datetime"`
}

func main() {
	//DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")
	user := "dns"
	password := "Dns-123"
	server := "202.173.14.104:13307"
	database := "dns"
	maxidleconns := 50
	maxopenconns := 50

	openSql := user + ":" + password + "@tcp(" + server + ")/" + database

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

	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(maxidleconns)
	//设置最大打开连接数
	engine.SetMaxOpenConns(maxopenconns)
	//engine.SetTableMapper(core.SnakeMapper{})
	engine.ShowSQL(true)

	session := engine.NewSession()
	defer session.Close()

	rrType := "A"
	rrFullDomain := "dns1.example.com"
	resultRrModels := make([]*RrModel, 0)
	fmt.Println("queryDb(): rrType:", rrType, "  rrFullDomain:", rrFullDomain)
	if rrType == dnsutil.DNS_TYPE_STR_ANY {
		sql := `select o.origin, r.rrName, r.rrFullDomain, r.rrType, r.rrClass, IFNULL(r.rrTtl,o.ttl) as rrTtl, r.rrData  
			from lab_dns_rr r,	lab_dns_origin o 
			where r.originId = o.id and r.rrFullDomain = ?  
			group by r.id `
		err = engine.SQL(sql, rrFullDomain).Find(&resultRrModels)
	} else {
		sql := `select o.origin, r.rrName, r.rrFullDomain, r.rrType, r.rrClass, IFNULL(r.rrTtl,o.ttl) as rrTtl, r.rrData  
			from lab_dns_rr r, lab_dns_origin o 
			where r.originId = o.id and r.rrFullDomain = ? and r.rrType = ? 
			group by r.id `
		err = engine.SQL(sql, rrFullDomain, rrType).Find(&resultRrModels)
	}
	if err != nil {
		fmt.Println("queryDb(): lab_dns_rr fail:", err)
		return
	}
	fmt.Println("queryDb(): resultRrModels:", jsonutil.MarshalJson(resultRrModels))

}
