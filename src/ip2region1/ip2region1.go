package main

import (
	"fmt"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

func main() {
	var dbPath = `./ip2region.xdb`
	searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		fmt.Printf("failed to create searcher: %s\n", err.Error())
		return
	}

	defer searcher.Close()

	// do the search
	var ip = "114.242.3.2"
	var tStart = time.Now()
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
		return
	}

	fmt.Printf("{region: %s, took: %s}\n", region, time.Since(tStart))

	// 备注：并发使用，每个 goroutine 需要创建一个独立的 searcher 对象。
}
