package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type Node5 struct {
	IpFamliyBytes []byte `json:"value"`
	Node6s        Node6  `json:"nodes"`
}
type Node6 struct {
	IpAddresses []IpAddress
}
type IpAddress struct {
	IpAddressPrefix []byte `json:"value"`
}

func main() {
	s := ` 
	{
		{
			"value": "Ag=="
		}, 
		{
			"nodes": [
				{
					"value": "ACABBnwgjA=="
				}
			]
		}
	}`
	n := Node5{}
	err := jsonutil.UnmarshalJson(s, &n)
	fmt.Println(n, err)

}
