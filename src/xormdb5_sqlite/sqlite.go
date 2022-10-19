package main

import (
	"github.com/cpusoft/goutil/belogs"
	_ "github.com/cpusoft/goutil/conf"
	_ "github.com/cpusoft/goutil/logs"
	"github.com/cpusoft/goutil/xormdb"
	_ "github.com/mattn/go-sqlite3"
)

//set CGO_ENABLED=1

func main() {
	// start mysql
	err := xormdb.InitSqlite()
	if err != nil {
		belogs.Error("main(): start InitSqlite failed:", err)
		return
	}
	defer xormdb.XormEngine.Close()

	session, err := xormdb.NewSession()
	if err != nil {
		belogs.Error("main(): NewSession fail: err:", err)
		return
	}
	defer session.Close()
	var initSqls []string = []string{
		`CREATE TABLE IF NOT EXISTS lab_dns_origin(
		 id integer primary key autoincrement,
		 origin varchar(20) not null,
		 ttl int(10) not null,
		 updateTime text
		)`,
		`CREATE TABLE IF NOT EXISTS lab_dns_rr (
			id integer primary key autoincrement,
			originId int(10) not null, 
			rrName varchar(512) not null ,
			rrFullDomain varchar(512) not null ,
			rrType varchar(256) not null ,
			rrClass varchar(256),
			rrTtl int(10),
			rrData varchar(1024),
			updateTime text NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS rrFullDomain ON lab_dns_rr(rrFullDomain)`,
		`CREATE INDEX IF NOT EXISTS rrFullDomainAndRrType ON lab_dns_rr(rrFullDomain,rrType)`,
	}
	belogs.Debug("initSqls:", initSqls)
	for _, sql := range initSqls {
		_, err = session.Exec(sql)
		if err != nil {
			xormdb.RollbackAndLogError(session, "sql fail:"+sql, err)
			return
		}
	}

	insertSql := `insert into lab_dns_rr(	
		originId, rrName,rrFullDomain,
		rrType,	rrClass, rrTtl,
		rrData, updateTime) values (
		?,?,?,
		?,?,?,
		?,?)
		`
	_, err = session.Exec(insertSql,
		1, "example.com", "sqlite.example.com",
		"A", "IN", 1000,
		"1.1.1.1", "2022-10-19T16:51:00")
	if err != nil {
		xormdb.RollbackAndLogError(session, "insertSql fail:"+insertSql, err)
		return
	}

	err = xormdb.CommitSession(session)
	if err != nil {
		xormdb.RollbackAndLogError(session, "CommitSession fail", err)
		return
	}
	belogs.Info("main(): CommitSession ok")

}
