package main

import (
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"io/ioutil"
)

func main() {
	path := `G:\Download\cert\verify\1\`

	roots := x509.NewCertPool()
	rootFile := path + `arin-rpki-ta.cer`
	rootFileByte, err := ioutil.ReadFile(rootFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	rootcert, err := x509.ParseCertificate(rootFileByte)
	if err != nil {
		fmt.Println(err)
		return
	}
	rootcert.UnhandledCriticalExtensions = make([]asn1.ObjectIdentifier, 0)
	roots.AddCert(rootcert)

	inter := x509.NewCertPool()
	inter3File := path + `3_5e4a23ea-e80a-403e-b08c-2171da2157d3.cer`
	inter3FileByte, err := ioutil.ReadFile(inter3File)
	if err != nil {
		fmt.Println(err)
		return
	}
	inter3cert, err := x509.ParseCertificate(inter3FileByte)
	if err != nil {
		fmt.Println(err)
		return
	}
	inter3cert.UnhandledCriticalExtensions = make([]asn1.ObjectIdentifier, 0)
	inter.AddCert(inter3cert)

	inter2File := path + `2_f60c9f32-a87c-4339-a2f3-6299a3b02e29.cer`
	inter2FileByte, err := ioutil.ReadFile(inter2File)
	if err != nil {
		fmt.Println(err)
		return
	}
	inter2cert, err := x509.ParseCertificate(inter2FileByte)
	if err != nil {
		fmt.Println(err)
		return
	}
	inter2cert.UnhandledCriticalExtensions = make([]asn1.ObjectIdentifier, 0)
	inter.AddCert(inter2cert)

	cer1File := path + `1_74939129-73e4-4131-8dea-7d0466eca8f3.cer`
	cer1FileByte, err := ioutil.ReadFile(cer1File)
	if err != nil {
		fmt.Println(err)
		return
	}
	cer1cert, err := x509.ParseCertificate(cer1FileByte)
	if err != nil {
		fmt.Println(err)
		return
	}
	cer1cert.UnhandledCriticalExtensions = make([]asn1.ObjectIdentifier, 0)

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: inter,
		//KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		//KeyUsages: []x509.ExtKeyUsage{x509.KeyUsageCertSign},
	}

	if chains, err := cer1cert.Verify(opts); err != nil {
		fmt.Println("failed to verify certificate: ", err, len(chains))
		fmt.Println("UnhandledCriticalExtensions   www: len:", len(cer1cert.UnhandledCriticalExtensions))
		for _, one := range chains {
			fmt.Println("chains ", one)
		}
		return
	}
	fmt.Printf("Success!\n")

}
