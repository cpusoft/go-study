package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
	"github.com/openrdap/rdap"
)

func main() {
	req1 := &rdap.Request{
		Type:  rdap.DomainRequest,
		Query: "google.com",
	}

	client1 := &rdap.Client{}
	resp1, err := client1.Do(req1)
	fmt.Println(resp1, err)
	if domain, ok := resp1.Object.(*rdap.Domain); ok {
		fmt.Printf("Handle=%s Domain=%s\n", domain.Handle, domain.LDHName)
		fmt.Println(jsonutil.MarshalJson(domain))
	}

	req2 := &rdap.Request{
		Type:  rdap.AutnumRequest,
		Query: "2846",
	}
	client2 := &rdap.Client{}
	resp2, err := client2.Do(req2)
	fmt.Println(resp2, err)
	if autnum, ok := resp2.Object.(*rdap.Autnum); ok {
		fmt.Println(jsonutil.MarshalJson(autnum))
	}

	req3 := &rdap.Request{
		Type:  rdap.IPRequest,
		Query: "8.8.8.0/24",
	}
	client3 := &rdap.Client{}
	resp3, err := client3.Do(req3)
	fmt.Println(resp3, err)
	if ip, ok := resp3.Object.(*rdap.IPNetwork); ok {
		fmt.Println(jsonutil.MarshalJson(ip))
	}

}
