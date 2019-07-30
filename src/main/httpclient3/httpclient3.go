package main

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
)

func main() {
	b, _ := ioutil.ReadFile(`G:\Download\cert\02DB5704B0F211E58974874FC4F9AE02.roa`)
	fmt.Println(len(b))
	resp, body, err := gorequest.New().Post("http://127.0.0.1:8080/parsecert/upload").
		Type("multipart").
		SendFile(b, "02DB5704B0F211E58974874FC4F9AE02.roa", "certfile").
		End()
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp == nil || resp.StatusCode != 200 || len(body) == 0 {
		fmt.Println(" resp == nil || resp.StatusCode != 200 || len(body) == 0 ")
		return
	}
	str := string(body)
	fmt.Println(str)
}
