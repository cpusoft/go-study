package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	/*
		resp, err := http.Get("http://localhost:8080/lookup/sina.com.cn")
		if err != nil {
			// handle error
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}

		fmt.Println(string(body))
	*/
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/lookup/sohu.com.cn",
		strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/rpki-slurm")
	req.Header.Set("Cookie", "name=anny")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		// handle error
	}
	fmt.Println("Content-Type:", contentType)
	fmt.Println(string(body))

}
