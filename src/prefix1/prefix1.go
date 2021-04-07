package main

import (
	"fmt"
	"net"

	"github.com/cpusoft/goutil/iputil"
)

func main() {
	_, ipNet1, err := net.ParseCIDR("204.2.135.0/24")
	fmt.Println(ipNet1, err)

	_, ipNet2, err := net.ParseCIDR("204.2.226.240/28")
	fmt.Println(ipNet2, err)

	ipPrefix1 := iputil.Prefix{*ipNet1}
	fmt.Println(ipPrefix1)
	ipPrefix2 := iputil.Prefix{*ipNet2}
	fmt.Println(ipPrefix2)

	c1 := ipPrefix1.Contains(&ipPrefix2)
	fmt.Println(c1)

	c2 := ipPrefix2.Contains(&ipPrefix1)
	fmt.Println(c2)

}
