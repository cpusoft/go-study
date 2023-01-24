package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

func main() {
	s := `"41.0.0.0/8"`

	ss := make([]string, 0)
	s1 := ""
	jsonutil.UnmarshalJson(s, &ss)
	fmt.Println(ss)
	jsonutil.UnmarshalJson(s, &s1)
	fmt.Println(s1)

	s = `["137.63.0.0/16","137.64.0.0/16"]`
	ss = make([]string, 0)
	jsonutil.UnmarshalJson(s, &ss)
	fmt.Println(ss)
}
