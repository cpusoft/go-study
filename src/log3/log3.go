package main

import (
	belogs "github.com/beego/beego/v2/core/logs"
)

func main() {
	Init()
	belogs.Debug("ssss")
	belogs.Warn(map[string]int{"key": 2016, "key2": 2017})

}
