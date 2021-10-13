package main

import (
	"fmt"

	"github.com/cpusoft/goutil/belogs"
	_ "github.com/cpusoft/goutil/logs"
	"github.com/cpusoft/goutil/xormdb"
	sys "labscm.zdns.cn/rpstir2-mod/rpstir2-sys"
)

func main() {
	fmt.Println("start export")
	// start mysql
	err := xormdb.InitMySql()
	if err != nil {
		belogs.Error("main(): start InitMySql failed:", err)
		fmt.Println("failed to start, ", err)
		return
	}
	defer xormdb.XormEngine.Close()

	err = sys.ExportRtrForManrsConsole()
	if err != nil {
		fmt.Println("fail to export:", err)
		return
	}
	fmt.Println("success to export")
}
