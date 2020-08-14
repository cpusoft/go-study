package main

import (
	"encoding/asn1"
	"fmt"
	"time"

	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type CertificateList struct {
	Oid asn1.ObjectIdentifier
	Mft []asn1.RawValue `asn1:"optional,explicit,default:0,tag:0""`
}

type MftCert struct {
	ManifestParse ManifestParse
	Certificate   Certificate `asn1:"optional,explicit,default:0,tag:0""`
}

type ManifestParse struct {
	ManifestNumber asn1.RawValue         `json:"manifestNumber"`
	ThisUpdate     time.Time             `asn1:"generalized" json:"thisUpdate"`
	NextUpdate     time.Time             `asn1:"generalized" json:"nextUpdate"`
	FileHashAlg    asn1.ObjectIdentifier `json:"fileHashAlg"`
	FileList       []FileAndHashParse    `json:"fileList"`
}
type FileAndHashParse struct {
	File string         `asn1:"ia5" json:"file"`
	Hash asn1.BitString `json:"hash"`
}

func main() {
	var file string
	file = `E:\Go\go-study\src\asnmft1\1.mft`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}
	certificate := CertificateList{}
	asn1.Unmarshal(b, &certificate)
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate))

	mftCert := MftCert{}
	_, err = asn1.Unmarshal(certificate.Mft[0].Bytes, &mftCert)
	fmt.Println("mftCert:", jsonutil.MarshallJsonIndent(mftCert), err)

}
