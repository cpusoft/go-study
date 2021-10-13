package main

import (
	"flag"
	"fmt"
	"log"

	belogs "github.com/beego/beego/v2/core/logs"
	glog "github.com/golang/glog"

	//	"log/syslog"
	"encoding/json"
	"os"
	"path/filepath"

	config "github.com/beego/beego/v2/core/config"
)

func main() {
	arr := []int{2, 3}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {

	}
	fmt.Println("dir:", dir)
	pwd, _ := os.Getwd()
	fmt.Println("pwd:", pwd)

	//log
	log.SetFlags(log.Lshortfile | log.LUTC)
	log.Print("Print array ", arr, "\n")

	//golang glog
	flag.Parse()
	defer glog.Flush()

	flag.Set("log_dir", "E:\\Go\\test1\\log")
	flag.Set("dailyRolling", "true")
	p, err := os.Getwd()
	if err != nil {
		glog.Info("Getwd: ", err)
	} else {
		glog.Info("Getwd: ", p)
	}
	glog.Info("This is a Info log", arr)

	//beego logs
	config1 := make(map[string]interface{})
	config1["filename"] = "E:\\Go\\test1\\log\\beegologs.log"
	config1["level"] = belogs.LevelDebug
	configStr, err := json.Marshal(config1)
	fmt.Println("config1:", config1)
	//belogs.SetLogger(belogs.AdapterFile, string(configStr))
	//belogs.Debug("this is a test, my name is %s", "stu01")

	//beego logs , config is in ini file
	conf, err := config.NewConfig("ini", "E:\\Go\\test1\\src\\main\\slurm.conf")
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	log_level := conf.String("logs::log_level")
	fmt.Println("log_level:", log_level)

	if len(log_level) == 0 {
		log_level = "debug"
	}
	log_path := conf.String("logs::log_path")
	fmt.Println("log_path:", log_path)
	fmt.Println("log_path:", pwd+log_path)

	config1["filename"] = pwd + log_path
	config1["level"] = belogs.LevelDebug
	fmt.Println("config1:", config1)
	configStr, err = json.Marshal(config1)
	belogs.SetLogger(belogs.AdapterFile, string(configStr))
	belogs.Debug("new log, my name is %s", "stu01")
	belogs.Info("sssss")

}
