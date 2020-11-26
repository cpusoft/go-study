package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
	"github.com/guregu/null"
)

type TransferHead struct {
	Operate       string `json:"operate"`
	Time          string `json:"time"`
	FromVersion   uint64 `json:"fromversion"`
	ToVersion     uint64 `json:"toversion"`
	Uuid          string `json:"uuid"`
	Sha256withRSA string `json:"sha256withRSA"`
}

type TransferModel struct {
	Head TransferHead `json:"head"`
	Data Slurm        `json:"data"`
}
type PrefixAssertions struct {
	Asn             null.Int `json:"asn"`
	Prefix          string   `json:"prefix"`
	MaxPrefixLength uint64   `json:"maxPrefixLength"`
	Comment         string   `json:"comment"`
}

// set asn==-1 means asn is empty
type BgpsecAssertions struct {
	Asn             null.Int `json:"asn"`
	Comment         string   `json:"comment"`
	SKI             string   `json:"SKI"`
	RouterPublicKey string   `json:"RouterPublicKey"`
}

type LocallyAddedAssertions struct {
	PrefixAssertions []PrefixAssertions `json:"prefixAssertions"`
	BgpsecAssertions []BgpsecAssertions `json:"bgpsecAssertions"`
}
type PrefixFilters struct {
	Asn             null.Int `json:"asn"`
	Prefix          string   `json:"prefix"`
	FormatPrefix    string   `json:"-"`
	PrefixLength    uint64   `json:"-"`
	MaxPrefixLength uint64   `json:"-"`
	Comment         string   `json:"comment"`
}

// set asn==-1 means asn is empty
type BgpsecFilters struct {
	Asn     null.Int `json:"asn"`
	SKI     string   `json:"SKI"`
	Comment string   `json:"comment"`
}

type ValidationOutputFilters struct {
	PrefixFilters []PrefixFilters `json:"prefixFilters"`
	BgpsecFilters []BgpsecFilters `json:"bgpsecFilters"`
}
type Slurm struct {
	SlurmVersion            int                     `json:"slurmVersion"`
	ValidationOutputFilters ValidationOutputFilters `json:"validationOutputFilters"`
	LocallyAddedAssertions  LocallyAddedAssertions  `json:"locallyAddedAssertions"`
}

func main() {
	body := `{"head":{"operate":"all","time":"2020-10-09 22:17:24","fromversion":0,"toversion":13,"uuid":"4bf8e10f-49d1-420b-927c-46a7b2721309","sha256withRSA":"a72422398dd1b87da73bb182877f8ee4e490445a"},"data":{"slurmVersion":1,"validationOutputFilters":{"prefixFilters":null,"bgpsecFilters":null},"locallyAddedAssertions":{"prefixAssertions":[{"asn":11,"prefix":"101.101.96/22","maxPrefixLength":22,"comment":""},{"asn":null,"prefix":"101.97.43/24","maxPrefixLength":24,"comment":""}],"bgpsecAssertions":null}}}`
	transferModel := TransferModel{}
	err := jsonutil.UnmarshalJson(body, &transferModel)
	fmt.Println(transferModel, err)
	body = jsonutil.MarshalJson(transferModel)
	fmt.Println(body)

	body = `{"asn":11}`
	slurmAsnModel := null.Int{}
	err = jsonutil.UnmarshalJson(body, &slurmAsnModel)
	fmt.Println(slurmAsnModel, err)
}
