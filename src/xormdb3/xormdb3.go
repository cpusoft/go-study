package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"

	"github.com/cpusoft/goutil/iputil"
	"github.com/cpusoft/goutil/jsonutil"
	_ "github.com/go-sql-driver/mysql"
	model "labscm.zdns.cn/rpstir2-mod/rpstir2-model"
	"xorm.io/xorm"
)

func main() {
	//DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")
	user := "rpstir2"
	password := "Rpstir-123"
	server := "202.173.14.104:13306"
	database := "rpstir2"
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
	asn := 136825
	prefixLength := 22
	maxLength := 0
	address := `103.91.24.0`
	curSerialNumber := 1

	eng := engine.Table("lab_rpki_rtr_full_log").Where(" serialNumber= ? ", curSerialNumber)
	change := false

	if asn > 0 {
		change = true
		eng = eng.And(` asn = ? `, asn)
	}

	if prefixLength > 0 {
		change = true
		eng = eng.And(` prefixLength = ? `, prefixLength)
	}
	if maxLength > 0 {
		change = true
		eng = eng.And(` maxLength = ? `, maxLength)
	}
	if len(address) > 0 {
		change = true
		addressNew, _ := iputil.TrimAddressPrefixZero(address, iputil.GetIpType(address))
		eng = eng.And(` address = ? `, addressNew)
	}
	if !change {
		fmt.Println("getEffectSlurmsFromSlurmDb():not found delete condition from slurm, continue to next, :")
		return
	}

	// because asn may be nil or be 0, so using  sql.NullInt64
	type SlurmToRtrFullLog struct {
		Id             uint64        `json:"id" xorm:"id int"`
		Style          string        `json:"style" xorm:"style varchar(128)"`
		Asn            sql.NullInt64 `json:"asn" xorm:"asn int"`
		Address        string        `json:"address" xorm:"address varchar(256)"`
		PrefixLength   uint64        `json:"prefixLength" xorm:"prefixLength int"`
		MaxLength      uint64        `json:"maxLength" xorm:"maxLength int"`
		SlurmId        uint64        `json:"slurmId" xorm:"slurmId int"`
		SlurmLogId     uint64        `json:"slurmLogId" xorm:"slurmLogId int"`
		SlurmLogFileId uint64        `json:"slurmLogFileId" xorm:"slurmLogFileId int"`
	}
	filterSlurms := make([]model.EffectSlurmToRtrFullLog, 0)
	err = eng.Cols("asn,address,prefixLength,maxLength").Find(&filterSlurms)
	if err != nil {
		fmt.Println("getEffectSlurmsFromSlurmDb(): get lab_rpki_rtr_full_log fail:", err)
		return
	}
	fmt.Println("getEffectSlurmsFromSlurmDb():curSerialNumber:", curSerialNumber,
		"     filterSlurms:", jsonutil.MarshalJson(filterSlurms))

}
