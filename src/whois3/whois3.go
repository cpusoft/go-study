package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func main() {

	// 查询域名
	fmt.Println("example.com---------------------------")
	whois_raw, err := whois.Whois("example.com")
	fmt.Println(whois_raw, err)
	result, err := whoisparser.Parse(whois_raw)
	if err == nil {
		fmt.Println("result", jsonutil.MarshalJson(result))
		fmt.Println("result.Domain---------------------------")
		if result.Domain != nil {
			fmt.Println(result.Domain)
		}
		fmt.Println("result.Registrar---------------------------")
		// Print the registrar name
		if result.Registrar != nil {
			fmt.Println(result.Registrar)
		}
		fmt.Println("result.Registrant---------------------------")
		// Print the registrant name
		if result.Registrant != nil {
			fmt.Println(result.Registrant)
		}
		fmt.Println("result.Administrative---------------------------")
		if result.Administrative != nil {
			fmt.Println(result.Administrative)
		}
		fmt.Println("result.Technical---------------------------")
		if result.Technical != nil {
			fmt.Println(result.Technical)
		}
		fmt.Println("result.Billing---------------------------")
		if result.Billing != nil {
			fmt.Println(result.Billing)
		}
	}
	fmt.Println("ssssssssssssssssssssssssssss")
	// 查询 IP（底层会自动找到对应 RIR 服务器）
	whois_raw, err = whois.Whois("1.1.1.1")
	fmt.Println(whois_raw, err)
	result, err = whoisparser.Parse(whois_raw)
	if err == nil {
		fmt.Println("result", jsonutil.MarshalJson(result))
	}

	// 查询 ASN
	whois_raw, err = whois.Whois("AS60614")
	fmt.Println(whois_raw, err)
	result, err = whoisparser.Parse(whois_raw)
	if err == nil {
		// Print the domain status
		fmt.Println("result", jsonutil.MarshalJson(result))
	}

}
