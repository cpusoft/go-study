package main

import (
	"fmt"

	"github.com/cpusoft/goutil/httpclient"
)

func main() {
	url := `https://rdap.apnic.net/autnum/38082`
	_, body, err := httpclient.GetHttpsWithConfig(url, httpclient.NewHttpClientConfigWithParam(5, 3, "all", true))
	fmt.Println(body)
	fmt.Println(err)
}
