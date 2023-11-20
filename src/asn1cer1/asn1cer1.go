package main

import (
	"encoding/asn1"
	"fmt"
	"reflect"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type Certificate struct {
	TBSCertificate TBSCertificate

	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type TBSCertificate struct {
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       asn1.RawValue
	SignatureAlgorithm AlgorithmIdentifier
	Issuer             RDNSequence
	Validity           Validity
	Subject            RDNSequence
	PublicKey          PublicKeyInfo

	CerRawValue asn1.RawValue
}

type AlgorithmIdentifier struct {
	Algorithm asn1.ObjectIdentifier
}

type RDNSequence []RelativeDistinguishedNameSET

type RelativeDistinguishedNameSET []AttributeTypeAndValue

type AttributeTypeAndValue struct {
	Type  asn1.ObjectIdentifier
	Value interface{}
}

type Validity struct {
	NotBefore time.Time `asn1:"generalized"`
	NotAfter  time.Time `asn1:"generalized"`
}

type PublicKeyInfo struct {
	Algorithm AlgorithmIdentifier
	PublicKey asn1.BitString
}

type CerParse struct {
	SubjectKeyIdentifier   ObjectIdentifierAndRawValue
	AuthorityKeyIdentifier ObjectIdentifierAndRawValue
	KeyUsage               ObjectIdentifierAndBoolAndRawValue
	BasicConstraints       ObjectIdentifierAndBoolAndRawValue
	CRLDistributionPoints  ObjectIdentifierAndRawValue
	AuthorityInfoAccess    ObjectIdentifierAndRawValue
	CertificatePolicies    ObjectIdentifierAndBoolAndRawValue
	SubjectInfoAccess      ObjectIdentifierAndRawValue
	IpAddrBlocks           ObjectIdentifierAndBoolAndRawValue
}

type ObjectIdentifierAndRawValue struct {
	Type     asn1.ObjectIdentifier
	RawValue asn1.RawValue
}
type ObjectIdentifierAndBoolAndRawValue struct {
	Type     asn1.ObjectIdentifier
	Bool     bool
	RawValue asn1.RawValue
}

func Asn1ParseFromRawValue(rawValue *asn1.RawValue, v interface{}) {
	s := convert.PrintBytes(rawValue.Bytes, 8)
	fmt.Println("Asn1ParseFromRawValue():", s)
	cerParse := CerParse{}
	fmt.Println("Asn1ParseFromRawValue()  cerParse:", cerParse)
	asn1.Unmarshal(rawValue.Bytes, &cerParse)
	fmt.Println("Asn1ParseFromRawValue():", jsonutil.MarshallJsonIndent(cerParse))
}
func Asn1ParseReflectFromRawValue(rawValue *asn1.RawValue, v interface{}) (r interface{}, err error) {
	//s := convert.PrintBytes(rawValue.Bytes, 8)
	//fmt.Println("Asn1ParseReflectFromRawValue():s:", s)
	rt := reflect.TypeOf(v)
	fmt.Printf("Asn1ParseReflectFromRawValue():rt:%v,%T\n", jsonutil.MarshallJsonIndent(rt), rt)
	asn1.Unmarshal(rawValue.Bytes, &rt)
	fmt.Println("Asn1ParseFromRawValue():rt:", jsonutil.MarshallJsonIndent(rt))

	rt = rt.Elem()
	fmt.Printf("Asn1ParseReflectFromRawValue()  new rt %v,%T\n", jsonutil.MarshallJsonIndent(rt), rt)
	asn1.Unmarshal(rawValue.Bytes, &rt)
	fmt.Println("Asn1ParseFromRawValue():new rt:", jsonutil.MarshallJsonIndent(rt))

	nt := reflect.New(rt) // 调用反射创建对象
	fmt.Printf("Asn1ParseReflectFromRawValue()  nt: %v,%T\n", jsonutil.MarshallJsonIndent(nt), nt)
	asn1.Unmarshal(rawValue.Bytes, &nt)
	fmt.Println("Asn1ParseFromRawValue():nt:", jsonutil.MarshallJsonIndent(nt))

	nt = nt.Elem()
	fmt.Printf("Asn1ParseReflectFromRawValue()  new nt:%v,%T\n", jsonutil.MarshallJsonIndent(nt), nt)
	asn1.Unmarshal(rawValue.Bytes, &nt)
	fmt.Println("Asn1ParseFromRawValue():new nt", jsonutil.MarshallJsonIndent(nt))
	return nt, nil
}

/*
Go type                | ASN.1 universal tag
-----------------------|--------------------
bool                   | BOOLEAN
All int and uint types | INTEGER
*big.Int               | INTEGER
string                 | OCTET STRING
[]byte                 | OCTET STRING
asn1.Oid               | OBJECT INDETIFIER
asn1.Null              | NULL
Any array or slice     | SEQUENCE OF
Any struct             | SEQUENCE
*/
func main() {
	file := `E:\Go\go-study\src\asncer1\0.cer`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}
	fmt.Println(len(b))
	s := convert.PrintBytes(b, 8)
	fmt.Println(s)

	//var roaAllParse asn1.RawValue
	certificate := Certificate{}

	//roaAllParse := make([]RoaAllParse, 0)
	//roaAllParse := RoaAllParse{}
	asn1.Unmarshal(b, &certificate)
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate))
	/*
		v := CerParse{}
		Asn1ParseFromRawValue(&certificate.TBSCertificate.CerRawValue, &v)
		fmt.Println("Asn1ParseFromRawValue cerParse:", jsonutil.MarshallJsonIndent(v))
	*/
	rawValue := certificate.TBSCertificate.CerRawValue
	v := CerParse{}
	asn1.Unmarshal(rawValue.Bytes, &v)
	fmt.Println("Asn1ParseFromRawValue():v", jsonutil.MarshallJsonIndent(v))

	e := reflect.New(reflect.TypeOf(CerParse{})).Elem()
	i := e.Interface()
	fmt.Println("Asn1ParseFromRawValue():e,i:", e, "--------", i)
	c := i.(CerParse)
	fmt.Println("Asn1ParseFromRawValue():c:", c)
	asn1.Unmarshal(rawValue.Bytes, &c)
	fmt.Println("Asn1ParseFromRawValue():c:", jsonutil.MarshallJsonIndent(c))
	/*
		rt := reflect.TypeOf(v)
		fmt.Printf("Asn1ParseReflectFromRawValue():rt:%v,%T\n", rt, rt)

		asn1.Unmarshal(rawValue.Bytes, &rt)
		fmt.Println("Asn1ParseFromRawValue():rt:", rt)

		nt := reflect.New(rt) // 调用反射创建对象
		fmt.Printf("Asn1ParseReflectFromRawValue()  nt: %v,%T\n", nt, nt)
		asn1.Unmarshal(rawValue.Bytes, &nt)
		fmt.Println("Asn1ParseFromRawValue():nt:", nt)

		nt = nt.Elem()
		fmt.Printf("Asn1ParseReflectFromRawValue()  new nt:%v,%T\n", nt, nt)
		asn1.Unmarshal(rawValue.Bytes, nt)
		fmt.Println("Asn1ParseFromRawValue():new nt", nt)
	*/
	//	cerParse, err := Asn1ParseReflectFromRawValue(&certificate.TBSCertificate.CerRawValue, (*CerParse)(nil))
	//	fmt.Println("Asn1ParseReflectFromRawValue cerParse:", jsonutil.MarshallJsonIndent(cerParse), err)
	/*
		// ok
		b = certificate.TBSCertificate.CerRawValue.Bytes
		s = convert.PrintBytes(b, 8)
		fmt.Println(s)
		cerParse = CerParse{}
		asn1.Unmarshal(b, &cerParse)
		fmt.Println("cerParse:", jsonutil.MarshallJsonIndent(cerParse))
	*/
	/*
		b = cerParse.SubjectKeyIdentifier.V
		s = convert.PrintBytes(b, 8)
		fmt.Println(s)
		octectString := make([]byte, 0)
		asn1.Unmarshal(b, &octectString)
		fmt.Println("octectString:", jsonutil.MarshallJsonIndent(octectString))
		s = convert.PrintBytes(octectString, 8)
		fmt.Println(s)
	*/
	/*
		b = octectString.Bytes
		s = convert.PrintBytes(b, 8)
		fmt.Println(s)
	*/
}
