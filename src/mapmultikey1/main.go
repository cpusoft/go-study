package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type Model struct {
	Name string
}

func main() {
	model1 := &Model{"111"}
	model2 := &Model{"222"}
	model3 := &Model{"333"}
	model4 := &Model{"444"}

	mapKeys := make([]string, 0)
	mapKeys = append(mapKeys, "cer")
	mapKeys = append(mapKeys, "roa")
	mapKeys = append(mapKeys, "crl")
	mapKeys = append(mapKeys, "mft")
	c := NewConcurrentMutilMaps(mapKeys, 20000)
	c.Set("cer", "111", model1)
	c.Set("cer", "222", model2)
	c.Set("cer", "333", model3)
	c.Set("cer", "444", model4)

	c.Set("roa", "222", model2)
	c.Set("crl", "333", model3)
	c.Set("mft", "444", model4)

	m, _ := c.Get("cer", "111")
	fmt.Println(m)
	ms := c.GetMap("cer")
	fmt.Println(jsonutil.MarshalJson(ms))

	c.Remove("cer", "222")

	ms = c.GetMap("cer")
	fmt.Println(jsonutil.MarshalJson(ms))
}
