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

	ip := net.ParseIP("192.168.5.2").To4()

	fmt.Println("ip:", ip, "    len(ip)", len(ip))
	ipI, _ := IPToString(ip)
	fmt.Println("ipI:", ipI)

	network := "192.168.5/24"
	network, _ = FillIP(network, 1)
	fmt.Println("network:", network)
	_, subnet, _ := net.ParseCIDR(network)
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

	if minI <= ipI && ipI <= maxI {
		fmt.Println("ipI:", ipI, "  in minI:", minI, " maxI:", maxI)
	} else {
		fmt.Println("ipI:", ipI, " not in minI:", minI, " maxI:", maxI)
	}

	network = "2803:d380/28"
	network, _ = FillIP(network, 2)
	fmt.Println("network:", network)
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

func FillIP(ip string, ipType int) (string, error) {

	prefix := ""
	ipp := ip
	pos := strings.Index(ip, "/")
	if pos > 0 {
		prefix = string(ip[pos:])
		ipp = string(ip[:pos])
	}
	fmt.Println(" ipp:", ipp, "   prefix:", prefix, "   pos:", pos)

	if ipType == 1 {
		countComma := strings.Count(ipp, ".")
		if countComma == 3 {
			return ipp + prefix, nil
		} else if countComma < 3 {
			return ipp + strings.Repeat(".0", net.IPv4len-countComma-1) + prefix, nil
		} else {
			return "", errors.New("illegal ipv4")
		}
	} else if ipType == 2 {
		countColon := strings.Count(ipp, ":")
		if countColon == 7 {
			return ipp + prefix, nil
		} else if strings.HasSuffix(ipp, "::") {
			return ipp + prefix, nil
		} else {
			return ipp + "::" + prefix, nil
		}

	} else {
		return "", errors.New("illegal ipType")
	}

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
