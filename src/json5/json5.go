package main

import (
	"errors"
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type TestJson struct {
	Name string `json:"name"`
	Err  error  `json:"err"`
}

func main() {
	t := TestJson{
		Name: "111",
		Err:  errors.New("test error"),
	}
	fmt.Println(jsonutil.MarshalJson(t))
}
