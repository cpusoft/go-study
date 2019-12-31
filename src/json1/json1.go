package main

import (
	"encoding/json"
	"fmt"
)

//https://github.com/golang/go/issues/16426
// `json:"alloc"      gencodec:"required"`

// filter
type PrefixFilters struct {
	Asn     JSONInt `json:"asn"`
	Prefix  string  `json:"prefix"`
	Comment string  `json:"comment"`
}

type BgpsecFilters struct {
	Asn       int64  `json:"asn,omitempty"`
	RouterSKI string `json:"routerSKI"`
	Comment   string `json:"comment"`
}

type ValidationOutputFilters struct {
	PrefixFilters []PrefixFilters `json:"prefixFilters"`
	BgpsecFilters []BgpsecFilters `json:"bgpsecFilters"`
}

// assertion
type PrefixAssertions struct {
	Asn             int64  `json:"asn,omitempty"`
	Prefix          string `json:"prefix"`
	MaxPrefixLength int    `json:"maxPrefixLength"`
	Comment         string `json:"comment"`
}

type BgpsecAssertions struct {
	Asn       int64  `json:"asn,omitempty"`
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
	ValidationOutputFilters ValidationOutputFilters `json:"validationOutputFilters"`
	LocallyAddedAssertions  LocallyAddedAssertions  `json:"locallyAddedAssertions"`
}
type JSONInt struct {
	Value    int
	IsNotNil bool
}

func (i *JSONInt) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.IsNotNil = true

	// The key isn't set to null
	var temp int
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	return nil
}

func main() {

	slurm := Slurm{}
	str := `{
       "slurmVersion": 1,
       "validationOutputFilters": {
         "prefixFilters": [
           {
             "prefix": "64496",
             "comment": "no asn"
           },
           {
             "asn":0,
             "prefix": "192.0.2.0/24",
             "comment": "asn is 0"
           },
           {
             "asn":1,
             "prefix": "192.0.2.0/24",
             "comment": "asn is 1"
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
	/*
		var personFromJSON interface{}

		decoder := json.NewDecoder(bytes.NewReader(str))
		decoder.UseNumber()
		decoder.
			decoder.Decode(&personFromJSON)
	*/
	for _, one := range slurm.ValidationOutputFilters.PrefixFilters {
		fmt.Println(one)
		fmt.Println(one.Asn)

	}
	tmp, _ := (json.Marshal(slurm))
	fmt.Println(string(tmp))

}
