package main

import (
	_ "database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/cpusoft/goutil/jsonutil"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type SyncLogFileModel struct {
	Id        uint64      `json:"id" xorm:"pk autoincr"`
	SyncLogId uint64      `json:"syncLogId" xorm:"syncLogId int"`
	FilePath  string      `json:"filePath" xorm:"filePath varchar(512)"`
	FileName  string      `json:"fileName" xorm:"fileName varchar(128)"`
	FileType  string      `json:"fileType" xorm:"fileType varchar(16)"`
	SyncType  string      `json:"syncType" xorm:"syncType varchar(16)"`
	CertModel interface{} `json:"-"`

	//cerId / mftId / roaId / crlId
	CertId uint64 `json:"certId" xorm:"certId int"`
}

func main() {

	syncLogFileModelCh := make(chan *SyncLogFileModel, runtime.NumCPU())
	endCh := make(chan bool, 1)
	go callParseValidate(syncLogFileModelCh, endCh)
	getFromDb(syncLogFileModelCh, endCh)
	select {}
}

func callParseValidate(syncLogFileModelCh chan *SyncLogFileModel, endCh chan bool) {
	for {
		select {
		case syncLogFileModel, ok := <-syncLogFileModelCh:
			fmt.Println("callParseValidate(): id:", syncLogFileModel.Id,
				"  fileType:", syncLogFileModel.FileType, "   syncType:", syncLogFileModel.SyncType, ok)
			if ok && syncLogFileModel.Id > 0 {
				parseValidate(syncLogFileModel)
			}

		case end := <-endCh:
			if end {
				fmt.Println("callParseValidate(): end:")
				return
			}
		}
	}
}
func parseValidate(syncLogFileModel *SyncLogFileModel) {

}
func getFromDb(syncLogFileModelCh chan *SyncLogFileModel, endCh chan bool) {

	defer func() {
		endCh <- true
	}()
	start := time.Now()
	//DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")
	user := "rpstir2"
	password := "Rpstir-123"
	server := "10.1.135.22:13306"
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
	defer engine.Close()
	//连接测试
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(maxidleconns)
	//设置最大打开连接数
	engine.SetMaxOpenConns(maxopenconns)
	engine.SetTableMapper(names.SnakeMapper{})
	engine.ShowSQL(true)

	labRpkiSyncLogId := 21
	dbSyncLogFileModel := new(SyncLogFileModel)
	sql := `select s.id,s.syncLogId,s.filePath,s.fileName, s.fileType, s.syncType, 
			cast(CONCAT(IFNULL(c.id,''),IFNULL(m.id,''),IFNULL(l.id,''),IFNULL(r.id,''),IFNULL(a.id,'')) as unsigned int) as certId from lab_rpki_sync_log_file s 
		left join lab_rpki_cer c on c.filePath = s.filePath and c.fileName = s.fileName  
		left join lab_rpki_mft m on m.filePath = s.filePath and m.fileName = s.fileName  
		left join lab_rpki_crl l on l.filePath = s.filePath and l.fileName = s.fileName  
		left join lab_rpki_roa r on r.filePath = s.filePath and r.fileName = s.fileName 
		left join lab_rpki_asa a on a.filePath = s.filePath and a.fileName = s.fileName 
		where s.syncLogId=? order by s.id `
	rows, err := engine.SQL(sql, labRpkiSyncLogId).Rows(dbSyncLogFileModel)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("rows:", rows, " time(s):", time.Since(start))

	defer rows.Close()
	var count int
	for rows.Next() {
		err = rows.Scan(dbSyncLogFileModel)
		fmt.Println("dbSyncLogFileModel:", jsonutil.MarshalJson(dbSyncLogFileModel), "  , time(s):", time.Since(start))
		syncLogFileModelCh <- dbSyncLogFileModel
		count++
	}
	fmt.Println("getFromDb end, count:", count, "  time(s):", time.Since(start))
}
