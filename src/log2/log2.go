package main

func main() {
	//	belogs.Info("ssss")
}

/*
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug

func init() {

	logLevel := conf.String("logs::level")
	// get process file name as log name
	logName := filepath.Base(os.Args[0])
	if logName != "" {
		logName = strings.Split(logName, ".")[0] + ".log"
	} else {
		logName = conf.String("logs::name")
	}
	async := conf.DefaultBool("logs::async", false)
	//fmt.Println("log", logLevel, logName)

	var logLevelInt int = belogs.LevelInformational
	switch logLevel {
	case "LevelEmergency":
		logLevelInt = belogs.LevelEmergency
	case "LevelAlert":
		logLevelInt = belogs.LevelAlert
	case "LevelCritical":
		logLevelInt = belogs.LevelCritical
	case "LevelError":
		logLevelInt = belogs.LevelError
	case "LevelWarning":
		logLevelInt = belogs.LevelWarning
	case "LevelNotice":
		logLevelInt = belogs.LevelNotice
	case "LevelInformational":
		logLevelInt = belogs.LevelInformational
	case "LevelDebug":
		logLevelInt = belogs.LevelDebug
	}
	//ts := time.Now().Format("2006-01-02")

	//
	path, err := osutil.GetCurrentOrParentAbsolutePath("log")
	if err != nil {
		panic("found " + path + " failed, " + err.Error())
	}
	log := path + string(os.PathSeparator) + logName
	fmt.Println("log file is ", log)

	logConfig := make(map[string]interface{})
	logConfig["filename"] = log // + "." + ts
	logConfig["level"] = logLevelInt
	// no max lines
	logConfig["maxlines"] = 0
	logConfig["maxsize"] = 0
	logConfig["daily"] = true
	logConfig["maxdays"] = 30

	logConfigStr, _ := json.Marshal(logConfig)
	//fmt.Println("log:logConfigStr", string(logConfigStr))
	belogs.NewLogger(1000000)
	belogs.SetLogger(belogs.AdapterFile, string(logConfigStr))
	if async {
		belogs.Async()
	}

}
*/
