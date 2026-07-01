package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Domain string
	IPv4   []string
	IPv6   []string
	AllIPs []string
	Error  string
}

func extractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil // 去掉端口，只留域名
}
func extractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil // 去掉端口，只留域名
}
func main() {

	urls := []string{
		"https://www.example.com:5432/path?a=1",
		"https://192.168.1.1:5432/db",
		"https://[2001:db8::1]:5432/db",
		"postgresql://user:pass@host.example.com:5432/dbname", // 也兼容其他协议
	}
	domains := make([]string, 0, len(urls))
	for _, s := range urls {
		domain, err := extractDomain(s)
		if err != nil {
			fmt.Printf("ERR: %s -> %v\n", s, err)
			continue
		}
		fmt.Printf("%s -> %s\n", s, domain)
		domains = append(domains, domain)
	}

	results := resolveAll(domains, 50, 10*time.Second)

	// 打印到控制台
	for _, r := range results {
		fmt.Printf("%s\tIPv4: %v\tIPv6: %v\tAll: %v\tErr: %s\n",
			r.Domain, r.IPv4, r.IPv6, r.AllIPs, r.Error)
	}

	// 可选：写入 CSV
	writeCSV("output.csv", results)
}

func resolveAll(domains []string, workers int, timeout time.Duration) []Result {
	var wg sync.WaitGroup
	in := make(chan string, len(domains))
	out := make(chan Result, len(domains))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range in {
				out <- resolveOne(d, timeout)
			}
		}()
	}

	for _, d := range domains {
		in <- d
	}
	close(in)

	go func() {
		wg.Wait()
		close(out)
	}()

	var results []Result
	for r := range out {
		results = append(results, r)
	}
	return results
}

func resolveOne(domain string, timeout time.Duration) Result {
	r := Result{Domain: domain}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// IPv4
	if ips, err := net.DefaultResolver.LookupIP(ctx, "ip4", domain); err == nil {
		for _, ip := range ips {
			r.IPv4 = append(r.IPv4, ip.String())
			r.AllIPs = append(r.AllIPs, ip.String())
		}
	} else if !isNoRecord(err) {
		r.Error = err.Error()
	}

	// IPv6
	if ips, err := net.DefaultResolver.LookupIP(ctx, "ip6", domain); err == nil {
		for _, ip := range ips {
			r.IPv6 = append(r.IPv6, ip.String())
			r.AllIPs = append(r.AllIPs, ip.String())
		}
	} else if !isNoRecord(err) && r.Error == "" {
		r.Error = err.Error()
	}

	if len(r.AllIPs) == 0 && r.Error == "" {
		r.Error = "no records"
	}

	return r
}

func isNoRecord(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "no such host") || strings.Contains(s, "NXDOMAIN")
}

func writeCSV(filename string, results []Result) {
	f, _ := os.Create(filename)
	defer f.Close()
	w := csv.NewWriter(f)
	w.Write([]string{"domain", "ipv4", "ipv6", "all_ips", "error"})
	for _, r := range results {
		w.Write([]string{
			r.Domain,
			strings.Join(r.IPv4, ";"),
			strings.Join(r.IPv6, ";"),
			strings.Join(r.AllIPs, ";"),
			r.Error,
		})
	}
	w.Flush()
}
