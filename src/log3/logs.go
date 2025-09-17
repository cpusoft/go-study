package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	logs "github.com/beego/beego/v2/core/logs"
	"github.com/cpusoft/goutil/osutil"
)

/*
LevelEmergency = iota
LevelAlert
LevelCritical
LevelError
LevelWarning
LevelNotice
LevelInformational
LevelDebug
*/
func Init() {

	logLevel := "LevelDebug" // conf.String("logs::level")
	// get process file name as log name
	logName := filepath.Base(os.Args[0])
	if logName != "" {
		logName = strings.Split(logName, ".")[0] + ".log"
	} else {
		logName = `D:\share\我的坚果云\Code\common\go-study\src\log3\test.log` //conf.String("logs::name")
	}
	fmt.Println("logName", logName)
	async := false // conf.DefaultBool("logs::async", false)
	//fmt.Println("log", logLevel, logName)

	var logLevelInt int = logs.LevelInformational
	switch logLevel {
	case "LevelEmergency":
		logLevelInt = logs.LevelEmergency
	case "LevelAlert":
		logLevelInt = logs.LevelAlert
	case "LevelCritical":
		logLevelInt = logs.LevelCritical
	case "LevelError":
		logLevelInt = logs.LevelError
	case "LevelWarning":
		logLevelInt = logs.LevelWarning
	case "LevelNotice":
		logLevelInt = logs.LevelNotice
	case "LevelInformational":
		logLevelInt = logs.LevelInformational
	case "LevelDebug":
		logLevelInt = logs.LevelDebug
	}
	fmt.Println("logLevelInt", logLevelInt)
	//ts := time.Now().Format("2006-01-02")

	path, err := osutil.GetCurrentOrParentAbsolutePath("log")
	if err != nil {
		fmt.Println("found " + path + " failed, " + err.Error())
	}
	filePath := path + string(os.PathSeparator) + logName
	fmt.Println("log file is ", filePath)

	logConfig := make(map[string]interface{})
	logConfig["daily"] = true
	logConfig["hourly"] = false
	logConfig["filename"] = filePath // + "." + ts
	logConfig["maxlines"] = 0
	logConfig["maxfiles"] = 0
	logConfig["maxsize"] = 0
	logConfig["maxdays"] = 30
	logConfig["maxhours"] = 0
	logConfig["formatter"] = "json"
	logConfig["level"] = logLevelInt

	logConfigStr, _ := json.Marshal(logConfig)
	fmt.Println("log:logConfigStr", string(logConfigStr))
	//logs.NewLogger(1000000)
	err = logs.SetLogger(logs.AdapterFile, string(logConfigStr))
	if err != nil {
		fmt.Println(filePath + " SetLogger failed, " + err.Error() + ",   " + string(logConfigStr))
	}
	if async {
		logs.Async()
	}
	fmt.Println("log init ok, log file is ", filePath)
}
