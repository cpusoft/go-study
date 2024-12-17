package main

import (
	"time"

	"github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/conf"
	"github.com/cpusoft/goutil/jsonutil"
	"github.com/cpusoft/goutil/xormdb"
	"xorm.io/xorm"
)

var XormEnginePg = &xorm.Engine{}

func main() {
	belogs.Info("main(): start")
	err := InitPostgreSQL()
	if err != nil {
		ingoreNoMySQLError := conf.String("mysql::ingoreNoPgError")
		if ingoreNoMySQLError == "true" {
			belogs.Error("main(): start InitPostgreSQL failed, but is still running:", err)
		} else {
			belogs.Error("main(): start InitPostgreSQL failed:", err)
			return
		}
	}
	defer XormEnginePg.Close()
	XormEnginePg.ShowSQL(true)

	//fz, err := getFzAgentCollectItemsDb()
	//	belogs.Error(fz, err)
	saveToDb()
}
func saveToDb() {
	sql := `insert into fz_agent_collect_item
	(code, attr_value, create_time, status) values
	(?,?,?,?)`
	session, err := NewPostgreSQLSession()
	if err != nil {
		belogs.Error("saveToDb(): NewSession fail: err:", err)
		return
	}
	defer session.Close()
	_, err = session.Exec(sql, "KSHYMJKRPKIAPNICCASTATUS",
		`{"KSHYMJKRPKIFZCACASTATUS":{"value":"running"}"}`, time.Now(), "COLLECTED")
	if err != nil {
		belogs.Error("saveCaProcessLogDb(): sql fail:", sql, err)
		xormdb.RollbackAndLogError(session, "sql fail:"+sql, err)
	}
	if err = session.Commit(); err != nil {
		xormdb.RollbackAndLogError(session, "saveToAgentDb():commit fail", err)
		return
	}
}
func getFzAgentCollectItemsDb() ([]FzAgentCollectItem, error) {
	fzAgentCollectItems := make([]FzAgentCollectItem, 0)
	sql := `select id,task_id,code,attr_value,create_time,status
			from fz_agent_collect_item order by id`
	err := XormEnginePg.SQL(sql).Find(&fzAgentCollectItems)
	if err != nil {
		belogs.Error("getFzAgentCollectItemsDb(): fz_agent_collect_item fail:", err)
		return nil, err
	}
	belogs.Debug("getFzAgentCollectItemsDb(): fzAgentCollectItems:", jsonutil.MarshalJson(fzAgentCollectItems))
	return fzAgentCollectItems, nil
}

type FzAgentCollectItem struct {
	TaskId     int       `json:"taskId" xorm:"task_id integer"` // --> task_id
	Code       string    `json:"code" xorm:"code varchar"`
	AttrValue  string    `json:"attrValue" xorm:"attr_value varchar"`
	CreateTime time.Time `json:"createTime" xorm:"create_time datetime"`
	Status     string    `json:"status" xorm:"status varchar"`
}

func InitPostgreSQL() (err error) {
	user := conf.String("postgresql::user")
	password := conf.String("postgresql::password")
	server := conf.String("postgresql::server")
	database := conf.String("postgresql::database")
	maxidleconns := conf.Int("postgresql::maxidleconns")
	maxopenconns := conf.Int("postgresql::maxopenconns")
	XormEnginePg, err = xormdb.InitPostgreSQLParameter(user, password, server, database, maxidleconns, maxopenconns)
	if err != nil {
		belogs.Error("InitPostgreSQL(): fail: ", err)
		return err
	}
	return nil
}

// get new session, and begin session
func NewPostgreSQLSession() (*xorm.Session, error) {
	// open mysql session
	session := XormEnginePg.NewSession()
	if err := session.Begin(); err != nil {
		return nil, xormdb.RollbackAndLogError(session, "session.Begin() XormEnginePg fail", err)
	}
	return session, nil
}
