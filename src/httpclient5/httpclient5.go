package main

import (
	"fmt"

	"github.com/cpusoft/goutil/httpclient"
)

func main() {
	url := `http://202.173.14.103:58085/allReset`
	var httpResponse HttpResponse
	err := httpclient.PostAndUnmarshalStruct(url, "", false, httpResponse)
	fmt.Println(httpResponse, err)
}

type HttpResponse struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}
