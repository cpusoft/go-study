package main

import (
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type ContentInfo struct {
	ContentType ObjectIdentifier
	Seqs        []RawValue `asn1:"optional,explicit,default:0,tag:0""`
}
type RoaSignedData struct {
	Version             uint64   `json:"version"`
	AlgorithmIdentifier string   `json:"algorithmIdentifier"`
	RoaModel            RoaModel `json:"roaModel"`
}

// //////////////////////////////////////
// Roa
type RoaModel struct {
	// must be 0, but always is not in file actually
	//The version number of this version of the roa specification MUST be 0.
	Version int `json:"version"`

	Asn      int64  `json:"asn" xorm:"asn bigint"`
	Ski      string `json:"ski" xorm:"ski varchar(128)"`
	Aki      string `json:"aki" xorm:"aki varchar(128)"`
	FilePath string `json:"filePath" xorm:"filePath varchar(512)"`
	FileName string `json:"fileName" xorm:"fileName varchar(128)"`
	FileHash string `json:"fileHash" xorm:"fileHash varchar(512)"`

	//OID: 1.2.240.113549.1.9.16.1.24
	EContentType string `json:"eContentType"`

	RoaIpAddressModels []RoaIpAddressModel `json:"roaIpAddressModels"`
	SiaModel           SiaModel            `json:"siaModel"`
	AiaModel           AiaModel            `json:"aiaModel"`

	EeCertModel     EeCertModel     `json:"eeCertModel"`
	SignerInfoModel SignerInfoModel `json:"signerInfoModel"`
}

func (c RoaModel) String() string {
	m := make(map[string]interface{})
	m["ski"] = c.Ski
	m["aki"] = c.Aki
	m["asn"] = c.Asn
	m["filePath"] = c.FilePath
	m["fileName"] = c.FileName
	m["aiaModel"] = c.AiaModel.String()
	m["siaModel"] = c.SiaModel.String()
	m["len(roaIpAddressModels)"] = len(c.RoaIpAddressModels)
	return jsonutil.MarshalJson(m)
}

type RoaIpAddressModel struct {
	AddressFamily uint64 `json:"addressFamily" xorm:"addressFamily int unsigned"`
	AddressPrefix string `json:"addressPrefix" xorm:"addressPrefix varchar(512)"`
	MaxLength     uint64 `json:"maxLength" xorm:"maxLength int unsigned"`
	//min address range from addressPrefix or min/max, in hex:  63.60.00.00'
	RangeStart string `json:"rangeStart" xorm:"rangeStart varchar(512)"`
	//max address range from addressPrefix or min/max, in hex:  63.69.7f.ff'
	RangeEnd string `json:"rangeEnd" xorm:"rangeEnd varchar(512)"`
	//min--max, such as 192.0.2.0--192.0.2.130, will convert to addressprefix range in json:{192.0.2.0/25, 192.0.2.128/31, 192.0.2.130/32}
	AddressPrefixRange string `json:"addressPrefixRange" xorm:"addressPrefixRange json"`
}

func (my *RoaModel) Compare(other *RoaModel) bool {
	if my.Asn != other.Asn {
		return false
	}
	if len(my.RoaIpAddressModels) != len(other.RoaIpAddressModels) {
		return false
	}
	for i := range my.RoaIpAddressModels {
		found := false
		for j := range other.RoaIpAddressModels {
			if my.RoaIpAddressModels[i].AddressFamily == other.RoaIpAddressModels[j].AddressFamily &&
				my.RoaIpAddressModels[i].AddressPrefix == other.RoaIpAddressModels[j].AddressPrefix &&
				my.RoaIpAddressModels[i].MaxLength == other.RoaIpAddressModels[j].MaxLength {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// /////////////////////////////////////////////
// EE
// EE in CerModel, MftModel, RoaModel, to get X509 Info and aia/sia/aki/ski
// https://datatracker.ietf.org/doc/rfc6488/?include_text=1
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

/* rfc5280
KeyUsage ::= BIT STRING {
   digitalSignature        (0),
   nonRepudiation          (1),  -- recent editions of X.509 have
                              -- renamed this bit to contentCommitment
   keyEncipherment         (2),
   dataEncipherment        (3),
   keyAgreement            (4),
   keyCertSign             (5),
   cRLSign                 (6),
   encipherOnly            (7),
   decipherOnly            (8) }
*/

// https://datatracker.ietf.org/doc/rfc6488/?include_text=1
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

// AIA
type AiaModel struct {
	CaIssuers string `json:"caIssuers" xorm:"caIssuers varchar(512)"`
	Critical  bool   `json:"critical"`
}

func (c AiaModel) String() string {
	return jsonutil.MarshalJson(c.CaIssuers)
}

// SIA
type SiaModel struct {
	RpkiManifest string `json:"rpkiManifest" xorm:"rpkiManifest varchar(512)"`
	RpkiNotify   string `json:"rpkiNotify" xorm:"rpkiNotify varchar(512)"`
	CaRepository string `json:"caRepository" xorm:"caRepository varchar(512)"`
	SignedObject string `json:"signedObject" xorm:"signedObject varchar(512)"`
	Critical     bool   `json:"critical"`
}

func (c SiaModel) String() string {
	return jsonutil.MarshalJson(c)
}

type KeyUsageModel struct {
	KeyUsage      int    `json:"keyUsage"`
	Critical      bool   `json:"critical"`
	KeyUsageValue string `json:"keyUsageValue"`
}

// IPAddress
type CerIpAddressModel struct {
	CerIpAddresses []CerIpAddress `json:"cerIpAddresses"`
	Critical       bool           `json:"critical"`
}

func (c CerIpAddressModel) String() string {
	return jsonutil.MarshalJson(c.CerIpAddresses)
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
type CrldpModel struct {
	Crldps   []string `json:"crldps" xorm:"crldps varchar(512)"`
	Critical bool     `json:"critical"`
}

type Certificate struct {
	TBSCertificate     TBSCertificate
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     BitString
}

type TBSCertificate struct {
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       RawValue
	SignatureAlgorithm AlgorithmIdentifier
	Issuer             RDNSequence
	Validity           Validity
	Subject            RDNSequence
	PublicKey          PublicKeyInfo
}

type RDNSequence []RelativeDistinguishedNameSET

type RelativeDistinguishedNameSET []AttributeTypeAndValue

type AttributeTypeAndValue struct {
	Type  ObjectIdentifier
	Value any
}

type Validity struct {
	NotBefore, NotAfter time.Time
}

type TbsCertificate struct {
	Raw                RawContent
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm AlgorithmIdentifier
	Issuer             RawValue
	Validity           Validity
	Subject            RawValue
	PublicKey          PublicKeyInfo
	Extensions         []Extension `asn1:"optional,explicit,tag:3"`
}
type AlgorithmIdentifier struct {
	Algorithm ObjectIdentifier
	//Parameters RawValue `asn1:"optional"`
}

type PublicKeyInfo struct {
	Raw       RawContent
	Algorithm AlgorithmIdentifier
	PublicKey BitString
}

type Extension struct {
	Raw      RawContent
	Oid      ObjectIdentifier
	Critical bool `asn1:"optional"`
	Value    []byte
}

//////////////////////////////////////////////

// asID as in rfc6482
type RouteOriginAttestation struct {
	AsID         ASID                 `json:"asID"`
	IpAddrBlocks []ROAIPAddressFamily `json:"ipAddrBlocks"`
}
type ASID int64
type ROAIPAddressFamily struct {
	AddressFamily []byte         `json:"addressFamily"`
	Addresses     []ROAIPAddress `json:"addresses"`
}
type ROAIPAddress struct {
	Address   BitString `json:"address"`
	MaxLength int       `asn1:"optional,default:-1" json:"maxLength"`
}

/*
Version          int
DigestAlgorighms AlgorithmIdentifier
EncapContentInfo EncapsulatedContentInfo
TbsCertificate   TBSCertificate
SignerInfos      RawValue
*/

type EncapsulatedContentInfo struct {
	EContentType string     `json:"eContentType"`
	EContent     []RawValue `asn1:"optional,explicit,default:0,tag:0""`
}
type OctString []byte

type Sha256 struct {
	Oid  ObjectIdentifier
	Null RawValue
}

type OctetString []byte
type RoaOctetString struct {
	EContentType ObjectIdentifier
	OctetString  OctetString `asn1:"tag:0,explicit,optional"`
}

func main() {
	var file string
	file = `F:\share\我的坚果云\Go\common\go-study\src\asnroa1\1.roa`
	file = `F:\share\我的坚果云\Go\common\go-study\src\asnroa1\asn0.roa`
	file = `F:\share\我的坚果云\Go\common\go-study\src\asnroa1\ok.roa`
	//file = `F:\share\我的坚果云\Go\common\go-study\src\asnroa1\fail1.roa`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}

	contentInfo := ContentInfo{}
	_, err = Unmarshal(b, &contentInfo)
	if err != nil {
		fmt.Println("file:", file, err)
		return
	}
	contentTypeOid := contentInfo.ContentType.String()
	fmt.Println("contentTypeOid:", contentTypeOid)

	roaSignedData := RoaSignedData{}
	for _, seq := range contentInfo.Seqs {
		fmt.Println("seq:", jsonutil.MarshallJsonIndent(seq))

		if seq.Class == 0 && seq.Tag == 2 && !seq.IsCompound {
			// version:       version CMSVersion INTEGER 3
			roaSignedData.Version = convert.Bytes2Uint64(seq.Bytes)
		} else if seq.Class == 0 && seq.Tag == 17 && seq.IsCompound && len(seq.Bytes) < 100 {
			// digestAlgorithms DigestAlgorithmIdentifiers or signerInfos SignerInfos SET (1 elem)
			var algorithmIdentifier AlgorithmIdentifier
			_, err = Unmarshal(seq.Bytes, &algorithmIdentifier)
			if err != nil {
				fmt.Println("algorithmIdentifier fail:", err)
			} else {
				roaSignedData.AlgorithmIdentifier = algorithmIdentifier.Algorithm.String()
			}

		} else if seq.Class == 0 && seq.Tag == 16 && seq.IsCompound {
			//  encapContentInfo EncapsulatedContentInfo
			fmt.Println("\n\nRoaOctetString start-----okokok-----------------------------")
			var roaOctetString RoaOctetString
			_, err = Unmarshal(seq.FullBytes, &roaOctetString)
			if err != nil {
				fmt.Println("Unmarshal roaOctetString fail:", err)
			} else {
				fmt.Println("roaOctetString:", roaOctetString)
			}
			routeOriginAttestation := RouteOriginAttestation{}
			_, err = Unmarshal([]byte(roaOctetString.OctetString), &routeOriginAttestation)
			if err != nil {
				fmt.Println("Unmarshal routeOriginAttestation fail:", err)
			}
			fmt.Println("RoaOctetString: routeOriginAttestation", jsonutil.MarshallJsonIndent(routeOriginAttestation))
			fmt.Println("RoaOctetString end----------------------------------\n\n\n\n")

			roaSignedData.RoaModel.Version = int(routeOriginAttestation.AsID)
			roaIpAddressModels := make([]RoaIpAddressModel, 0)
			for i := range routeOriginAttestation.IpAddrBlocks {
				ipAddrBlock := routeOriginAttestation.IpAddrBlocks[i]
				addressFamily := convert.BytesToBigInt(ipAddrBlock.AddressFamily)
				fmt.Println("addressFamily:", addressFamily)
				var size int
				if addressFamily.Uint64() == 1 {
					size = 4
				} else if addressFamily.Uint64() == 2 {
					size = 16
				}

				for j := range ipAddrBlock.Addresses {

					ipAddr := make([]byte, size)
					copy(ipAddr, ipAddrBlock.Addresses[j].Address.Bytes)
					mask := net.CIDRMask(ipAddrBlock.Addresses[j].Address.BitLength, size*8)
					fmt.Println("ipAddr:", convert.PrintBytesOneLine(ipAddr),
						jsonutil.MarshalJson(ipAddr), "  mask:", mask)
					ipNet := net.IPNet{
						IP:   net.IP(ipAddr),
						Mask: mask,
					}
					maxlength := ipAddrBlock.Addresses[j].MaxLength
					fmt.Println("ipNet:", ipNet.String(), "  maxlength:", maxlength)
				}
			}
			fmt.Println("roaIpAddressModels:", jsonutil.MarshallJsonIndent(roaIpAddressModels))
		}
	}
	fmt.Println("roaSignedData:", jsonutil.MarshallJsonIndent(roaSignedData))

}
