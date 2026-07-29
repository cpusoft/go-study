package main

import (
	"context"
	"fmt"

	"github.com/lgforsberg/go-whois/whois"
)

func main() {
	ctx := context.Background()
	client, _ := whois.NewClient()

	// 域名 WHOIS
	result, _ := client.Query(ctx, "example.com")
	fmt.Println(result)

	// IP WHOIS（自动向 ARIN/RIPE/APNIC 等 RIR 查询）
	ipResult, _ := client.QueryIP(ctx, "8.8.8.8")
	for _, net := range ipResult.ParsedWhois.Networks {
		fmt.Println(net.Inetnum, net.Org)
	}
}
