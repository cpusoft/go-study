package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/cpusoft/goutil/jsonutil"
	"github.com/guregu/null"
)

type HroaFilters struct {
	Asn               null.Int `json:"asn"`
	SubtreeIdentifier *big.Int `json:"subtree_identifier"`
	EncodedSubtree    null.Int `json:"encoded_subtree"`
}

func main() {
	b := big.NewInt(0)
	b, _ = b.SetString("51240900000000000222", 10)
	h := HroaFilters{
		Asn:               null.IntFrom(int64(99)),
		SubtreeIdentifier: b,
		EncodedSubtree:    null.IntFrom(int64(199)),
	}
	s := jsonutil.MarshalJson(h)
	fmt.Println(s)
	v := b.FillBytes(make([]byte, 32))
	fmt.Println(hex.EncodeToString(v))
	b2 := big.NewInt(0)
	b2.SetBytes(v)
	fmt.Println(b2)

	s = `{"asn":99,"subtree_identifier":512409557484068888888832,"encoded_subtree":199}`
	err := jsonutil.UnmarshalJson(s, &h)
	fmt.Println(h, err)

}
