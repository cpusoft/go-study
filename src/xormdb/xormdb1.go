package main

import (
	_ "database/sql"
	"fmt"

	belogs "github.com/astaxie/beego/logs"

	"github.com/cpusoft/goutil/jsonutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

func main() {
	//DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")
	user := "rpstir2"
	password := "Rpstir-123"
	server := "202.173.14.105:13306"
	database := "rpstir2"
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
	engine.ShowSQL(true)

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
		// get current , the set new value
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
	/*
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
	*/
	/*
		type ChainIpAddress struct {
			Id         uint64 `json:"id" xorm:"id int"`
			RangeStart string `json:"rangeStart" xorm:"rangeStart"`
			//max address range from addressPrefix or min/max, in hex:  63.69.7f.ff'
			RangeEnd string `json:"rangeEnd" xorm:"rangeEnd"`
		}
	*/

	/*
		type ChainIpAddress struct {
			Id            uint64 `json:"id" xorm:"id int"`
			AddressFamily uint64 `json:"-"  xorm:"addressFamily int"`
			//address prefix: 147.28.83.0/24 '
			AddressPrefix string `json:"-"  xorm:"addressPrefix varchar(512)"`
			MaxLength     uint64 `json:"-"  xorm:"maxLength int"`
			//min address:  99.96.0.0
			Min string `json:"-" xorm:"min varchar(512)`
			//max address:   99.105.127.255
			Max string `json:"-" xorm:"max varchar(512)`
			//min address range from addressPrefix or min/max, in hex:  63.60.00.00'
			RangeStart string `json:"rangeStart" xorm:"rangeStart"`
			//max address range from addressPrefix or min/max, in hex:  63.69.7f.ff'
			RangeEnd string `json:"rangeEnd" xorm:"rangeEnd"`
			//min--max, such as 192.0.2.0--192.0.2.130, will convert to addressprefix range in json:{192.0.2.0/25, 192.0.2.128/31, 192.0.2.130/32}
			AddressPrefixRange string `json:"-" xorm:"addressPrefixRange json"`
		}
		chainIpAddresses := make([]ChainIpAddress, 0)
		cerId := 24217
		err = engine.Table("lab_rpki_cer_ipaddress").
			Select("id,rangeStart,rangeEnd").
			Where("cerId=?", cerId).
			OrderBy("id").Find(&chainIpAddresses)
		fmt.Println(chainIpAddresses, err)
	*/

	/*
		type ChainRoa struct {
			Id          uint64 `json:"id" xorm:"id int"`
			Asn         uint64 `json:"asn" xorm:"asn int"`
			FilePath    string `json:"-" xorm:"filePath varchar(512)"`
			FileName    string `json:"-" xorm:"fileName varchar(128)"`
			Ski         string `json:"-" xorm:"ski varchar(128)"`
			Aki         string `json:"-" xorm:"aki varchar(128)"`
			State       string `json:"-" xorm:"state json"`
			EeCertStart uint64 `json:"-" xorm:"eeCertStart int"`
			EeCertEnd   uint64 `json:"-" xorm:"eeCertEnd int"`
		}
		chainRoa := ChainRoa{}
		roaId := 61591
		_, err = engine.Table("lab_rpki_roa").
			Select("id,asn,ski,aki,filePath,fileName,state,jsonAll->'$.eeCertModel.eeCertStart' as eeCertStart,jsonAll->'$.eeCertModel.eeCertEnd' as eeCertEnd").
			Where("id=?", roaId).Get(&chainRoa)
		fmt.Println(chainRoa, err)


	*/
	/*
		type LabRpkiRtrFullLog struct {
			Id           uint64 `json:"id" xorm:"id int"`
			SerialNumber uint64 `json:"serialNumber" xorm:"serialNumber bigint"`
			Asn          uint64 `json:"asn" xorm:"asn int"`
			//address: 63.60.00.00
			Address      string `json:"address" xorm:"address varchar(512)"`
			PrefixLength uint64 `json:"prefixLength" xorm:"prefixLength int"`
			MaxLength    uint64 `json:"maxLength" xorm:"maxLength int"`
			//'come from : {souce:sync/slurm/transfer,syncLogId/syncLogFileId/slurmId/slurmFileId/transferLogId}',
			SourceFrom string `json:"sourceFrom" xorm:"sourceFrom json"`
		}

		labRpkiRtrFullLog := LabRpkiRtrFullLog{}
		Asn := uint64(1)
		PrefixLength := uint64(1)
		MaxLength := uint64(0)
		Address := "1.1.1.1"
		if Asn > 0 {
			labRpkiRtrFullLog.Asn = Asn

		}
		if PrefixLength > 0 {
			labRpkiRtrFullLog.PrefixLength = PrefixLength

		}
		if MaxLength > 0 {
			labRpkiRtrFullLog.MaxLength = MaxLength

		}
		if len(Address) > 0 {
			labRpkiRtrFullLog.Address = Address

		}

		fmt.Println(jsonutil.MarshalJson(labRpkiRtrFullLog))
		aff, err := session.Delete(&labRpkiRtrFullLog)
		fmt.Println(aff, err)

	*/
	/*
		// /root/rpki/data/repo/rpki.arin.net/repository/arin-rpki-ta/5e4a23ea-e80a-403e-b08c-2171da2157d3/f60c9f32-a87c-4339-a2f3-6299a3b02e29/
		filePathPrefix := `/root/rpki/data/repo/rpki.arin.net/`
		cerId := make([]uint64, 0, 1000)
		//err = session.SQL("select id from lab_rpki_cer Where filePath like '" + filePathPrefix + "%'").Find(&cerId)
		err = session.SQL("select id from lab_rpki_cer Where filePath like ? ", filePathPrefix+"%").Find(&cerId)
		fmt.Println(jsonutil.MarshalJson(cerId), err)

	*/
	/*
		var mftId uint32
		filePath := `/root/rpki/data/reporrdp/rpki.ripe.net/repository/DEFAULT/`
		fileName := `KpSo3VVK5wEHIJnHC2QHVV3d5mk.mft`
		has, err := engine.Table("lab_rpki_mft").Where("filePath=?", filePath).And("fileName=?", fileName).Cols("id").Get(&mftId)
		if err != nil {
			fmt.Println("GetChainMftId(): lab_rpki_mft id fail, filePath, fileName:", filePath, fileName, err)
			return
		}
		if !has {
			fmt.Println("GetChainMftId(): lab_rpki_mft id has not found, filePath, fileName:", filePath, fileName, err)
			return
		}
		fmt.Println("GetChainMftId():", mftId)
	*/
	/*
		type ChainFileHash struct {
			File string `json:"-" xorm:"file varchar(1024)"`
			Hash string `json:"-" xorm:"hash varchar(1024)"`
		}
		chainFileHashs := make([]ChainFileHash, 0)
		mftId := 31382
		err = engine.Table("lab_rpki_mft_file_hash").
			Cols("file,hash").
			Where("mftId=?", mftId).
			OrderBy("id").Find(&chainFileHashs)
		if err != nil {
			fmt.Println("getChainFileHashs(): lab_rpki_mft_file_hash fail:", err)
			return
		}
		fmt.Println("getChainFileHashs():mftId, len(chainFileHashs):",
			mftId, jsonutil.MarshalJson(chainFileHashs), len(chainFileHashs))
	*/
	/*
		cerId := 10
		var tmpIds []int64
		err = session.Table("lab_rpki_cer_sia").Where("cerId=?", cerId).Cols("id").Find(&tmpIds)
		if err != nil {
			fmt.Println("lab_rpki_cer_sia find fail:", err)
			return
		}
		fmt.Println(tmpIds)
		if len(tmpIds) == 0 {
			return
		}

		ids, err := session.In("id", tmpIds).Delete("lab_rpki_cer_sia")
		if err != nil {
			fmt.Println("lab_rpki_cer_sia delete fail:", err)
			return
		}
		fmt.Println("lab_rpki_cer_sia :", ids)
		err = session.Rollback()
		if err != nil {
			fmt.Println("lab_rpki_cer_sia rollback fail:", err)
			return
		}
	*/
	/*
		type ChainMft struct {
			Id          uint64 `json:"id" xorm:"id int"`
			FilePath    string `json:"-" xorm:"filePath varchar(512)"`
			FileName    string `json:"-" xorm:"fileName varchar(128)"`
			Ski         string `json:"-" xorm:"ski varchar(128)"`
			Aki         string `json:"-" xorm:"aki varchar(128)"`
			MftNumber   string `json:"-" xorm:"mftNumber varchar(1024)"`
			State       string `json:"-" xorm:"state json"`
			EeCertStart uint64 `json:"-" xorm:"eeCertStart int"`
			EeCertEnd   uint64 `json:"-" xorm:"eeCertEnd int"`
		}
		var mftId uint64
		mftId = 1
		chainMft := ChainMft{}
		_, err = session.Table("lab_rpki_mft").
			Select("id,ski,aki,filePath,fileName,mftNumber,state,jsonAll->'$.eeCertModel.eeCertStart' as eeCertStart,jsonAll->'$.eeCertModel.eeCertEnd' as eeCertEnd").
			Where("id=?", mftId).Get(&chainMft)
		if err != nil {
			fmt.Println("GetChainMft(): lab_rpki_mft fail:", mftId, err)
			return
		}
		fmt.Println(chainMft)
	*/
	/*
		type ChainSnInCrlRevoked struct {
			CrlFileName    string    `json:"-" xorm:"fileName varchar(512)"`
			RevocationTime time.Time `json:"-" xorm:"revocationTime datetime"`
		}
		cerId := 100
		chainSnInCrlRevoked := ChainSnInCrlRevoked{}
		sql := `select l.fileName, r.revocationTime from lab_rpki_cer c, lab_rpki_crl l, lab_rpki_crl_revoked_cert r
		 where  c.sn = r.sn and r.crlId = l.id and c.aki = l.aki and c.id=` + convert.ToString(cerId)
		has, err := engine.
			Sql(sql).Get(&chainSnInCrlRevoked)
		if err != nil {
			fmt.Println("select fail:", has, err)
			return
		}
		fmt.Println(chainSnInCrlRevoked)

	*/
	type SyncLogFileModel struct {
		Id        uint64 `json:"id" xorm:"pk autoincr"`
		SyncLogId uint64 `json:"syncLogId" xorm:"syncLogId int"`
		FilePath  string `json:"filePath" xorm:"filePath varchar(512)"`
		FileName  string `json:"fileName" xorm:"fileName varchar(128)"`
		FileType  string `json:"fileType" xorm:"fileType varchar(16)"`
		SyncType  string `json:"syncType" xorm:"syncType varchar(16)"`
		JsonAll   string `json:"jsonAll"`
		//cerId / mftId / roaId / crlId
		CertId uint64 `json:"certId"`
	}
	labRpkiSyncLogId := 8
	dbSyncLogFileModels := make([]SyncLogFileModel, 0)
	err = engine.Table("lab_rpki_sync_log_file").Select("id,syncLogId,filePath,fileName, fileType, syncType").
		Where("state->'$.updateCertTable'=?", "notYet").And("syncLogId=?", labRpkiSyncLogId).
		And("fileType=?", "mft").
		OrderBy("id").Find(&dbSyncLogFileModels)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, _ := range dbSyncLogFileModels {
		if dbSyncLogFileModels[i].FileName == "fIPsfy5eLDPb2Ki9g3-aa8fzomM.mft" {
			fmt.Println(jsonutil.MarshalJson(dbSyncLogFileModels[i]))
		}
	}

}
