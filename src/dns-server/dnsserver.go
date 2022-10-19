package main

import (
	"fmt"

	dnsserver "labscm.zdns.cn/dns-mod/dns-server"
)

func main() {

	err := dnsserver.StartDnsServer(bool, 8888, 8889)
	fmt.Println(err)
}
