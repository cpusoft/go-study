package main

import (
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"io/ioutil"
)

func main() {
	path := `G:\Download\cert\verify\2\`

	roots := x509.NewCertPool()
	rootFile := path + `inter.cer`
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
	fmt.Println("root issuer:", rootcert.Issuer.String(), "   subject:", rootcert.Subject.String())
	roots.AddCert(rootcert)
	/*
		inter := x509.NewCertPool()
		inter3File := path + `inter.cer`
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
	*/
	cer1File := path + `A9.cer`
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
	fmt.Println("cert issuer:", cer1cert.Issuer.String(), "   subject:", cer1cert.Subject.String())

	opts := x509.VerifyOptions{
		Roots: roots,
		//Intermediates: inter,
		//KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		//KeyUsages: []x509.ExtKeyUsage{x509.KeyUsageCertSign},
	}
	if chains, err := cer1cert.Verify(opts); err != nil {
		fmt.Println("failed to verify certificate: ", err, len(chains))
		for _, one := range chains {
			fmt.Println("chains ", one)
		}
		return
	}
	fmt.Printf("Success!\n")

}
