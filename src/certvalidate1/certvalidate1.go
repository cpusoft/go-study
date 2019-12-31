package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"strings"

	jsonutil "github.com/cpusoft/goutil/jsonutil"
	"github.com/square/certigo/pkcs7"
)

var oid = map[string]string{
	"2.5.4.3":                    "CN",
	"2.5.4.4":                    "SN",
	"2.5.4.5":                    "serialNumber",
	"2.5.4.6":                    "C",
	"2.5.4.7":                    "L",
	"2.5.4.8":                    "ST",
	"2.5.4.9":                    "streetAddress",
	"2.5.4.10":                   "O",
	"2.5.4.11":                   "OU",
	"2.5.4.12":                   "title",
	"2.5.4.17":                   "postalCode",
	"2.5.4.42":                   "GN",
	"2.5.4.43":                   "initials",
	"2.5.4.44":                   "generationQualifier",
	"2.5.4.46":                   "dnQualifier",
	"2.5.4.65":                   "pseudonym",
	"0.9.2342.19200300.100.1.25": "DC",
	"1.2.840.113549.1.9.1":       "emailAddress",
	"0.9.2342.19200300.100.1.1":  "userid",
	"2.5.29.20":                  "CRL Number",
}

func GetDNFromName(namespace pkix.Name, sep string) (string, error) {
	return GetDNFromRDNSeq(namespace.ToRDNSequence(), sep)
}

func GetDNFromRDNSeq(rdns pkix.RDNSequence, sep string) (string, error) {
	subject := []string{}
	for _, s := range rdns {
		for _, i := range s {
			if v, ok := i.Value.(string); ok {
				if name, ok := oid[i.Type.String()]; ok {
					// <oid name>=<value>
					subject = append(subject, fmt.Sprintf("%s=%s", name, v))
				} else {
					// <oid>=<value> if no <oid name> is found
					subject = append(subject, fmt.Sprintf("%s=%s", i.Type.String(), v))
				}
			} else {
				// <oid>=<value in default format> if value is not string
				subject = append(subject, fmt.Sprintf("%s=%v", i.Type.String, v))
			}
		}
	}
	return strings.Join(subject, sep), nil
}

//https://github.com/square/certigo/blob/master/lib/certs.go
func main() {
	file := `G:\Download\cert\cache\rpki.afrinic.net\repository\afrinic\0mlPWs5d97eiscdbRLGLWndoOGg.cer`
	data, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Errorf("unable to read input: %s\n", err)
		return
	}
	x509Certs, err0 := x509.ParseCertificates(data)
	if err0 == nil {
		fmt.Println("x509Certs:", len(x509Certs))
		for _, cert := range x509Certs {
			//	callback(EncodeX509ToPEM(cert, headers))
			fmt.Println(jsonutil.MarshalJson(cert))
			/*
							type KeyUsage int

				const (
					KeyUsageDigitalSignature KeyUsage = 1 << iota
					KeyUsageContentCommitment
					KeyUsageKeyEncipherment
					KeyUsageDataEncipherment
					KeyUsageKeyAgreement
					KeyUsageCertSign
					KeyUsageCRLSign
					KeyUsageEncipherOnly
					KeyUsageDecipherOnly
				)
			*/
			fmt.Println(fmt.Sprintf("keyusage:%0x", cert.KeyUsage))
			fmt.Println(x509.KeyUsageCertSign, x509.KeyUsageCRLSign)
			k := (x509.KeyUsageCertSign + x509.KeyUsageCRLSign)
			fmt.Println(fmt.Sprintf("keyCertSign + cRLSign:%0x", k))

			cs := cert.Subject.String()
			fmt.Println("Subject.String() cs:", cs)

			cs, _ = GetDNFromName(cert.Subject, ",")
			fmt.Println("GetDNFromRDNSeq cs:", cs)

			for _, ext := range cert.Extensions {
				fmt.Println(ext.Id)
				fmt.Println(ext.Critical)
			}

		}
		//	return
	} else {
		fmt.Println("x509Certs: err:", err0)
	}

	file = `G:\Download\cert\cache\rpki.afrinic.net\repository\afrinic\K1eJenypZMPIt_e92qek2jSpj4A.mft`
	data, err = ioutil.ReadFile(file)
	p7bBlocks, err1 := pkcs7.ParseSignedData(data)
	if err1 == nil {
		fmt.Println("p7bBlocks:", len(p7bBlocks))
		for _, block := range p7bBlocks {
			//	callback(pkcs7ToPem(block, headers))
			fmt.Println(jsonutil.MarshalJson(block))
		}
		//		return
	} else {
		fmt.Println("p7bBlocks: err:", err1)
	}
}
