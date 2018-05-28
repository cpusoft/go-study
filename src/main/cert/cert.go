package main

import (
	_ "crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	_ "encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// cer 文件
type CerInfo struct {
	Version               int      `json:"version"`
	SN                    string   `json:"sn"`
	NotBefore             string   `json:"notBefore"`
	NotAfter              string   `json:"notAfter"`
	BasicConstraintsValid bool     `json:"basicConstraintsValid"`
	IsRoot                bool     `json:"isRoot"`
	DNSNames              []string `json:"dnsNames"`
	EmailAddresses        []string `json:"emailAddresses"`
	IPAddresses           []net.IP `json:"ipAddresses"`
	Subject               string   `json:"subject"`
	SubjectAll            string   `json:"subjectAll"`
	Issuer                string   `json:"issuer"`
	IssuerAll             string   `json:"issuerAll"`
}

type CrlRevokedCert struct {
	SN             string `jsong:"sn"`
	RevocationTime string `json:"revocationTime"`
}

type CrlInfo struct {
	Version         int              `json:"version"`
	Issuer          string           `json:"issuer"`
	ThisUpdate      string           `json:"thisUpdate"`
	NextUpdate      string           `json:"nextUpdate"`
	HasExpired      string           `json:"hasExpired"`
	CrlRevokedCerts []CrlRevokedCert `json:"CrlRevokedCerts"`
}

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

func getDNFromName(namespace pkix.Name, sep string) (string, error) {
	return getDNFromRDNSeq(namespace.ToRDNSequence(), sep)
}

func getDNFromRDNSeq(rdns pkix.RDNSequence, sep string) (string, error) {
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
	return sep + strings.Join(subject, sep), nil
}

func parseCer(file string) error {
	//300E300C040200013006030402B9A6FC     16
	/*  ` 0x0c 0x04 0x02 0x00 0x01 0x30 0x06 0x03 0x04 0x02 0xb9 0xa6 0xfc
	300E300C040200013006030402B9A6FC
	300D300B04020001300503030084FC

	oidValue:
	0x30 0x0e      0x30是SEQUENCE类型固定的， 0e是后面长度
		0x30 0x0c  0x30是SEQUENCE类型固定的， 0c是后面长度, 从这里开始
			0x04 0x02 0x00 0x01     0x04, 0x02, 0x00, 0x01, // address family: IPv4    对比：0x04, 0x02, 0x00, 0x02, // address family: IPv6
				0x30 0x06
					0x03 0x04
					 0x02 0xb9 0xa6 0xfc
					      185.166.252/22
	type: 48
	len: 14
	oidIP:
	0x30 0x0c 0x04 0x02 0x00 0x01 0x30 0x06 0x03 0x04 0x02 0xb9 0xa6 0xfc
	`*/

	rootCa := file
	caBlock, err := ioutil.ReadFile(rootCa)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return err
	}

	cert, err := x509.ParseCertificate(caBlock)
	if err != nil {
		fmt.Println("ParseCertificate err:", err)
		return err
	}

	cerInfo := CerInfo{}
	cerInfo.SN = fmt.Sprintf("%x", cert.SerialNumber)
	cerInfo.Version = cert.Version
	cerInfo.DNSNames = cert.DNSNames
	cerInfo.EmailAddresses = cert.EmailAddresses
	cerInfo.IPAddresses = cert.IPAddresses
	cerInfo.BasicConstraintsValid = cert.BasicConstraintsValid
	cerInfo.IsRoot = cert.IsCA
	cerInfo.NotBefore = cert.NotBefore.Format("2006-01-02 15:04:05")
	cerInfo.NotAfter = cert.NotAfter.Format("2006-01-02 15:04:05")
	cerInfo.Subject = cert.Subject.CommonName
	cerInfo.SubjectAll, _ = getDNFromName(cert.Subject, "/")
	cerInfo.Issuer = cert.Issuer.CommonName
	cerInfo.IssuerAll, _ = getDNFromName(cert.Issuer, "/")
	jsonCer, _ := json.Marshal(cerInfo)
	fmt.Printf("%+v", string(jsonCer))
	/*
		fmt.Println("valfrom:", cert.NotBefore.Format("2006-01-02 15:04:05"))
		fmt.Println("valto:", cert.NotAfter.Format("2006-01-02 15:04:05"))
		fmt.Printf("subject: /CN=%s/serialNumber=%s\r\n", cert.Subject.CommonName, cert.Subject.SerialNumber)
		fmt.Printf("issuer: /CN=%s/serialNumber=%s\r\n", cert.Issuer.CommonName, cert.Issuer.SerialNumber)
		sn := cert.SerialNumber.Uint64()
		fmt.Printf("sn: %d\r\n", sn)
		// flags ??
		ski := cert.SubjectKeyId
		printBytes("ski", ski)
		fmt.Println(printBase64(ski))

		publicInfo := cert.RawSubjectPublicKeyInfo
		printBytes("publicInfo", publicInfo)
		publicKey := cert.PublicKey
		//printBytes("publicKey", publicKey)
		fmt.Println(publicKey)

		aki := cert.AuthorityKeyId
		printBytes("aki", aki)
		fmt.Println("len(ExtKeyUsage):", len(cert.ExtKeyUsage))
		for _, eku := range cert.ExtKeyUsage {
			fmt.Println("%v\r\n", eku)
		}
		//颁发机构信息访问
		fmt.Println("aia:", cert.IssuingCertificateURL)

		crldp := cert.CRLDistributionPoints
		fmt.Printf("crldp:%+v\r\n", crldp)

		fmt.Printf("Extensions: %+v\r\n", cert.Extensions)
		fmt.Printf("ExtraExtensions: %+v\r\n", cert.ExtraExtensions)

		for _, extension := range cert.Extensions {
			oid := extension.Id
			value := extension.Value
			fmt.Println("oid", oid)
			printBytes("value", value)
			fmt.Println(printBase64(value))

		}
	*/
	return nil
}

func parseCrl(file string) error {
	crl := file
	crlBlock, err := ioutil.ReadFile(crl)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return err
	}
	certList, err := x509.ParseCRL(crlBlock)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tbsCertList := certList.TBSCertList
	crlInfo := CrlInfo{}
	crlInfo.Version = tbsCertList.Version
	crlInfo.Issuer, _ = getDNFromRDNSeq(tbsCertList.Issuer, "/")
	crlInfo.ThisUpdate = tbsCertList.ThisUpdate.Local().Format("2006-01-02 15:04:05") //LoadLocation("Local")
	crlInfo.NextUpdate = tbsCertList.NextUpdate.Local().Format("2006-01-02 15:04:05")
	crlInfo.HasExpired = strconv.FormatBool(certList.HasExpired(time.Now()))
	//exts := tbsCertList.Extensions
	crlInfo.CrlRevokedCerts = make([]CrlRevokedCert, 0)
	revokedCerts := tbsCertList.RevokedCertificates
	for _, revokedCert := range revokedCerts {
		crlRevokedCert := CrlRevokedCert{}
		crlRevokedCert.SN = fmt.Sprintf("%x", revokedCert.SerialNumber)
		crlRevokedCert.RevocationTime = revokedCert.RevocationTime.Local().Format("2006-01-02 15:04:05")
		crlInfo.CrlRevokedCerts = append(crlInfo.CrlRevokedCerts, crlRevokedCert)
	}

	jsonCrl, _ := json.Marshal(crlInfo)
	fmt.Printf("%+v", string(jsonCrl))

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: ./cert 1.cer")
		return
	}
	//`E:\Go\go-study\src\main\cert\root.cer`
	var err error
	certFile := strings.ToLower(os.Args[1])
	if strings.HasSuffix(certFile, ".cer") {
		err = parseCer(certFile)
	} else if strings.HasSuffix(certFile, ".crl") {
		err = parseCrl(certFile)
	}
	if err != nil {
		fmt.Println(err)
	}

}

func printBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func printAsn(name string, typ byte, ln byte, byt []byte) {
	fmt.Println(fmt.Sprintf(name+"Type:0x%02x (%d)", typ, typ))
	fmt.Println(fmt.Sprintf(name+"Len:0x%02x (%d)", ln, ln))
	printBytes(name+"Value:", byt)
}

func printBytes(name string, byt []byte) {
	fmt.Println(name)
	for _, i := range byt {
		fmt.Print(fmt.Sprintf("0x%02x ", i))
	}
	fmt.Println("")
}
