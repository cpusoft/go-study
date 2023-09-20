package main

import (
	"fmt"
	"time"

	"github.com/cpusoft/goutil/jsonutil"
)

type Domain struct {
	N int         `json:"n"`
	P string      `json:"p"`
	T interface{} `json:"t"`
}

func main() {
	ds := make([]Domain, 0)

	d1 := Domain{N: 1, P: "p1"}
	ds = append(ds, d1)

	d2 := Domain{N: 2, P: "p2"}
	ds = append(ds, d2)

	d3 := Domain{N: 3, P: "p3"}
	ds = append(ds, d3)
	fmt.Println("before:", jsonutil.MarshalJson(ds))

	for i := range ds {
		go test(&ds[i], i)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("after:", jsonutil.MarshalJson(ds))
}

func test(d *Domain, i int) {
	fmt.Println(jsonutil.MarshalJson(d), i)
	t := time.Now()

	d.T = t.AddDate(0, 0, i)
	fmt.Println("in:", jsonutil.MarshalJson(d))
}
