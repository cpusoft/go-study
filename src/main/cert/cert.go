package main

import (
	_ "crypto/tls"
	"crypto/x509"
	"encoding/base64"
	_ "encoding/pem"
	"fmt"
	"io/ioutil"
)

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
	fmt.Println(len(caBlock))

	cert, err := x509.ParseCertificate(caBlock)
	if err != nil {
		fmt.Println("ParseCertificate err:", err)
		return err
	}

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
	return nil
}
func main() {
	err := parseCer(`E:\Go\go-study\src\main\cert\ROUTER-0000FBF0_new.cer`)
	if err != nil {
		fmt.Println(err)
	}
	err = parseCer(`E:\Go\go-study\src\main\cert\ROUTER-00010000_new.cer`)
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
