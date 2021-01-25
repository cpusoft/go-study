package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type CerIdModel struct {
	Id uint64 `json:"id"`
}

func main() {
	json := `
	[
	{
		"id": 19913
	},
	{
		"id": 19938
	},
	{
		"id": 22075
	}
]`
	cerIdModels := make([]CerIdModel, 0)
	jsonutil.UnmarshalJson(json, &cerIdModels)
	fmt.Println(cerIdModels)

	c1 := CerIdModel{Id: 1}
	c2 := CerIdModel{Id: 2}
	c3 := CerIdModel{Id: 3}
	c4 := CerIdModel{Id: 4}
	cerIdModels = append(cerIdModels, c1)
	cerIdModels = append(cerIdModels, c2)
	cerIdModels = append(cerIdModels, c3)
	cerIdModels = append(cerIdModels, c4)
	json = jsonutil.MarshalJson(cerIdModels)
	fmt.Println(json)

	cerIdModels2 := make([]CerIdModel, 0)
	jsonutil.UnmarshalJson(json, &cerIdModels2)
	fmt.Printf("%v\n", cerIdModels2)
}
