package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	skip := "/static/*"
	test1 := strings.HasSuffix(skip, "*")
	fmt.Println(test1)

	url := "/static/js/commons.js"
	//url := "/static/"
	//pattern := `^` + skip
	//reqPath:     skipUrl: /static/*
	reg := regexp.MustCompile(skip)
	test := reg.MatchString(url)
	fmt.Println(test)

}
