package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36"
	DefaultTimeout   = 5
)

func Get(urlStr string) (resp gorequest.Response, body string, errors []error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		errs := make([]error, 0)
		errs[0] = err
		return nil, "", errs
	}
	return gorequest.New().Get(urlStr).
		Timeout(DefaultTimeout*time.Minute).
		Set("User-Agent", DefaultUserAgent).
		Set("Referrer", url.Host).
		End()

}

func Post(urlStr string, postJson string) (resp gorequest.Response, body string, errors []error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		errs := make([]error, 0)
		errs[0] = err
		return nil, "", errs
	}
	return gorequest.New().Post(urlStr).
		Timeout(DefaultTimeout*time.Minute).
		Set("User-Agent", DefaultUserAgent).
		Set("Referrer", url.Host).
		Send(postJson).
		End()

}

func main() {

	urlStr := "http://localhost:8080/lookup/baidu.com"

	resp, body, errs := Get(urlStr)
	fmt.Println("resp,   ", resp)

	fmt.Println(len(body))
	fmt.Println(errs)
	/*
		urlStr = "http://localhost:8080/countries"
		content := `[{"code":"Code1", "name":"Name1"},{"code":"Code2", "name":"Name2"}]`
		//slice := make([](map[string]string), 5)

			content := make(map[string]string, 0)
			content["code"] = "CODEMAP1"
			content["name"] = "NAMEMAP1"
			content["address"] = "ADDRESS1"

		fmt.Println(content)
		resp1, body1, errs1 := Post(urlStr, content)
		fmt.Println("resp1,  ", resp1)
		fmt.Println(len(body1))
		fmt.Println(errs1)
	*/

}
