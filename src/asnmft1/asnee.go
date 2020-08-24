package main

import (
	"encoding/asn1"
	"fmt"
	"time"

	"github.com/cpusoft/goutil/convert"
)

type SignedData struct {
	Version     int    `asn1:"default:3"`
	Tmp1        []byte `asn1:"tag:0"`
	SignSha2561 SignSha256
	Sign        asn1.RawValue `asn1:"explicit,tag:0"`

	SignSha2562 SignSha256
	Tmp2        []byte
}

type SignSha256 struct {
	Oid  asn1.ObjectIdentifier
	Null asn1.RawValue
}
type OidOid struct {
	Oid  asn1.ObjectIdentifier
	Oids Oids
}
type Oids []interface{} //asn1.RawValue
type OidTime struct {
	Oid       asn1.ObjectIdentifier
	SignTimes []asn1.RawValue //[]time.Time `asn:"utc"`
}
type Sign struct {
	OidOid1 OidOid
	OidTime asn1.RawValue //OidTime
	OidOid2 asn1.RawValue //OidOid
}

func GetOidSetBytes(value []byte) (oid asn1.ObjectIdentifier, b []byte, err error) {
	oid, raw, err := GetOidSetValue(value)
	if err != nil {
		fmt.Println("GetOidSetBytes(): GetOidSetValue fail:", err)
		return
	}

	_, err = asn1.Unmarshal(raw.Bytes, &b)
	if err != nil {
		fmt.Println("GetOidSetBytes(): Unmarshal fail:", err)
		return
	}
	return
}

func GetOidSetOid(value []byte) (oid asn1.ObjectIdentifier, oid1 asn1.ObjectIdentifier, err error) {
	oid, raw, err := GetOidSetValue(value)
	if err != nil {
		fmt.Println("GetOidSetOid(): GetOidSetValue fail:", err)
		return
	}
	oid1 = asn1.ObjectIdentifier{}
	_, err = asn1.Unmarshal(raw.Bytes, &oid1)
	if err != nil {
		fmt.Println("GetOidSetOid(): Unmarshal fail:", err)
		return
	}
	return
}

func GetOidSetTime(value []byte) (oid asn1.ObjectIdentifier, time time.Time, err error) {
	oid, raw, err := GetOidSetValue(value)
	if err != nil {

		return
	}
	//Tm{}
	_, err = asn1.Unmarshal(raw.Bytes, &time)
	if err != nil {
		return
	}
	return
}

func GetOidSetValue(value []byte) (oid asn1.ObjectIdentifier, raw asn1.RawValue, err error) {
	oid = asn1.ObjectIdentifier{}
	rawRest, err := asn1.Unmarshal(value, &oid)
	if err != nil {
		fmt.Println("GetOidSetValue(): Unmarshal oid fail:", err)
		return
	}
	raw = asn1.RawValue{}
	_, err = asn1.Unmarshal(rawRest, &raw)
	if err != nil {
		fmt.Println("GetOidSetValue(): Unmarshal raw fail:", err)
		return
	}
	return
}

func GetSignedData(value []byte) (signedData SignedData, sign Sign, err error) {
	_, err = asn1.Unmarshal(value, &signedData)
	if err != nil {
		return signedData, sign, err
	}

	raw := asn1.RawValue{}
	rest, err := asn1.Unmarshal(signedData.Sign.Bytes, &raw)
	fmt.Println("rest::", len(rest))
	/*
		oid := asn1.ObjectIdentifier{}
		rawRest, err := asn1.Unmarshal(raw.Bytes, &oid)
		raw = asn1.RawValue{}
		_, err = asn1.Unmarshal(rawRest, &raw)
		oid = asn1.ObjectIdentifier{}
		_, err = asn1.Unmarshal(raw.Bytes, &oid)
		fmt.Println("oid:", oid, err)
	*/
	oid, oid1, err := GetOidSetOid(raw.Bytes)
	fmt.Println("oid:", oid, oid1, err)

	raw = asn1.RawValue{}
	rest, err = asn1.Unmarshal(rest, &raw)
	oid, time, err := GetOidSetTime(raw.Bytes)
	fmt.Println("time:", oid, time, err)

	raw = asn1.RawValue{}
	rest, err = asn1.Unmarshal(rest, &raw)
	oid, b, err := GetOidSetBytes(raw.Bytes)
	fmt.Println("[]byte:", oid, convert.Bytes2String(b), err)
	return signedData, sign, nil
}
