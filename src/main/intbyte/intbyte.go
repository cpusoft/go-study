package main

import (
	"bytes"

	"fmt"
	"strconv"
	"strings"
)

func main() {
	var a []byte = []byte{80, 128, 0, 0}
	fmt.Println(a)

	aaa := fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
	fmt.Println(aaa)

	ip := "199.99.99.0"
	ss := IpToRtrFormat(ip)
	fmt.Println(ss)
	s := []byte(ss)
	bb0, _ := strconv.ParseInt(string(s[0:2]), 16, 0)
	bb1, _ := strconv.ParseInt(string(s[2:4]), 16, 0)
	bb2, _ := strconv.ParseInt(string(s[4:6]), 16, 0)
	bb3, _ := strconv.ParseInt(string(s[6:8]), 16, 0)
	fmt.Println(bb0, bb1, bb2, bb3)
	//200107F8001900000000000000000000
	c := []byte{32, 1, 6, 120, 0, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	fmt.Println(c)

	cc0 := fmt.Sprintf("%0x", c[0:2])
	fmt.Println(cc0)

	cc0 = fmt.Sprintf("%s", c[0:2])
	fmt.Println(cc0)

	str := "2001:DB8::"
	di := IpToRtrFormat(str)
	fmt.Println(di)
	ips := di
	fmt.Println(len(ips))
	ip0 := string(ips[0:8])
	ip1 := string(ips[8:16])
	ip2 := string(ips[16:24])
	ip3 := string(ips[24:32])
	fmt.Println("%s:%s:%s:%s", ip0, ip1, ip2, ip3)

	var buffer bytes.Buffer
	buffer.WriteString(ips[0:4])
	buffer.WriteString(":")
	buffer.WriteString(ips[4:8])
	buffer.WriteString(":")
	buffer.WriteString(ips[8:12])
	buffer.WriteString(":")
	buffer.WriteString(ips[12:16])
	buffer.WriteString(":")
	buffer.WriteString(ips[16:20])
	buffer.WriteString(":")
	buffer.WriteString(ips[20:24])
	buffer.WriteString(":")
	buffer.WriteString(ips[24:28])
	buffer.WriteString(":")
	buffer.WriteString(ips[28:32])
	fmt.Println(buffer.String())
}
func IpToRtrFormat(ip string) string {
	formatIp := ""

	// format  ipv4
	ipsV4 := strings.Split(ip, ".")
	if len(ipsV4) > 1 {
		for _, ipV4 := range ipsV4 {
			ip, _ := strconv.Atoi(ipV4)
			formatIp += fmt.Sprintf("%02x", ip)
		}
		return formatIp
	}

	// format ipv6
	count := strings.Count(ip, ":")
	if count > 0 {
		count := strings.Count(ip, ":")
		if count < 7 { // total colon is 8
			needCount := 7 - count + 2 //2 is current "::", need add
			colon := strings.Repeat(":", needCount)
			ip = strings.Replace(ip, "::", colon, -1)

		}
		ipsV6 := strings.Split(ip, ":")

		for _, ipV6 := range ipsV6 {
			formatIp += fmt.Sprintf("%04s", ipV6)
		}
		return formatIp
	}
	return ""
}
