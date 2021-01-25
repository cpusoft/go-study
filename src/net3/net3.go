package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	txt, err := net.LookupTXT("AS266087.asn.cymru.com")
	fmt.Println(txt, err)
	splits := strings.Split(txt[0], "|")
	if len(splits) < 2 {
		return
	}

	s := splits[len(splits)-1] + "," + splits[2]
	fmt.Println(s)
}
