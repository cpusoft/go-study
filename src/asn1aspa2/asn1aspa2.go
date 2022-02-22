package main

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type AsProviderAttestation struct {
	CustomerAsId  int            `json:"customerAsId"`
	ProviderAsIds []ProviderAsId `json:"ProviderAsIds"`
}
type ProviderAsId struct {
	ProviderAsId            int `json:"providerAsId"`
	AddressFamilyIdentifier Afi `json:"addressFamilyIdentifier" asn1:"optional"`
}
type Afi []byte

func (a Afi) MarshalText() ([]byte, error) {
	//return []byte(`[` + convert.PrintBytesOneLine(a) + `]`), nil
	s := fmt.Sprintf("%#x", a)
	return []byte(s), nil
}
func main() {
	hexStr := `30240203033979301D3005020300FDE83009020300FDE9040200013009020300FDEA04020002`
	by, err := hex.DecodeString(hexStr)
	fmt.Println(len(by), err)
	asProviderAttestation := AsProviderAttestation{}
	_, err = asn1.Unmarshal(by, &asProviderAttestation)
	fmt.Println(jsonutil.MarshalJson(asProviderAttestation))
	fmt.Println(asProviderAttestation, err)
	/*
		node, err := asn1node.ParseHex(hexStr)
		fmt.Println(jsonutil.MarshalJson(node), err)

		customerAsIdBy := node.Nodes[0].Data
		fmt.Println(convert.PrintBytesOneLine(customerAsIdBy))
		customerAsId, err := asn1base.ParseInt64(customerAsIdBy)
		fmt.Println(customerAsId, err)

		providerAssBy := node.Nodes[1].FullData
		providerAss := make([]ProviderAs, 0)
		rest, err := asn1.Unmarshal(providerAssBy, &providerAss)
		fmt.Println(providerAss, "\n", jsonutil.MarshalJson(providerAss), len(rest), convert.PrintBytesOneLine(rest), err)

		for i := range node.Nodes[1].Nodes {
			psBy := node.Nodes[1].Nodes[i].FullData
			fmt.Println(i, convert.PrintBytesOneLine(psBy))
			providerAs := ProviderAs{}
			rest, err := asn1.Unmarshal(psBy, &providerAs)
			fmt.Println(jsonutil.MarshalJson(providerAs), len(rest), convert.PrintBytesOneLine(rest), err)
		}
	*/
	/*
		by, err := hex.DecodeString(hexStr)
		fmt.Println(len(by), err)

		asProviderAttestations := make([]int, 0)
		rest, err := asn1.Unmarshal(by, &asProviderAttestations)
		fmt.Println(jsonutil.MarshalJson(asProviderAttestations), len(rest), convert.PrintBytesOneLine(rest), err)
	*/
}
