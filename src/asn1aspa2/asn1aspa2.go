package main

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"

	"github.com/cpusoft/goutil/asn1util/asn1base"
	"github.com/cpusoft/goutil/asn1util/asn1node"
	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/jsonutil"
)

type AsProviderAttestation struct {
	Version int `json:"version"` //default 1
	//AddressFamilyIdentifier Afi            `json:"addressFamilyIdentifier" asn1:"optional"`
	CustomerAsId int   `json:"customerAsId"`
	ProviderAss  []int `json:"ProviderAsIds"`
}
type VersionModel struct {
	Version int
}
type ProviderAs struct {
	ProviderAsId int `json:"providerAsId"`
	//AddressFamilyIdentifier Afi `json:"addressFamilyIdentifier" asn1:"optional"`
}
type Afi []byte

type AsProviderRaw struct {
	Version    asn1.RawValue   //`asn1:"explicit,tag:5"`
	CustomerAs asn1.RawValue   //`asn1:"explicit,tag:5"`
	Proviers   []asn1.RawValue //`asn1:"explicit,tag:5"`
}

func (a Afi) MarshalText() ([]byte, error) {
	//return []byte(`[` + convert.PrintBytesOneLine(a) + `]`), nil
	s := fmt.Sprintf("%#x", a)
	return []byte(s), nil
}
func main() {
	//hexStr := `30240203033979301D3005020300FDE83009020300FDE9040200013009020300FDEA04020002`
	hexStr := `30820223A003020101020500DB6BA98430820213020412CB8E690204133428C5020414034E5002041612C2420204168B3A0C02041808030902041C6255F002041D56D703020421AD6D400204253D73B50204258B40870204271BC84202042825DCA30204288D5952020429473F2702042C818AB602042D9C3357020433356D6E02043B640A4102043DD0E31C02043F0FDB6802043F3865E802044231ECFF02044568AC25020448DDB35C02044BE5281502044D852AF702044F894F8D020450ECE9D902045113554D020456D887E0020456FE3A710204574CEE8102045759ACFB020457670DE5020457DB3AA4020461F5C8AF0204646EC452020465E91D6A02046A7538CC02046FBBBB1F0204780B8BE902047F6C060202050080BC3E1802050083CEDAAC02050086AE1B5902050087896B550205008B43507A0205008EE6CE2B02050092F6EAE40205009364BAAE0205009488CB2F020500953CFEB802050095E471610205009A6950CF0205009EC5D82E0205009F2DAAFF020500A0C41015020500A4F3DECF020500A656DCAA020500A6F2B427020500ACE44005020500B104D0A9020500B1B82C9C020500B390DA09020500BE774AD7020500BF2CC02C020500C076766B020500CC03558B020500CEBFBE08020500D0E71E01020500D419F444020500D4C1A711020500D4E57671020500D67E6C12020500D7400BD4020500DC99D146020500E26B7D04020500E28FF1F1020500F2B95D92020500FBC65D85020500FD9B3B52`
	by, err := hex.DecodeString(hexStr)
	fmt.Println(len(by), err)
	asProviderAttestation := AsProviderAttestation{}
	_, err = asn1.Unmarshal(by, &asProviderAttestation)
	fmt.Println(jsonutil.MarshalJson(asProviderAttestation))
	fmt.Println(asProviderAttestation, err)

	node, err := asn1node.ParseHex(hexStr)
	//fmt.Println(jsonutil.MarshalJson(node), err)
	versionBy := node.Nodes[0].Data
	fmt.Println("versionBy:", convert.PrintBytesOneLine(versionBy))
	var version int
	_, err = asn1.Unmarshal(versionBy, &version)
	//version, err := asn1base.ParseInt64(versionBy)
	fmt.Println("version:", version, err)

	customerAsIdBy := node.Nodes[1].Data
	fmt.Println("customerAsIdBy:", convert.PrintBytesOneLine(customerAsIdBy))
	//var customerAs int
	//_, err = asn1.Unmarshal(customerAsIdBy, &customerAs)
	customerAs, err := asn1base.ParseInt64(customerAsIdBy)
	fmt.Println("customerAs:", customerAs, err)

	//providerAssBy := node.Nodes[2].FullData
	//providerAss := make([]int, 0)
	//rest, err := asn1.Unmarshal(providerAssBy, &providerAss)
	//fmt.Println(providerAss, "\n", jsonutil.MarshalJson(providerAss), len(rest), convert.PrintBytesOneLine(rest), err)

	for i := range node.Nodes[2].Nodes {
		psBy := node.Nodes[2].Nodes[i].Data
		fmt.Println("provideras:", i, convert.PrintBytesOneLine(psBy))
		ps, err := asn1base.ParseInt64(psBy)
		fmt.Println("ps:", ps, err)
	}

}
