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
type AsaModel struct {

	// must be 0, or no in file
	//  The version be 0.
	Version  int    `json:"version"`
	Ski      string `json:"ski" xorm:"ski varchar(128)"`
	Aki      string `json:"aki" xorm:"aki varchar(128)"`
	FilePath string `json:"filePath" xorm:"filePath varchar(512)"`
	FileName string `json:"fileName" xorm:"fileName varchar(128)"`
	FileHash string `json:"fileHash" xorm:"fileHash varchar(512)"`

	CustomerAsns []CustomerAsn `json:"customerAsns"`

	EContentType    string          `json:"eContentType"`
	AiaModel        AiaModel        `json:"aiaModel"`
	SiaModel        SiaModel        `json:"siaModel"`
	EeCertModel     EeCertModel     `json:"eeCertModel"`
	SignerInfoModel SignerInfoModel `json:"signerInfoModel"`
}
type CustomerAsn struct {
	Version      uint64   `json:"version"`
	CustomerAsn  uint64   `json:"customerAsn"`
	ProviderAsns []uint64 `json:"providerAsns"`
}

// AIA
type AiaModel struct {
	CaIssuers string `json:"caIssuers" xorm:"caIssuers varchar(512)"`
	Critical  bool   `json:"critical"`
}

type EeCertModel struct {
	// must be 3
	Version int `json:"version"`
	// SHA256-RSA: x509.SignatureAlgorithm
	DigestAlgorithm string        `json:"digestAlgorithm"`
	Sn              string        `json:"sn"`
	NotBefore       time.Time     `json:"notBefore"`
	NotAfter        time.Time     `json:"notAfter"`
	KeyUsageModel   KeyUsageModel `json:"keyUsageModel"`
	ExtKeyUsages    []int         `json:"extKeyUsages"`

	BasicConstraintsValid bool `json:"basicConstraintsValid"`
	IsCa                  bool `json:"isCa"`

	SubjectAll string `json:"subjectAll"`
	IssuerAll  string `json:"issuerAll"`

	SiaModel SiaModel `json:"siaModel"`
	// in roa, ee cert also has ip address
	CerIpAddressModel CerIpAddressModel `json:"cerIpAddressModel"`

	CrldpModel CrldpModel `json:"crldpModel"`

	// in mft/roa , eecert start/end byte:
	// eeCertByte:=OraByte[EeCertStart:EeCertEnd] eeCertByte:=MftByte[EeCertStart:EeCertEnd]
	EeCertStart uint64 `json:"eeCertStart"`
	EeCertEnd   uint64 `json:"eeCertEnd"`
}

// SIA
type SiaModel struct {
	RpkiManifest string `json:"rpkiManifest" xorm:"rpkiManifest varchar(512)"`
	RpkiNotify   string `json:"rpkiNotify" xorm:"rpkiNotify varchar(512)"`
	CaRepository string `json:"caRepository" xorm:"caRepository varchar(512)"`
	SignedObject string `json:"signedObject" xorm:"signedObject varchar(512)"`
	Critical     bool   `json:"critical"`
}

type KeyUsageModel struct {
	KeyUsage      int    `json:"keyUsage"`
	Critical      bool   `json:"critical"`
	KeyUsageValue string `json:"keyUsageValue"`
}
type CrldpModel struct {
	Crldps   []string `json:"crldps" xorm:"crldps varchar(512)"`
	Critical bool     `json:"critical"`
}
type CerIpAddressModel struct {
	CerIpAddresses []CerIpAddress `json:"cerIpAddresses"`
	Critical       bool           `json:"critical"`
}

type CerIpAddress struct {
	AddressFamily uint64 `json:"addressFamily"`
	//address prefix: 147.28.83.0/24 '
	AddressPrefix string `json:"addressPrefix" xorm:"addressPrefix varchar(512)"`
	//min address:  99.96.0.0
	Min string `json:"min" xorm:"min varchar(512)"`
	//max address:   99.105.127.255
	Max string `json:"max" xorm:"max varchar(512)"`
	//min address range from addressPrefix or min/max, in hex:  63.60.00.00'
	RangeStart string `json:"rangeStart" xorm:"rangeStart varchar(512)"`
	//max address range from addressPrefix or min/max, in hex:  63.69.7f.ff'
	RangeEnd string `json:"rangeEnd" xorm:"rangeEnd varchar(512)"`
	//min--max, such as 192.0.2.0--192.0.2.130, will convert to addressprefix range in json:{192.0.2.0/25, 192.0.2.128/31, 192.0.2.130/32}
	AddressPrefixRange string `json:"addressPrefixRange" xorm:"addressPrefixRange json"`
}
type SignerInfoModel struct {
	// must be 3
	Version int `json:"version"`
	// 2.16.840.1.101.3.4.2.1 sha-256, must be sha256
	DigestAlgorithm string `json:"digestAlgorithm"`

	// 1.2.840.113549.1.9.3 --> roa:1.2.840.113549.1.9.16.1.24  mft:1.2.840.113549.1.9.16.1.26
	ContentType string `json:"contentType"`
	// 1.2.840.113549.1.9.5
	SigningTime time.Time `json:"signingTime"`
	// 1.2.840.113549.1.9.4
	MessageDigest string `json:"messageDigest"`
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

func ParseAsaModelByAsn1(fileByte []byte, asaModel *AsaModel) (err error) {
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
		asaModel := AsaModel{}
		err = ParseAsaModelByAsn1(b, &asaModel)
		if err != nil {
			fmt.Println("ParseAsaModelByAsn1() fail:", file, err)
			continue
		}
		fmt.Println(file)
		fmt.Println(asaModel)
	}
}
