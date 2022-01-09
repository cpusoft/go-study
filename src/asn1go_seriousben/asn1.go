package main

import (
	"fmt"

	goasn1 "github.com/seriousben/go-asn1"
)

func main() {

	enc, err := goasn1.ParsePemFile("ec384-public.pem")
	fmt.Println(enc, err)
}
