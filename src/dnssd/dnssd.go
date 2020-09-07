package main

import (
	"fmt"
	"net"
)

func main() {
	ips, _ := net.LookupIP("www.bing.com")
	fmt.Println(ips)

	cname, _ := net.LookupCNAME("www.bilibili.com")
	fmt.Println(cname)

	cname, srvs, _ := net.LookupSRV("xmpp-server", "tcp", "google.com")
	fmt.Println(cname)
	fmt.Println("优先级    权重 目标 端口")
	for _, value := range srvs {
		fmt.Printf("%d  %d  %s %d\n", value.Priority, value.Weight, value.Target, value.Port)
	}

}
