package main

import (
	"fmt"
	"os"
	"strings"
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

	fmt.Println(os.Args)
	if len(os.Args) > 1 {
		args := strings.Split(os.Args[1], "=")
		fmt.Println(args)
		if len(args) > 1 && (args[0] == "conf" || args[0] == "-conf" || args[0] == "--conf") {
			fmt.Println(args[1])
		}
	}

}
