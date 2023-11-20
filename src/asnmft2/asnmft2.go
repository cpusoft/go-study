package main

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"time"

	"github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
	model "labscm.zdns.cn/rpstir2-mod/rpstir2-model"
)

type OctetString []byte
type MftFileAsn1 struct {
	SignedDataOidAsn1 asn1.ObjectIdentifier `json:"signedDataOidAsn1"`
	SignedDataAsn1s   []asn1.RawValue       `json:"signedDataAsn1s" asn1:"optional,explicit,default:0,tag:0"`
}

type MftOctetStringAsn1 struct {
	RoaOidAsn1         asn1.ObjectIdentifier
	MftOctetStringAsn1 OctetString `asn1:"tag:0,explicit,optional"`
}

// asID as in rfc6482
type MftBlockAsn1 struct {
	MftNumberAsn1         *big.Int              `json:"mftNumberAsn1"`
	ThisUpdateAsn1        time.Time             `asn1:"generalized" json:"thisUpdateAsn1"`
	NextUpdateAsn1        time.Time             `asn1:"generalized" json:"nextUpdateAsn1"`
	FileHashAlgorithmAsn1 asn1.ObjectIdentifier `json:"fileHashAlgorithmAsn1"`
	FileAndHashAsn1s      []FileAndHashAsn1     `json:"fileAndHashAsn1s"`
}
type FileAndHashAsn1 struct {
	FileAsn1 string         `asn1:"ia5" json:"fileAsn1"`
	HashAsn1 asn1.BitString `json:"hashAsn1"`
}

// data: asn1.FullBytes
func ParseToMftBlockAsn1(data []byte) (mftBlockAsn1 MftBlockAsn1, err error) {
	belogs.Debug("ParseToMftBlockAsn1(): len(data):", len(data))
	var mftOctetStringAsn1 MftOctetStringAsn1
	_, err = asn1.Unmarshal(data, &mftOctetStringAsn1)
	if err != nil {
		belogs.Error("ParseToMftBlockAsn1(): Unmarshal to mftOctetStringAsn1 fail:", err)
		return
	}
	belogs.Debug("ParseToMftBlockAsn1(): mftOctetStringAsn1:", jsonutil.MarshalJson(mftOctetStringAsn1))

	_, err = asn1.Unmarshal([]byte(mftOctetStringAsn1.MftOctetStringAsn1), &mftBlockAsn1)
	if err != nil {
		belogs.Error("ParseToMftBlockAsn1(): Unmarshal to mftBlockAsn1 fail:", err)
		return
	}
	belogs.Debug("ParseToMftBlockAsn1():mftBlockAsn1:", jsonutil.MarshalJson(mftBlockAsn1))
	return
}

type Sha256 struct {
	Oid  asn1.ObjectIdentifier
	Null asn1.RawValue
}

type OidAndValueAsn1 struct {
	OidAsn1   asn1.ObjectIdentifier `json:"oidAsn1"`
	ValueAsn1 asn1.RawValue         `json:"valueAsn1" asn1:"optional"`
}

type SignedAttributeAsn1s struct {
	OidAndValueAsn1s []OidAndValueAsn1 `json:"attributeTypeAndValues" asn1:"tag:0"`
}
type SignerInfoAsn1 struct {
	Version                int               `json:"version"`
	Sid                    OctetString       `json:"sid" asn1:"tag:0"`
	DigestAlgorithmAsn1    OidAndValueAsn1   `json:"digestAlgorithm"`
	OidAndValueAsn1s       []OidAndValueAsn1 `json:"attributeTypeAndValues" asn1:"tag:0"`
	SignatureAlgorithmAsn1 OidAndValueAsn1   `json:"signatureAlgorithm"`
	Sinagture              OctetString       `json:"sinagture"`
}

type OidAndValuesAsn1 struct {
	OidAsn1    asn1.ObjectIdentifier `json:"oidAsn1"`
	ValueAsn1s []asn1.RawValue       `json:"valueAsn1"`
}
type SignerInfoAsn1_new struct {
	Version             int               `json:"version"`
	Sid                 jsonutil.HexBytes `json:"sid" asn1:"tag:0"`
	DigestAlgorithmAsn1 OidAndValuesAsn1  `json:"digestAlgorithm"`
	//SignedAttributeAsn1s   SignedAttributeAsn1s `json:"signedAttributeAsn1s" asn1:"tag:0"`
	OidAndValueAsn1s       []asn1.RawValue `json:"attributeTypeAndValues" asn1:"tag:0"`
	SignatureAlgorithmAsn1 asn1.RawValue   `json:"signatureAlgorithm"`
	//Sinagture              OctetString     `json:"sinagture"`
}

type SignerInfoAsn1_raw struct {
	Version             int               `json:"version"`
	Sid                 jsonutil.HexBytes `json:"sid" asn1:"tag:0"`
	DigestAlgorithmAsn1 asn1.RawValue     `json:"digestAlgorithm"`
	//SignedAttributeAsn1s   SignedAttributeAsn1s `json:"signedAttributeAsn1s" asn1:"tag:0"`
	OidAndValueAsn1s       []asn1.RawValue `json:"attributeTypeAndValues"`
	SignatureAlgorithmAsn1 asn1.RawValue   `json:"signatureAlgorithm"`
	//Sinagture              OctetString     `json:"sinagture"`
}

// https://github.com/blacktop/ipsw/blob/master/internal/codesign/cms/cms.go
type Attribute struct {
	Type asn1.ObjectIdentifier
	// This should be a SET OF ANY, but Go's asn1 parser can't handle slices of
	// RawValues. Use value() to get an AnySet of the value.
	RawValue asn1.RawValue
}

type Attributes []Attribute

type SignerInfo_ipsw struct {
	Version            int
	SID                asn1.RawValue
	DigestAlgorithm    pkix.AlgorithmIdentifier
	SignedAttrs        Attributes `asn1:"optional,tag:0"`
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Signature          []byte
	UnsignedAttrs      Attributes `asn1:"set,optional,tag:1"`
}

////////////////////////// cms.go

func ParseMftModelByAsn1(fileByte []byte, mftModel *model.MftModel) (err error) {
	start := time.Now()
	//	fmt.Println("ParseMftModelByAsn1(): len(fileByte):", len(fileByte))
	mftFileAsn1 := MftFileAsn1{}
	_, err = asn1.Unmarshal(fileByte, &mftFileAsn1)
	if err != nil {
		fmt.Println("ParseMftModelByAsn1(): Unmarshal mftFileAsn1 fail, len(fileByte):", len(fileByte), err)
		return err
	}
	/*
	   ParseMftModelByAsn1(): seq.Tag: 2   seq.Class: 0   seq.IsCompound: false
	   ParseMftModelByAsn1(): seq.Tag: 17   seq.Class: 0   seq.IsCompound: true
	   ParseMftModelByAsn1(): seq.Tag: 16   seq.Class: 0   seq.IsCompound: true
	   ParseMftModelByAsn1(): seq.Tag: 0   seq.Class: 2   seq.IsCompound: true
	   ParseMftModelByAsn1(): seq.Tag: 17   seq.Class: 0   seq.IsCompound: true
	*/

	fmt.Println("ParseMftModelByAsn1(): len(SignedDataAsn1s):", len(mftFileAsn1.SignedDataAsn1s))
	for _, seq := range mftFileAsn1.SignedDataAsn1s {
		//	fmt.Println("ParseMftModelByAsn1(): seq.Tag:", seq.Tag, "  seq.Class:", seq.Class, "  seq.IsCompound:", seq.IsCompound)

		if seq.Class == 0 && seq.Tag == 2 && !seq.IsCompound {
			// version CMSVersion INTEGER 3: ignore
		} else if seq.Class == 0 && seq.Tag == 17 && seq.IsCompound && len(seq.Bytes) < 100 {
			// digestAlgorithms DigestAlgorithmIdentifiers SET (1 elem) : ignore
		} else if seq.Class == 0 && seq.Tag == 16 && seq.IsCompound {
			// //  encapContentInfo EncapsulatedContentInfo
			/*
				mftBlockAsn1, err := ParseToMftBlockAsn1(seq.FullBytes)
				if err != nil {
					fmt.Println("ParseMftModelByAsn1(): ParseToMftBlockAsn1 fail, len(seq.FullBytes):", len(seq.FullBytes), err)
					continue
				}
				fmt.Println("ParseMftModelByAsn1(): mftBlockAsn1:", jsonutil.MarshalJson(mftBlockAsn1))
			*/
		} else if seq.Class == 2 && seq.Tag == 0 && seq.IsCompound {
			// EeModel will call
		} else if seq.Class == 0 && seq.Tag == 17 && seq.IsCompound && len(seq.Bytes) > 100 {
			// signerInfos SignerInfos will call
			var signerInfoAsn1 SignerInfoAsn1
			_, err = asn1.Unmarshal(seq.Bytes, &signerInfoAsn1)
			if err != nil {
				fmt.Println("ParseToSignerInfoModel(): SignerInfoAsn1 fail, len(seq.Bytes):", len(seq.Bytes), err)
			}
			fmt.Println("ParseToSignerInfoModel(): signerInfoAsn1:", jsonutil.MarshalJson(signerInfoAsn1)+"\n")

			var signerInfoAsn1_new SignerInfoAsn1_new
			_, err = asn1.Unmarshal(seq.Bytes, &signerInfoAsn1_new)
			if err != nil {
				fmt.Println("ParseToSignerInfoModel(): signerInfoAsn1_new fail, len(seq.Bytes):", len(seq.Bytes), err)
			}
			fmt.Println("ParseToSignerInfoModel(): signerInfoAsn1_new:", jsonutil.MarshalJson(signerInfoAsn1_new)+"\n")

			var signerInfoAsn1_raw SignerInfoAsn1_raw
			_, err = asn1.Unmarshal(seq.Bytes, &signerInfoAsn1_raw)
			if err != nil {
				fmt.Println("ParseToSignerInfoModel(): signerInfoAsn1_raw fail, len(seq.Bytes):", len(seq.Bytes), err)
			}
			fmt.Println("ParseToSignerInfoModel(): signerInfoAsn1_raw:", jsonutil.MarshalJson(signerInfoAsn1_raw)+"\n")

			var signerInfo_ipsw SignerInfo_ipsw
			_, err = asn1.Unmarshal(seq.Bytes, &signerInfo_ipsw)
			if err != nil {
				fmt.Println("ParseToSignerInfoModel(): signerInfo_ipsw fail, len(seq.Bytes):", len(seq.Bytes), err)
			}
			fmt.Println("ParseToSignerInfoModel(): signerInfo_ipsw:", jsonutil.MarshalJson(signerInfo_ipsw)+"\n")

			/*
				var signerInfoAsn1_1 SignerInfoAsn1_1
				_, err = asn1.Unmarshal(seq.Bytes, &signerInfoAsn1_1)
				if err != nil {
					fmt.Println("ParseToSignerInfoModel(): SignerInfoAsn1 fail, len(seq.Bytes):", len(seq.Bytes), err)
					continue
				}
				fmt.Println("ParseToSignerInfoModel(): signerInfoAsn1_1:", jsonutil.MarshalJson(signerInfoAsn1_1))
			*/
		}
	}
	fmt.Println("ParseMftModelByAsn1(): ok", "  time(s):", time.Since(start))

	return
}

func main() {

	files := []string{`G:\Download\cert\mft\1.mft`, `G:\Download\cert\mft\c.mft`}
	for _, file := range files {
		fmt.Println(file)
		b, err := fileutil.ReadFileToBytes(file)
		if err != nil {
			fmt.Println(file, err)
			continue
		}
		mftModel := model.MftModel{}
		err = ParseMftModelByAsn1(b, &mftModel)
		if err != nil {
			fmt.Println("ParseMftModelByAsn1() fail:", file, err)
			continue
		}
		fmt.Println(mftModel.String() + "\n\n\n\n\n")
	}
}
