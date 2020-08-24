package main

import (
	"encoding/asn1"
	"fmt"
	_ "time"

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
	Oids Oids `asn1:"tag:17"`
}
type Oids []asn1.RawValue
type OidTime struct {
	Oid       asn1.ObjectIdentifier
	SignTimes []asn1.RawValue //[]time.Time `asn:"utc"`
}
type Sign struct {
	OidOid1 OidOid
	OidTime asn1.RawValue //OidTime
	OidOid2 asn1.RawValue //OidOid
}

func GetSignedData(value []byte) (signedData SignedData, sign Sign, err error) {
	_, err = asn1.Unmarshal(value, &signedData)
	if err != nil {
		return signedData, sign, err
	}

	raw := asn1.RawValue{}
	rest, err := asn1.Unmarshal(signedData.Sign.Bytes, &raw)
	fmt.Println("rest::", len(rest))

	oid := asn1.ObjectIdentifier{}
	rawRest, err := asn1.Unmarshal(raw.Bytes, &oid)
	raw = asn1.RawValue{}
	_, err = asn1.Unmarshal(rawRest, &raw)
	oid = asn1.ObjectIdentifier{}
	_, err = asn1.Unmarshal(raw.Bytes, &oid)
	fmt.Println("oid:", oid, err)
	return signedData, sign, nil
}
