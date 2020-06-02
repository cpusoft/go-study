package main

import (
	"fmt"
	jsonutil "github.com/cpusoft/goutil/jsonutil"
)

type Data struct {
	AcceptTime uint64 `json:"acceptTime"`
	AllotTime  uint64 `json:"allotTime"`
}
type DD struct {
	Data      Data   `json:"data"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type DD1 struct {
	DD
	ErrorCode string `json:"-"`
}

func main() {
	data := Data{AcceptTime: 1, AllotTime: 2}
	dd := DD{Data: data, Message: "test", ErrorCode: "error"}
	test := jsonutil.MarshalJson(dd)
	fmt.Println(test)

	dd1 := DD{Data: data, Message: "test1", ErrorCode: "error1"}
	test1 := jsonutil.MarshalJson(dd1)
	fmt.Println(test1)

	var ii interface{}
	ii = dd
	var ii1 interface{}
	ii1 = dd1

	test = jsonutil.MarshalJson(ii)
	fmt.Println(test)
	test1 = jsonutil.MarshalJson(ii1)
	fmt.Println(test1)

}
