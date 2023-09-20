package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	qs := []string{"123-456", "123-", "-456", "778  - 555"}
	for _, q := range qs {
		q = strings.Replace(q, " ", "", -1)
		split := strings.Split(q, "-")
		fmt.Println(q, split, len(split))
		if len(split) == 2 {
			customerAsn := split[0]
			providerAsn := split[1]
			fmt.Println(customerAsn)
			fmt.Println(providerAsn)
			ca, err1 := strconv.Atoi(customerAsn)
			pa, err2 := strconv.Atoi(providerAsn)
			fmt.Println(ca, err1)
			fmt.Println(pa, err2)
		}
	}
}
