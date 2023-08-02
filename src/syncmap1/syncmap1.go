package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cpusoft/goutil/jsonutil"
)

func main() {
	failUrls := jsonSyncMap{}
	snapshotFailUrls := make(map[string]string)
	snapshotFailUrls["1_https://rpki.telecentras.lt/"] = "1"
	snapshotFailUrls["2_https://rrdp.ripe.net/"] = "2"
	snapshotFailUrls["3_https://ca.rg.net"] = "3"
	snapshotFailUrls["4_https://google.com"] = "4"
	for k, v := range snapshotFailUrls {
		url := getUrlFromKey(k)
		failUrls.Store(url, v)
	}
	fmt.Println(failUrls)
	fmt.Println(jsonutil.MarshalJson(failUrls))
}

type jsonSyncMap struct {
	sync.Map
}

func (c jsonSyncMap) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	c.Range(func(key, value interface{}) bool {
		m[key.(string)] = value
		return true
	})
	return []byte(jsonutil.MarshalJson(m)), nil
}

func getUrlFromKey(key string) (url string) {
	if len(key) == 0 {
		return ""
	}
	split := strings.Split(key, "_")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}
