package main

import (
	"fmt"
	"net"
)

func main() {

	ipAddress, IPnet, err := net.ParseCIDR("198.162.0.0/16")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("IP address :", ipAddress)
	fmt.Printf("IP Network : %#v\n ", IPnet)
	fmt.Println("IPnet.IP : ", IPnet.IP)
	fmt.Println("Contains 192.162.0.0 : ", IPnet.Contains(net.ParseIP("192.162.0.0")))
	fmt.Println("Network : ", IPnet.Network())

	ipAddress, IPnet, err = net.ParseCIDR("2001:DB6::/32")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("IP address :", ipAddress)
	fmt.Printf("IP Network : %#v\n ", IPnet)
	fmt.Println("IPnet.IP : ", IPnet.IP)
	fmt.Println("Contains 192.162.0.0 : ", IPnet.Contains(net.ParseIP("192.162.0.0")))
	fmt.Println("Network : ", IPnet.Network())
}
