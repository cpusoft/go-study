package main

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type AsProviderAttestation struct {
	Version      int   `json:"version" asn1:"explicit,tag:0"` //default 1
	CustomerAsId int   `json:"customerAsId"`
	ProviderAss  []int `json:"providerAss"`
}

func main() {
	hexStr := `3017A003020101020207E9300C020202C0020202C1020202C2`
	by, err := hex.DecodeString(hexStr)
	fmt.Println(len(by), err)
	asProviderAttestation := AsProviderAttestation{}
	_, err = asn1.Unmarshal(by, &asProviderAttestation)
	fmt.Println(jsonutil.MarshalJson(asProviderAttestation))
	fmt.Println(asProviderAttestation, err)

}
