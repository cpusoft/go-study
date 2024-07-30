package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	/*
	   client与server进行通信时 client也要对server返回数字证书进行校验
	   因为server自签证书是无效的 为了client与server正常通信
	   通过设置客户端跳过证书校验
	   TLSClientConfig:{&tls.Config{InsecureSkipVerify: true}
	   true:跳过证书校验
	*/
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://tal.apnic.net/apnic.tal")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
