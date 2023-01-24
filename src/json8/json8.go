package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
	"github.com/guregu/null"
)

func main() {
	t := `[
			{
				"origin": "example.com.",
				"rrName": "test1",
				"rrType": "A",
				"rrClass": "NONE",
				"rrTtl": 0,
				"rrData": "1.1.1.1"
			},
			{
				"origin": "example.com.",
				"rrName": "test1",
				"rrType": "TXT",
				"rrClass": "IN",
				"rrTtl": 0,
				"rrData": "v=spf1 include:spf.mail.qq.com ip4:203.99.30.50 ~all"
			}
		]`
	rr := make([]*RrModel, 0)
	err := jsonutil.UnmarshalJson(t, &rr)
	fmt.Println(rr, err)
	fmt.Println(jsonutil.MarshalJson(rr))

}

type RrModel struct {
	RrFullDomain string `json:"rrFullDomain"` // lower: rrName+"."+Origin[-"."]

	RrType  string `json:"rrType"`  // upper
	RrClass string `json:"rrClass"` // upper
	// null.NewInt(0, false) or null.NewInt(i64, true)
	RrTtl  null.Int `json:"rrTtl"`
	RrData string   `json:"rrData"`
}
