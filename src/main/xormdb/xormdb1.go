package main

import (
	_ "database/sql"
	"fmt"
	_ "time"

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
	/*
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
	*/
	/*
		// get current rsyncState, the set new value
		var id int64
		_, err = session.Table("lab_rpki_sync_log").Select("max(id)").Get(&id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(id)
	*/
	//lab_rpki_sync_log_file
	/*
		type LabRpkiSyncLogFile struct {
			Id        uint64 `json:"id" xorm:"pk autoincr"`
			SyncLogId uint64 `json:"syncLogId" xorm:"syncLogId int"`
			//cer/roa/mft/crl, not dot
			FileType string `json:"fileType" xorm:"fileType varchar(16)"`
			//sync time for every file
			SyncTime                time.Time `json:"syncTime" xorm:"syncTime datetime"`
			Ski                     string    `json:"ski" xorm:"ski varchar(128)"`
			Aki                     string    `json:"aki" xorm:"aki varchar(128)"`
			FilePath                string    `json:"filePath" xorm:"filePath varchar(512)"`
			FileName                string    `json:"fileName" xorm:"fileName varchar(128)"`
			ParseValidateResultJson string    `json:"syncJson" xorm:"syncJson json"`
			JsonAll                 string    `json:"jsonAll" xorm:"jsonAll json"`
			LastJsonAll             string    `json:"jsonAll" xorm:"lastJsonAll json"`
			FileHash                string    `json:"fileHash" xorm:"fileHash varchar(512)"`
			//add/update/del
			SyncType string `json:"syncType" xorm:"syncType varchar(16)"`
			//LabRpkiSyncLogFileState:
			State string `json:"state" xorm:"state json"`
		}
		syncLogFiles := make([]LabRpkiSyncLogFile, 0)
		//err = xormdb.XormEngine.Table(&syncLogFile).Where("syncLogId = ?", syncLog.Id).And("state->>'$.updateCertTable'=?", "notYet").Asc("id").Find(&syncLogFiles)

		err = engine.Table("lab_rpki_sync_log_file").Select("id,filePath, fileName, fileType, syncType").
			Where("state->>'$.updateCertTable'=?", "notYet").And("syncLogId=?", 18).OrderBy("id").Find(&syncLogFiles)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(len(syncLogFiles))
	*/

	/*
			select filePath from lab_rpki_cer where fileName='eca135ce-a78f-437e-a32f-4714e9a373d2.cer'
		union
		select filePath from lab_rpki_crl where fileName='eca135ce-a78f-437e-a32f-4714e9a373d2.cer'
		union
		select filePath from lab_rpki_mft where fileName='eca135ce-a78f-437e-a32f-4714e9a373d2.cer'
		union
		select filePath from lab_rpki_mft where fileName='eca135ce-a78f-437e-a32f-4714e9a373d2.cer'
	*/
	sqls := `
	    select filePath from lab_rpki_cer where fileName='eca135ce-a78f-437e-a32f-4714e9a373d2.cer'
		union
		select filePath from lab_rpki_crl where fileName='eca135ce-a78f-437e-a32f-4714e9a373d2.cer'
		union
		select filePath from lab_rpki_mft where fileName='eca135ce-a78f-437e-a32f-4714e9a373d2.cer'
		union
		select filePath from lab_rpki_mft where fileName='eca135ce-a78f-437e-a32f-4714e9a373d2.cer'`
	var filePath string
	has, err := engine.SQL(sqls).Get(&filePath)
	fmt.Println(filePath, has, err)
}
