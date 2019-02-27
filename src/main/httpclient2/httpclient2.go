package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

func Get(urlStr string) (gorequest.Response, string, []error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		errs := make([]error, 0)
		errs[0] = err
		return nil, "", errs
	}
	return gorequest.New().Get(urlStr).
		Timeout(5*time.Minute).
		Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.90 Safari/537.36").
		Set("Referrer", url.Host).
		End()

}

func main() {

	urlStr := "http://localhost:8080/lookup/baidu.com"
	/*
		resp, body, errs := Get(urlStr)
		fmt.Println("resp,   ", resp)

		fmt.Println(len(body))
		fmt.Println(errs)
	*/

	urlStr = "http://localhost:8080/countries"
	body := `[{"Code":"Code1", "Name":"Name1"},{"Code":"Code2", "Name":"Name2"}]`
	resp1, body1, errs1 := gorequest.New().Post(urlStr).
		Timeout(5*time.Minute).
		Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.90 Safari/537.36").
		Send(body).
		End()
	fmt.Println("resp1,  ", resp1)
	fmt.Println(len(body1))
	fmt.Println(errs1)

}
