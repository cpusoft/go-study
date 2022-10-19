package main

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/cpusoft/goutil/convert"
)

func main() {
	ip := `1.1.1.1`
	addr := net.ParseIP(ip)
	fmt.Println("ip:", ip, " addr:", addr)

	b := addr.To4()
	fmt.Println("hex:", hex.EncodeToString(b), "  convert:", convert.PrintBytesOneLine(b))

	b = addr.To16()
	fmt.Println("hex:", hex.EncodeToString(b), "  convert:", convert.PrintBytesOneLine(b))

	ip6 := `2001:67c:1562::1c`
	addr6 := net.ParseIP(ip6)
	fmt.Println("ip6:", ip6, "  addr6:", addr6)

	b6 := addr6.To16()
	fmt.Println("hex6:", hex.EncodeToString(b6), "  convert:", convert.PrintBytesOneLine(b6))
}
