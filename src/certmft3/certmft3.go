package main

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
)

func DecodeBase64(oldBytes []byte) ([]byte, error) {
	isBinary := false

	for _, b := range oldBytes {
		t := int(b)

		if t < 32 && t != 9 && t != 10 && t != 13 {
			isBinary = true
			break
		}
	}

	if isBinary {
		return oldBytes, nil
	}
	txt := string(oldBytes)
	txt = strings.Replace(txt, "-----BEGIN CERTIFICATE-----", "", -1)
	txt = strings.Replace(txt, "-----END CERTIFICATE-----", "", -1)
	txt = strings.Replace(txt, "-", "", -1)
	txt = strings.Replace(txt, " ", "", -1)
	txt = strings.Replace(txt, "\r", "", -1)
	txt = strings.Replace(txt, "\n", "", -1)
	newBytes, err := base64.StdEncoding.DecodeString(txt)
	return newBytes, err

}

func Trim00(oldByte []byte) []byte {
	nullbytes := []byte{0x00, 0x00}
	if bytes.HasSuffix(oldByte, nullbytes) {
		oldByte = oldByte[:len(oldByte)-len(nullbytes)]
	}
	return oldByte
}

func ReadFile(file string) ([]byte, error) {
	oldByte, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	newByte, err := DecodeBase64(oldByte)
	if err != nil {
		return nil, err
	}
	return newByte, nil
}

func main() {
	certFile := `E:\Go\parse-cert\data\029d506f4cfb1a1e4eae7d68b5ebbc15c8b52c93.mft`
	certFile = `G:\Download\.mft698152231.cer`
	fileByte, err := ReadFile(certFile)
	fmt.Println(err)
	fileByte = Trim00(fileByte)
	cer, err := x509.ParseCertificate(fileByte)
	fmt.Println(cer.Version)
	fmt.Println(err)
}
