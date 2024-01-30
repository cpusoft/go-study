package main

import (
	"fmt"

	"github.com/cpusoft/goutil/httpclient"
)

func main() {
	url := `https://rdap.apnic.net/autnum/38082`
	_, body, err := httpclient.GetHttpsVerify(url, true)
	fmt.Println(body)
	fmt.Println(err)
}
