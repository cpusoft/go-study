package main

import (
	"fmt"

	"github.com/cpusoft/goutil/xormdb"
)

func main() {
	//err = xormdb.InitMySql()
	err := xormdb.InitPostgreSQL()
	xormdb.XormEngine.ShowSQL(true)
	defer xormdb.XormEngine.Close()

	serialNumbers := make([]int32, 0)
	serialNumbers = append(serialNumbers, 10929, 10930, 10932)
	rtrIncrementals := make([]LabRpkiRtrIncremental, 0)
	err = xormdb.XormEngine.
		In("serialNumber", serialNumbers).
		OrderBy(`"serialNumber" ASC`).Find(&rtrIncrementals)
	if err != nil {
		fmt.Println("getRecentRtrIncrementalsDb():get rtrIncrementals fail",
			"serialNumbers", serialNumbers, "err", err)
		return
	}
	fmt.Println("rtrIncrementals:", rtrIncrementals)
}

// lab_rpki_rtr_incremental
type LabRpkiRtrIncremental struct {
	Id           uint64 `json:"id" xorm:"id bigint"`
	SerialNumber uint64 `json:"serialNumber" xorm:"serialNumber bigint"`
	//announce/withdraw, is 1/0 in protocol
	Style string `json:"style" xorm:"style varchar(16)"`
	Asn   int64  `json:"asn" xorm:"asn bigint"`
	//address: 63.60.00.00
	Address      string `json:"address" xorm:"address varchar(512)"`
	PrefixLength uint64 `json:"prefixLength" xorm:"prefixLength int"`
	MaxLength    uint64 `json:"maxLength" xorm:"maxLength int"`
	//'come from : {souce:sync/slurm/transfer,syncLogId/syncLogFileId/slurmId/slurmFileId/transferLogId}',
	SourceFrom string `json:"sourceFrom" xorm:"sourceFrom json"`
}
