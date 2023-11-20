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
	fmt.Println("mft:", convert.PrintBytesOneLine(raw1.Bytes))
	fmt.Println("seqs[2]:mft:", jsonutil.MarshallJsonIndent(manifestParse))

	// ee cer
	cer := Certificate{}
	asn1.Unmarshal(certificate.Seqs[3].Bytes, &cer)
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(cer))
	fmt.Println(len(cer.TBSCertificate.Extensions))
	for i := range cer.TBSCertificate.Extensions {
		extension := &cer.TBSCertificate.Extensions[i]
		fmt.Println(extension.Oid.String())
		if extension.Oid.String() == "2.5.29.14" {
			// subjectKeyIdentifier
			fmt.Print("ski:")
			fmt.Println(GetOctetString(extension.Value))
		} else if extension.Oid.String() == "2.5.29.35" {
			// authorityKeyIdentifier
			fmt.Print("aki:")
			fmt.Println(GetOctetStringSequenceString(extension.Value))
		} else if extension.Oid.String() == "2.5.29.19" {
			// basicConstraints

			fmt.Print("basic constraints:", extension.Critical)
			fmt.Println(GetOctetStringSequenceBool(extension.Value))
		} else if extension.Oid.String() == "2.5.29.15" {
			// keyUsage
			usageValue, err := GetOctetStringBitString(extension.Value)
			fmt.Println("keyUsage:", extension.Critical, usageValue, err)

			var tmp int
			// usageValue: 0000011
			// 从左边开始数，从0开始计数，即第5,6位为1, 则对应KeyUsageCertSign  KeyUsageCRLSign
			for i := 0; i < 9; i++ {
				//当为1时挪动，即看是第几个进行挪动
				//fmt.Println(i, usageValue.At(i))
				if usageValue.At(i) != 0 {
					tmp |= 1 << uint(i)
				}
			}
			// 先写死吧
			usage := int(tmp)
			usageStr := "Certificate Sign, CRL Sign"
			fmt.Println(usage)
			fmt.Println(usageStr)
			/*
				fmt.Println(x509.KeyUsageDigitalSignature)
				fmt.Println(x509.KeyUsageContentCommitment)
				fmt.Println(x509.KeyUsageKeyEncipherment)
				fmt.Println(x509.KeyUsageDataEncipherment)
				fmt.Println(x509.KeyUsageKeyAgreement)
				fmt.Println(x509.KeyUsageCertSign)
				fmt.Println(x509.KeyUsageCRLSign)
				fmt.Println(x509.KeyUsageEncipherOnly)
				fmt.Println(x509.KeyUsageDecipherOnly)
			*/

		} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.1" {
			// authorityInfoAccess
			seqs, err := GetOctetStringSequenceOidString(extension.Value)
			fmt.Println("aia:", err)
			for i := range seqs {
				fmt.Println(seqs[i].Oid, string(seqs[i].Value))
			}
		} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.11" {
			// subjectInfoAccess
			seqs, err := GetOctetStringSequenceOidString(extension.Value)
			fmt.Println("sia:", err)
			for i := range seqs {
				fmt.Println(seqs[i].Oid, string(seqs[i].Value))
			}
		} else if extension.Oid.String() == "2.5.29.31" {
			// cRLDistributionPoints
			seqs, err := GetCrldp(extension.Value)
			fmt.Println("crl:", seqs, err)
		} else if extension.Oid.String() == "2.5.29.32" {
			// Policies
			seqs, err := GetPolicies(extension.Value)
			fmt.Println("plicies:", seqs, err)
		} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.7" {
			// IpBlock
			//seqs, err := GetIpBlocks(extension.Value)
			//fmt.Println(seqs, err)
		} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.8" {
			// Asn
			//GetAsns(extension.Value)

		}
	}

	sd, sign, err := GetSignedData(certificate.Seqs[4].Bytes)
	fmt.Println("sign")
	fmt.Println(jsonutil.MarshallJsonIndent(sd.SignSha2561), jsonutil.MarshallJsonIndent(sign))
	fmt.Println(jsonutil.MarshallJsonIndent(sd.SignSha2562), convert.Bytes2String(sd.Tmp2))

}
