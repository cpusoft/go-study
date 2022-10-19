package main

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var (
	Engine *xorm.Engine
)

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite3", `F:\share\我的坚果云\Go\go-study\src\sqlite1\test.db`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("create ok")

	err = Engine.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("connect ok")

}

var intiSqls []string = []string{
	`drop table if exists lab_rpki_cer`,
	`CREATE TABLE lab_rpki_cer (
		id  INTEGER PRIMARY KEY  AUTOINCREMENT,
		sn varchar(128) NOT NULL,
		notBefore datetime NOT NULL,
		notAfter datetime NOT NULL,
		subject varchar(512) ,
		issuer varchar(512) ,
		ski varchar(128) ,
		aki varchar(128) ,
		filePath varchar(512) NOT NULL ,
		fileName varchar(128) NOT NULL ,
		state varchar(1024),
		jsonAll varchar(2048),
		chainCerts varchar(1024),
		syncLogId int(10) not null,
		syncLogFileId int(10) not null,
		updateTime datetime NOT NULL,
		fileHash varchar(512) NOT NULL ,
		origin varchar(1024)
	) 
	`,
	`CREATE INDEX lab_rpki_cer_ski on lab_rpki_cer (ski)`,
	`CREATE INDEX lab_rpki_cer_aki on lab_rpki_cer (aki)`,
	`CREATE INDEX lab_rpki_cer_filePath on lab_rpki_cer (filePath)`,
	`CREATE INDEX lab_rpki_cer_fileName on lab_rpki_cer (fileName)`,
	`CREATE INDEX lab_rpki_cer_syncLogId on lab_rpki_cer (syncLogId)`,
	`CREATE INDEX lab_rpki_cer_syncLogFileId on lab_rpki_cer (syncLogFileId)`,
	`CREATE UNIQUE INDEX lab_rpki_cer_cerFilePathFileName on lab_rpki_cer (filePath,fileName)`,
	`CREATE UNIQUE INDEX lab_rpki_cer_cerSkiFilePathon on lab_rpki_cer (ski,filePath)`,
}

func createTable() {
	session := Engine.NewSession()
	defer session.Close()
	for _, sq := range intiSqls {
		sqlTime := time.Now()
		if _, err := session.Exec(sq); err != nil {
			fmt.Println("initResetImplDb():  "+sq+" fail", err)
			session.Rollback()
			return
		}
		fmt.Println("initResetImplDb(): sq:", sq, ", sql time(s):", time.Now().Sub(sqlTime).Seconds())
	}
	err := session.Commit()
	if err != nil {
		session.Rollback()
	}
}

func insertTable() (int64, error) {
	now := time.Now()
	session := Engine.NewSession()
	defer session.Close()
	sqlStr := `INSERT into lab_rpki_cer (
			sn, notBefore,notAfter,subject,
			issuer,ski,aki,filePath,fileName,
			fileHash,jsonAll,syncLogId,syncLogFileId,updateTime,
			state) 	
			VALUES 
			(?,?,?,?,
			?,?,?,?,?,
			?,?,?,?,?,
			?)`
	res, err := session.Exec(sqlStr,
		"sn", now, now, "subject",
		"issuer", "ski", "aki", "filePath", "fileName",
		"fileHash", "cerModel", 1, 2, now,
		"stateModel")
	if err != nil {
		fmt.Println("insertTable(): INSERT lab_rpki_cer Exec:", err)
		session.Rollback()
		return 0, err
	}

	cerId, err := res.LastInsertId()
	if err != nil {
		session.Rollback()
		fmt.Println("insertTable(): LastInsertId:", err)
		return 0, err
	}
	fmt.Println("insertTable(): cerId:", cerId)
	err = session.Commit()
	if err != nil {
		session.Rollback()
	}
	return cerId, nil
}
