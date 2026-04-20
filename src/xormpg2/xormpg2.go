package main

import (
	"fmt"

	"github.com/cpusoft/goutil/xormdb"
)

func main() {
	// start mysql/pg
	//err = xormdb.InitMySql()
	err := xormdb.InitPostgreSQL()
	xormdb.XormEngine.ShowSQL(true)
	defer xormdb.XormEngine.Close()

	session, err := xormdb.NewSession()
	if err != nil {
		fmt.Println("updateRsyncLogChainValidateStateStartDb(): NewSession fail:", err)
		return
	}
	defer session.Close()
	chainCerts := `{"id":41,"parentChainCers":[{"id":3593},{"id":3},{"id":1},{"id":7}]}`
	state := `{"state":"valid","errors":[],"warnings":[]}`
	originJson := `{"rir":"APNIC","repo":"rpki.apnic.net","notifyUrl":""}`
	roaId := 1111111
	sqlStr := `
			UPDATE lab_rpki_roa 
			SET 
				"chainCerts" = ?, 
				state = ?, 
				origin =  CASE
					WHEN origin = '{}' THEN ?
					WHEN origin = '""' THEN ?
					WHEN origin IS NULL THEN ?
					ELSE origin
				END 
			WHERE id = ?`
	_, err = session.Exec(sqlStr, chainCerts, state,
		originJson, originJson, originJson, roaId)
	if err != nil {
		fmt.Println("updateRoaDb(): UPDATE lab_rpki_roa fail :", roaId, "   chainCerts:", chainCerts,
			"   state:", state, "  originJson:", originJson, err)
		return
	}
	err = xormdb.CommitSession(session)
	if err != nil {
		return
	}
}
