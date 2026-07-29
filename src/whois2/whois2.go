package main

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/domainr/whois"
)

// 从 WHOIS 响应里匹配 RIR
var rirRegex = regexp.MustCompile(`(?i)(RIPE|APNIC|ARIN|LACNIC|AFRINIC)`)

// DomainToRIR 使用 domainr/whois 查询域名所属 RIR
func DomainToRIR(domain string) (string, net.IP, error) {
	// 1. 解析域名到 IP
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", nil, fmt.Errorf("解析IP失败: %w", err)
	}

	var ip net.IP
	for _, candidate := range ips {
		if candidate.To4() != nil { // 优先 IPv4
			ip = candidate
			break
		}
	}
	if ip == nil {
		return "", nil, errors.New("无可用IPv4")
	}

	// 2. 用 domainr/whois 查询 IP
	req, err := whois.NewRequest(ip.String())
	if err != nil {
		return "", ip, err
	}

	resp, err := whois.DefaultClient.Fetch(req)
	if err != nil {
		return "", ip, fmt.Errorf("whois查询失败: %w", err)
	}

	// 3. 提取 RIR
	result := string(resp.Body)
	match := rirRegex.FindString(result)
	if match == "" {
		return "UNKNOWN", ip, nil
	}

	return strings.ToUpper(match), ip, nil
}

func main() {
	domains := []string{
		"baidu.com",
		"github.com",
		"bbc.com",
		"facebook.com",
		"amazon.com",
		"qq.com",
	}

	for _, d := range domains {
		rir, ip, err := DomainToRIR(d)
		if err != nil {
			fmt.Printf("%-12s 错误: %v\n", d, err)
			continue
		}
		fmt.Printf("%-12s IP=%-16s RIR=%s\n", d, ip, rir)
	}
}
