package main

import (
	_ "database/sql"
	"fmt"

	belogs "github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

func main() {
	//DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")
	user := "rpstir"
	password := "Rpstir-123"
	server := "202.173.14.105:13306"
	database := "rpstir_test"
	maxidleconns := 50
	maxopenconns := 50

	openSql := user + ":" + password + "@tcp(" + server + ")/" + database

	//连接数据库
	engine, err := xorm.NewEngine("mysql", openSql)
	if err != nil {
		fmt.Println(err)
		belogs.Error("NewEngine failed: ", err)
		return
	}
	//连接测试
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		belogs.Error("Ping failed: ", err)
		return
	}

	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(maxidleconns)
	//设置最大打开连接数
	engine.SetMaxOpenConns(maxopenconns)
	engine.SetTableMapper(core.SnakeMapper{})

	session := engine.NewSession()
	defer session.Close()
	type RsyncFileHash struct {
		FilePath    string `json:"filePath" xorm:"filePath varchar(512)"`
		FileName    string `json:"fileName" xorm:"fileName varchar(128)"`
		FileHash    string `json:"fileHash" xorm:"fileHash varchar(512)"`
		JsonAll     string `json:"jsonAll" xorm:"jsonAll json"`
		LastJsonAll string `json:"lastJsonAll" xorm:"lastJsonAll json"`
		FileType    string `json:"jsonAll" xorm:"fileType  varchar(16)"`
	}

	cerFileHashs := make([]RsyncFileHash, 15000)
	err = engine.Table("lab_rpki_cer").
		Select("filePath , fileName, fileHash, jsonAll as lastJsonAll, 'cer' as fileType").
		Asc("id").Find(&cerFileHashs)
	if err != nil {
		fmt.Println(err)
		belogs.Error("GetFilesHashFromDb(): get lab_rpki_cer fail:", err)
		return
	}
	fmt.Println(len(cerFileHashs))
}
