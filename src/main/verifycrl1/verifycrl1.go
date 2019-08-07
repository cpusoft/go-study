package main

import (
	"crypto/x509"

	"fmt"
	"io/ioutil"
)

func main() {
	path := `G:\Download\cert\verify\4\`

	cerFile := path + `inter.cer` //err
	//cerFile := path + `bW-_qXU9uNhGQz21NR2ansB8lr0.cer`  //ok
	cerByte, err := ioutil.ReadFile(cerFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	cer, err := x509.ParseCertificate(cerByte)
	if err != nil {
		fmt.Println(err)
		return
	}

	crlFile := path + `bW-_qXU9uNhGQz21NR2ansB8lr0.crl`
	crlByte, err := ioutil.ReadFile(crlFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	crl, err := x509.ParseCRL(crlByte)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cer.CheckCRLSignature(crl)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("check crl ok")
}
