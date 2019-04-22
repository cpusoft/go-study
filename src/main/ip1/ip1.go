package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
)

func main() {
	//https://stackoverflow.com/questions/19882961/go-golang-check-ip-address-in-range
	network := "192.168.5.0/24"
	clientips := []string{
		"192.168.5.1",
		"192.168.6.0",
	}
	_, subnet, _ := net.ParseCIDR(network)
	fmt.Println("subnet:", subnet)
	for _, clientip := range clientips {
		ip := net.ParseIP(clientip)
		fmt.Println("clientip:", ip)
		if subnet.Contains(ip) {
			fmt.Println("IP in subnet", clientip)
		}
	}
	ipp := net.ParseIP(clientips[0])
	i, err := IP2Integer(&ipp)
	fmt.Println(i, err)

	ipv4Decimal := IP4toInt(net.ParseIP(clientips[0]))
	fmt.Println(ipv4Decimal)

	mask := uint64(0xFFFFFFFF<<(32-24)) & 0xFFFFFFFF //24 is the netmask
	fmt.Printf("mask:%x\n", mask)
	var dmask uint64
	dmask = 32
	fmt.Printf("dmask:%x\n", dmask)
	localmask := make([]uint64, 0, 4)
	fmt.Println("localmask:", localmask)
	for i := 1; i <= 4; i++ {
		tmp := mask >> (dmask - 8) & 0xFF
		localmask = append(localmask, tmp)
		dmask -= 8
	}
	fmt.Println("==================================")

	ip := net.ParseIP("192.168.5.1")
	fmt.Println("ip:", ip, "    len(ip)", len(ip))
	ipI, _ := IPToString(ip)
	fmt.Println("ipI:", ipI)

	network = "192.168.5.0/24"
	_, subnet, _ = net.ParseCIDR(network)
	fmt.Println("subnet.Mask:", subnet.Mask)
	min := make(net.IP, net.IPv4len)
	max := make(net.IP, net.IPv4len)
	for i := 0; i < net.IPv4len; i++ {
		min[i] = subnet.IP[i] & subnet.Mask[i]
		max[i] = subnet.IP[i] | (^subnet.Mask[i])
	}
	fmt.Println("min:", min, " max:", max)
	minI, _ := IPToString(min)
	maxI, _ := IPToString(max)
	fmt.Println("minI:", minI, " maxI:", maxI)

	network = "2803:d380::/28"
	_, subnet, _ = net.ParseCIDR(network)
	fmt.Println("subnet.Mask:", subnet.Mask)
	min = make(net.IP, net.IPv6len)
	max = make(net.IP, net.IPv6len)
	for i := 0; i < net.IPv6len; i++ {
		min[i] = subnet.IP[i] & subnet.Mask[i]
		max[i] = subnet.IP[i] | (^subnet.Mask[i])
	}
	fmt.Println("min:", min, " max:", max)
	minI, _ = IPToString(min)
	maxI, _ = IPToString(max)
	fmt.Println("minI:", minI, " maxI:", maxI)
}

func IP2Integer(ip *net.IP) (int64, error) {
	ip4 := ip.To4()
	if ip4 == nil {
		return 0, fmt.Errorf("illegal: %v", ip)
	}

	bin := make([]string, len(ip4))
	for i, v := range ip4 {
		bin[i] = fmt.Sprintf("%08b", v)
	}
	return strconv.ParseInt(strings.Join(bin, ""), 2, 64)
}

func IPToString(ip net.IP) (string, error) {

	var buffer bytes.Buffer
	if len(ip) == net.IPv4len {
		for i := 0; i < net.IPv4len; i++ {
			if i < net.IPv4len-1 {
				buffer.WriteString(fmt.Sprintf("%02x.", ip[i]))
			} else {
				buffer.WriteString(fmt.Sprintf("%02x", ip[i]))
			}
		}
		return buffer.String(), nil
	} else if len(ip) == net.IPv6len {
		for i := 0; i < net.IPv6len; i++ {
			if i < net.IPv6len-1 {
				buffer.WriteString(fmt.Sprintf("%02x:", ip[i]))
			} else {
				buffer.WriteString(fmt.Sprintf("%02x", ip[i]))
			}
		}
		return buffer.String(), nil
	}
	return "", errors.New("illegal ip type")

}

//https://www.socketloop.com/tutorials/golang-convert-ipv4-address-to-decimal-number-base-10-or-integer
func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}

func IP16toInt(IPv6Address net.IP) int64 {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(IPv6Address.To16())
	return IPv6Int.Int64()
}
