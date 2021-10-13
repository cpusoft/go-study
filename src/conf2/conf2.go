package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/cpusoft/goutil/osutil"
)

func main() {
	/*
			cannot use flag in init()
				flagFile := flag.String("conf", "", "")
				flag.Parse()
				fmt.Println("conf file is ", *flagFile, " from args")
				exists, err := osutil.IsExists(*flagFile)
				if err != nil || !exists {
					*flagFile = osutil.GetParentPath() + string(os.PathSeparator) + "conf" + string(os.PathSeparator) + "project.conf"
					fmt.Println("conf file is ", *flagFile, " default")
				}
		so ,use os.Args

	*/
	var err error
	var conf string
	if len(os.Args) > 1 {
		args := strings.Split(os.Args[1], "=")
		if len(args) > 0 && (args[0] == "conf" || args[0] == "-conf" || args[0] == "--conf") {
			conf = args[1]
		}
	}

	// decide by "conf" directory
	if conf == "" {
		path, err := osutil.GetCurrentOrParentAbsolutePath("conf")
		if err != nil {
			panic("found " + path + " failed, " + err.Error())

		}
		conf = path + string(os.PathSeparator) + "project.conf"

	}
	fmt.Println("conf file is ", conf)
	conf2, err := config.NewConfig("ini", conf)
	if err != nil {
		panic("load " + conf + " failed, " + err.Error())

	}
	value := conf2.String("mysql::server")
	fmt.Println(value)

	err = conf2.Set("db::name", "astaxie")
	fmt.Println(err)
	value = conf2.String("db::name")
	fmt.Println(value)

	err = conf2.SaveConfigFile(`F:\share\我的坚果云\Go\go-study\src\conf\project.conf`)
	fmt.Println(err)
}
