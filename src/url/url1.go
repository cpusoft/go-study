package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {
	s := "rsync://rpki.afrinic.net/repository/AfriNIC.cer"
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	pos := strings.LastIndex(u.Path, "/")
	fmt.Println(pos)
	fmt.Println(u.Host + string(u.Path[:pos+1]))
}
