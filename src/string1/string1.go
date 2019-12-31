package main

import (
	"fmt"
	"strings"
)

func main() {
	value := "${rpstir2::datadir}/rsyncrepo"
	start := strings.Index(value, "${")
	end := strings.Index(value, "}")
	fmt.Println(start, end)
	if start >= 0 && end > 0 && start < end {
		//${rpstir2::datadir}/rsyncrepo -->rpstir2::datadir
		replaceKey := string(value[start+len("${") : end])
		fmt.Println(replaceKey)

		replaceValue := "/root/rpki/data"
		prefix := string(value[:start])
		suffix := string(value[end+1:])
		newKey := prefix + replaceValue + suffix
		fmt.Println(newKey)
	}
	fmt.Println("no start")

}
