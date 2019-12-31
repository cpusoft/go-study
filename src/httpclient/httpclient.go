package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	seedUrl := "https://192.168.83.139:8443/https/status/"
	body, err := client.Get(seedUrl)
	if err != nil {
		fmt.Errorf("get https://192.168.83.139:8443/ error")
		panic(err)
	}

	fmt.Printf("%s\n", body)

}
