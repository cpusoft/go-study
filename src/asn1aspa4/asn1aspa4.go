package main

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cpusoft/goutil/belogs"
	//"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
	"github.com/cpusoft/goutil/osutil"
	model "labscm.zdns.cn/rpstir2-mod/rpstir2-model"
)

type OctetString []byte
type AsaFileAsn1 struct {
	SignedDataOidAsn1 asn1.ObjectIdentifier `json:"signedDataOidAsn1"`
	SignedDataAsn1s   []asn1.RawValue       `json:"signedDataAsn1s" asn1:"optional,explicit,default:0,tag:0"`
}

type AsaOctetStringAsn1 struct {
	AsaOidAsn1         asn1.ObjectIdentifier
	AsaOctetStringAsn1 OctetString `asn1:"tag:0,explicit,optional"`
}

// 1.2.840.113549.1.9.16.1.49
type AsaBlockAsn1 struct {
	//VersionAsn1      Version `json:"versionAsn1" asn1:"class:2,tag:0,optional"` //default 0
	CustomerAsnAsn1    int                 `json:"customerAsnAsn1"`
	ProviderBlockAsn1s []ProviderBlockAsn1 `json:"providerBlockAsn1s" asn1:"optional"`
}

type ProviderBlockAsn1 struct {
	ProviderAsnAsn1 int `json:"providerAsnAsn1"`
}

type AfiAsn1 struct {
	Afi int
}
type AsaBlockOldAsn1 struct {
	AfiAsn1          AfiAsn1 `asn1:"class:2,tag:0"` //asn1.RawValue
	CustomerAsnAsn1  int     //asn1.RawValue   //`asn1:"explicit,tag:5"`
	ProviderAsnAsn1s []int   //`asn1:"explicit,tag:5"`
}

func convertAsaBlockOldAsn1ToAsaBlockAsn1(old AsaBlockOldAsn1) AsaBlockAsn1 {
	as := AsaBlockAsn1{}
	as.CustomerAsnAsn1 = old.CustomerAsnAsn1
	as.ProviderBlockAsn1s = make([]ProviderBlockAsn1, 0)
	for i := range old.ProviderAsnAsn1s {
		providerBlockAsn1 := ProviderBlockAsn1{
			ProviderAsnAsn1: old.ProviderAsnAsn1s[i],
		}
		as.ProviderBlockAsn1s = append(as.ProviderBlockAsn1s, providerBlockAsn1)
	}
	return as
}

// data: asn1.FullBytes
func ParseToAsaBlockAsn1(data []byte) (asaBlockAsn1 AsaBlockAsn1, err error) {
	belogs.Debug("ParseToAsaBlockAsn1(): len(data):", len(data))
	var asaOctetStringAsn1 AsaOctetStringAsn1
	_, err = asn1.Unmarshal(data, &asaOctetStringAsn1)
	if err != nil {
		belogs.Error("ParseToAsaBlockAsn1(): Unmarshal to asaOctetStringAsn1 fail:", err)
		return
	}
	belogs.Debug("ParseToAsaBlockAsn1(): asaOctetStringAsn1:", jsonutil.MarshalJson(asaOctetStringAsn1),
		"    data:", hex.EncodeToString([]byte(asaOctetStringAsn1.AsaOctetStringAsn1)))

	asaBlockAsn1 = AsaBlockAsn1{}
	_, err = asn1.Unmarshal([]byte(asaOctetStringAsn1.AsaOctetStringAsn1), &asaBlockAsn1)
	if err != nil {
		belogs.Debug("ParseToAsaBlockAsn1(): Unmarshal to asaBlockAsn1, try AsaBlockOldAsn1:", err)

		asaBlockOldAsn1 := AsaBlockOldAsn1{}
		_, err = asn1.Unmarshal([]byte(asaOctetStringAsn1.AsaOctetStringAsn1), &asaBlockOldAsn1)
		if err != nil {
			belogs.Error("ParseToAsaBlockAsn1(): Unmarshal to asaBlockOldAsn1 fail:", hex.EncodeToString([]byte(asaOctetStringAsn1.AsaOctetStringAsn1)),
				err)
			return
		}
		belogs.Debug("ParseToAsaBlockAsn1(): asaBlockOldAsn1:", jsonutil.MarshalJson(asaBlockOldAsn1))
		asaBlockAsn1 = convertAsaBlockOldAsn1ToAsaBlockAsn1(asaBlockOldAsn1)
	}
	belogs.Debug("ParseToAsaBlockAsn1(): Unmarshal to asaBlockAsn1:", jsonutil.MarshalJson(asaBlockAsn1), err)

	return
}

////////////////////////// cms.go

func ParseAsaModelByAsn1(fileByte []byte, asaModel *model.AsaModel) (err error) {
	start := time.Now()
	//	fmt.Println("ParseAsaModelByAsn1(): len(fileByte):", len(fileByte))
	asaFileAsn1 := AsaFileAsn1{}
	_, err = asn1.Unmarshal(fileByte, &asaFileAsn1)
	if err != nil {
		fmt.Println("ParseAsaModelByAsn1(): Unmarshal asaFileAsn1 fail, len(fileByte):", len(fileByte), err)
		return err
	}
	/*
		ParseAsaModelByAsn1(): seq.Tag: 2   seq.Class: 0   seq.IsCompound: false
		ParseAsaModelByAsn1(): seq.Tag: 17   seq.Class: 0   seq.IsCompound: true
		ParseAsaModelByAsn1(): seq.Tag: 16   seq.Class: 0   seq.IsCompound: true
		ParseAsaModelByAsn1(): seq.Tag: 0   seq.Class: 2   seq.IsCompound: true
		ParseAsaModelByAsn1(): seq.Tag: 17   seq.Class: 0   seq.IsCompound: true
	*/

	fmt.Println("ParseAsaModelByAsn1(): len(SignedDataAsn1s):", len(asaFileAsn1.SignedDataAsn1s))
	for _, seq := range asaFileAsn1.SignedDataAsn1s {
		fmt.Println("ParseAsaModelByAsn1(): seq.Tag:", seq.Tag, "  seq.Class:", seq.Class, "  seq.IsCompound:", seq.IsCompound)

		if seq.Class == 0 && seq.Tag == 2 && !seq.IsCompound {
			// version CMSVersion INTEGER 3: ignore
		} else if seq.Class == 0 && seq.Tag == 17 && seq.IsCompound && len(seq.Bytes) < 100 {
			// digestAlgorithms DigestAlgorithmIdentifiers SET (1 elem) : ignore
		} else if seq.Class == 0 && seq.Tag == 16 && seq.IsCompound {
			// //  encapContentInfo EncapsulatedContentInfo

			asaBlockAsn1, err := ParseToAsaBlockAsn1(seq.FullBytes)
			if err != nil {
				fmt.Println("ParseAsaModelByAsn1(): ParseToAsaBlockAsn1 fail, len(seq.FullBytes):", len(seq.FullBytes), err)
				return err
			}
			fmt.Println("ParseAsaModelByAsn1(): asaBlockAsn1:", jsonutil.MarshalJson(asaBlockAsn1))

		} else if seq.Class == 2 && seq.Tag == 0 && seq.IsCompound {
			// EeModel will call
		} else if seq.Class == 0 && seq.Tag == 17 && seq.IsCompound && len(seq.Bytes) > 100 {

		}
	}
	fmt.Println("ParseAsaModelByAsn1(): ok", "  time(s):", time.Since(start))

	return
}

func main() {
	/*
		files := []string{`G:\Download\cert\asa2\AS970.asa`}
		path := ``
	*/
	path := `G:\Download\cert\asa2\`
	m := make(map[string]string, 0)
	m[".asa"] = ".asa"
	files, _ := osutil.GetFilesInDir(path, m)

	for _, file := range files {
		file = path + file
		fmt.Println(file)
		b, err := fileutil.ReadFileToBytes(file)
		if err != nil {
			fmt.Println(file, err)
			continue
		}
		asaModel := model.AsaModel{}
		err = ParseAsaModelByAsn1(b, &asaModel)
		if err != nil {
			fmt.Println("ParseAsaModelByAsn1() fail:", file, err)
			continue
		}
		fmt.Println(file + "\n" + asaModel.String() + "\n\n\n\n\n")
	}
}
