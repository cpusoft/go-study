package main

import (
	"fmt"
	"github.com/zpatrick/go-config"
	//	"http://gopkg.in/gcfg.v1"
)

func main() {
	iniFile := config.NewINIFile("E:\\Go\\test1\\src\\main\\slurm.conf")
	c := config.NewConfig([]config.Provider{iniFile})
	if err := c.Load(); err != nil {
		fmt.Println(err)
		return
	}

	host, err2 := c.String("db.host")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println(host)
}
