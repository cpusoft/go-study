package main

import (
	"encoding/json"
	"fmt"
)

//https://github.com/golang/go/issues/16426
// `json:"alloc"      gencodec:"required"`

type DD struct {
	Test *string `json:"test"  gencodec:"required" default:"111"`
}

//slurm head , 这里特殊处理，相当于每组只会有一个
type SlurmTarget struct {
	Asn      uint64 `json:"asn,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// filter
type PrefixFilters struct {
	Asn     uint64 `json:"asn"`
	Prefix  string `json:"prefix"`
	Comment string `json:"comment"`
}

type BgpsecFilters struct {
	Asn       uint64 `json:"asn"`
	RouterSKI string `json:"routerSKI"`
	Comment   string `json:"comment"`
}

type ValidationOutputFilters struct {
	PrefixFilters []PrefixFilters `json:"prefixFilters"`
	BgpsecFilters []BgpsecFilters `json:"bgpsecFilters"`
}

// assertion
type PrefixAssertions struct {
	Asn             uint64 `json:"asn"`
	Prefix          string `json:"prefix"`
	MaxPrefixLength int    `json:"maxPrefixLength"`
	Comment         string `json:"comment"`
}

type BgpsecAssertions struct {
	Asn       uint64 `json:"asn"`
	Comment   string `json:"comment"`
	SKI       string `json:"SKI"`
	PublicKey string `json:"publicKey"`
}

type LocallyAddedAssertions struct {
	PrefixAssertions []PrefixAssertions `json:"prefixAssertions"`
	BgpsecAssertions []BgpsecAssertions `json:"bgpsecAssertions"`
}

type Slurm struct {
	SlurmVersion            int                     `json:"slurmVersion"`
	SlurmTarget             []SlurmTarget           `json:"slurmTarget"`
	ValidationOutputFilters ValidationOutputFilters `json:"validationOutputFilters"`
	LocallyAddedAssertions  LocallyAddedAssertions  `json:"locallyAddedAssertions"`
}

func main() {
	dd := DD{}
	fmt.Println(dd)

	slurm := Slurm{}
	str := `{
       "slurmVersion": 1,
       "slurmTarget":[
         {
           "asn":"2"
         },
         {
           "asn":4
         },
         {
           "asn":3
         },
         {
           "asn":5
         },
         {
           "hostname":"rpki.example.com"
         }
       ],
       "validationOutputFilters": {
         "prefixFilters": [
           {
             "prefix": "64496",
             "comment": "All VRPs encompassed by prefix"
           },
           {
             "prefix": "192.0.2.0/24",
             "comment": "All VRPs encompassed by prefix"
           },
           {
             "asn": 64496,
             "comment": "All VRPs matching ASN"
           },
           {
             "prefix": "198.51.100.0/24",
             "asn": 64497,
             "comment": "All VRPs encompassed by prefix, matching ASN"
           }
         ],
         "bgpsecFilters": [
           {
             "asn": 64496,
             "comment": "All keys for ASN"
           },
           {
             "routerSKI": "Zm9v",
             "comment": "Key matching Router SKI"
           },
           {
             "asn": 64497,
             "routerSKI": "YmFy",
             "comment": "Key for ASN 64497 matching Router SKI"
           }
         ]
       },
       "locallyAddedAssertions": {
         "prefixAssertions": [
           {
             "asn": 64496,
             "prefix": "198.51.100.0/24",
             "comment": "My other important route"
           },
           {
             "asn": 64496,
             "prefix": "2001:DB8::/32",
             "maxPrefixLength": 48,
             "comment": "My other important de-aggregated routes"
           }
         ],
         "bgpsecAssertions": [
           {
             "asn": 64496,
             "comment" : "My known key for my important ASN",
             "SKI": "aksgnvKhg08qrg",
             "publicKey": "iuhnlkjHIOUHBjfkahvKhoi&*89qr"
           }
         ]
       }
     }
	`

	//https://tools.ietf.org/html/rfc7607
	//asn理论上不能0，需要从1开始
	err := json.Unmarshal([]byte(str), &slurm)
	if err != nil {
		fmt.Println("Para json failed: ", err)

	}
	fmt.Print(slurm)

	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	var f interface{}
	json.Unmarshal(b, &f)
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%+v", f))
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is array:")
			for i, j := range vv {
				fmt.Println(i, j)
			}
		default:
			fmt.Println(k, "dont know type")
		}
	}

}
