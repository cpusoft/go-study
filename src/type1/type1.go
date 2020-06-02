package main

import (
	"fmt"
	"github.com/cpusoft/goutil/jsonutil"
)

// Asn
type AsnModel struct {
	Asns     []Asn  `json:"asns"`
	AsnsP    []AsnP `json:"asnsP"`
	Critical bool   `json:"critical"`
}
type Asn struct {
	Asn uint64 `json:"asn" xorm:"asn int unsigned"`
	Min uint64 `json:"min" xorm:"min int unsigned"`
	Max uint64 `json:"max" xorm:"max int unsigned"`
}
type AsnP struct {
	Asn *uint64 `json:"asn" xorm:"asn int unsigned"`
	Min *uint64 `json:"min" xorm:"min int unsigned"`
	Max *uint64 `json:"max" xorm:"max int unsigned"`
}

func main() {
	asn1 := Asn{
		Asn: 0,
		Min: 1,
		Max: 222,
	}
	asn2 := Asn{
		Asn: 0,
		Min: 0,
		Max: 0,
	}
	asns := make([]Asn, 0)
	asns = append(asns, asn1)
	asns = append(asns, asn2)
	fmt.Println(asns)

	var i0 uint64 = 0
	var i2 uint64 = 22
	var i3 uint64 = 33
	asnP1 := AsnP{
		Asn: &i0,
		Min: &i2,
		Max: &i3,
	}
	asnP2 := AsnP{
		Asn: &i0,
		Min: &i0,
		Max: &i0,
	}

	asnsP := make([]AsnP, 0)
	asnsP = append(asnsP, asnP1)
	asnsP = append(asnsP, asnP2)
	fmt.Println(asnsP)

	asnModel := AsnModel{}
	asnModel.Asns = asns
	asnModel.AsnsP = asnsP
	fmt.Println(asnModel)

	asnModelStr := jsonutil.MarshalJson(asnModel)
	fmt.Println(asnModelStr)

	asnModel2 := AsnModel{}
	jsonutil.UnmarshalJson(asnModelStr, &asnModel2)
	fmt.Println(asnModel2)
}
