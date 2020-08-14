package main

import (
	"encoding/asn1"
	"fmt"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type CertificateList struct {
	Oid  asn1.ObjectIdentifier
	Seqs []asn1.RawValue `asn1:"optional,explicit,default:0,tag:0""`
}

type MftCertificateList struct {
	Oid   asn1.ObjectIdentifier
	Value asn1.RawValue `asn1:"explicit,tag:0"`
	//Seqs ManifestParse1 `asn1:"explicit,tag:0"`
}
type ManifestParse1 struct {
	ManifestParse2 []asn1.RawValue //ManifestParse2
	//ManifestParse []ManifestParse
}
type ManifestParse2 struct {
	ManifestParse asn1.RawValue
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

type Sha256 struct {
	Oid  asn1.ObjectIdentifier
	Null asn1.RawValue
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

	// int 3 ???
	fmt.Println("seqs[0]", convert.Bytes2Uint64(certificate.Seqs[0].Bytes), "\n")

	// sha256
	sha256 := Sha256{}
	asn1.Unmarshal(certificate.Seqs[1].Bytes, &sha256)
	fmt.Println("seqs[1]:Oids:", jsonutil.MarshallJsonIndent(sha256.Oid))

	// mft
	oid1 := asn1.ObjectIdentifier{} //make([]asn1.RawValue, 0)
	reset, err := asn1.Unmarshal(certificate.Seqs[2].Bytes, &oid1)
	fmt.Println("seqs[2]:oid:", jsonutil.MarshallJsonIndent(oid1), err)

	raw := asn1.RawValue{}
	_, err = asn1.Unmarshal(reset, &raw)
	//fmt.Println("seqs[2]:", convert.PrintBytes(raw.Bytes, 8))

	raw1 := asn1.RawValue{}
	_, err = asn1.Unmarshal(raw.Bytes, &raw1)
	//fmt.Println("seqs[2]:", convert.PrintBytes(raw1.Bytes, 8))

	manifestParse := ManifestParse{}
	_, err = asn1.Unmarshal(raw1.Bytes, &manifestParse)
	fmt.Println("seqs[2]:mft:", jsonutil.MarshallJsonIndent(manifestParse))

	mftCertificateList := MftCertificateList{}

	_, err = asn1.Unmarshal(certificate.Seqs[2].FullBytes, &mftCertificateList)
	fmt.Println("seqs[2]2:", jsonutil.MarshallJsonIndent(mftCertificateList), err)

	manifestParse1 := ManifestParse1{}
	_, err = asn1.Unmarshal(mftCertificateList.Value.Bytes, &manifestParse1)
	fmt.Println("seqs[2]2:", jsonutil.MarshallJsonIndent(manifestParse1), err)
}
