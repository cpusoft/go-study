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
func main() {
	// select '"' || url || '",' from lab_rpki_sync_url where "syncStyle" = 'rrdp' order by url
	urls := []string{
		"https://akane.maru.co.jp/oshirase.xml",
		"https://alias.rpki-measurement.site/rrdp/notification.xml",
		"https://ca.nat.moe/rrdp/notification.xml",
		"https://ca.rg.net/rrdp/notify.xml",
		"https://chloe.sobornost.net/rpki/news.xml",
		"https://cloudie.rpki.app/rrdp/notification.xml",
		"https://feo.tla.org/rrdp/notification.xml",
		"https://krill.47272.net/rrdp/notification.xml",
		"https://krill.accuristechnologies.ca/rrdp/notification.xml",
		"https://krill.byhill.de/rrdp/notification.xml",
		"https://krill.ca-bc-01.ssmidge.xyz/rrdp/notification.xml",
		"https://krill.cxy.hk/rrdp/notification.xml",
		"https://krill.dacs.utwente.nl/rrdp/notification.xml",
		"https://krill.immarket.space/rrdp/notification.xml",
		"https://krill.ipgua.com:3030/rrdp/notification.xml",
		"https://krill.peering.ee.columbia.edu/rrdp/notification.xml",
		"https://krill.peravix.group/rrdp/notification.xml",
		"https://krill.rg.net/rrdp/notification.xml",
		"https://krill.rpki-measurement.site/rrdp/notification.xml",
		"https://krill.signalx.cloud/rrdp/notification.xml",
		"https://krill.starlamp.su/rrdp/notification.xml",
		"https://krill.sy5.nsw.au.infininet.com.au/rrdp/notification.xml",
		"https://magellan.ipxo.com/rrdp/notification.xml",
		"https://orca.rg.net/rrdp/notification.xml",
		"https://oto.wakuwaku.ne.jp/pki/oshirase.xml",
		"https://pub.krill.ausra.cloud/rrdp/notification.xml",
		"https://repodepot.wildtky.com/rrdp/notification.xml",
		"https://repo.kagl.me/rpki/notification.xml",
		"https://repo.rpki.space/rrdp/notification.xml",
		"https://rki.plasmanodes.com/rrdp/notification.xml",
		"https://rov-measurements.nlnetlabs.net/rrdp/notification.xml",
		"https://rpki-01.pdxnet.uk/rrdp/notification.xml",
		"https://rpki.admin.freerangecloud.com/rrdp/notification.xml",
		"https://rpki.apernet.io/rrdp/notification.xml",
		"https://rpki.as207960.net/rrdp/notification.xml",
		"https://rpki.as215605.net/rrdp/notification.xml",
		"https://rpki.athene-center.net/rrdp/notification.xml",
		"https://rpki.axivora.net/rrdp/notification.xml",
		"https://rpki.blw.moe/rrdp/notification.xml",
		"https://rpki.cc/rrdp/notification.xml",
		"https://rpki.cernet.edu.cn/rrdp/notification.xml",
		"https://rpki.cernet.net/rrdp/notification.xml",
		"https://rpki.co/rrdp/notification.xml",
		"https://rpki.fengrui.link/rrdp/notification.xml",
		"https://rpki.fiti.net.cn/rrdp/notification.xml",
		"https://rpki.folf.systems/rrdp/notification.xml",
		"https://rpki.gns.net.br/rrdp/notification.xml",
		"https://rpki.gxt.network/rrdp/notification.xml",
		"https://rpki.komorebi.network:3030/rrdp/notification.xml",
		"https://rpki.leitecastro.com/notification.xml",
		"https://rpki.luys.cloud/rrdp/notification.xml",
		"https://rpki.miralium.net/rrdp/notification.xml",
		"https://rpki.multacom.com/rrdp/notification.xml",
		"https://rpki.nellicus.net/rrdp/notification.xml",
		"https://rpki.owl.net/rrdp/notification.xml",
		"https://rpki-pp.com/rrdp/notification.xml",
		"https://rpki-publication.haruue.net/rrdp/notification.xml",
		"https://rpki.pudu.be/rrdp/notification.xml",
		"https://rpki.qs.nu/rrdp/notification.xml",
		"https://rpki.rand.apnic.net/rrdp/notification.xml",
		"https://rpki-repo.canops.org/rrdp/notification.xml",
		"https://rpki-repo.nex3.com.br/rrdp/notification.xml",
		"https://rpki-repo.registro.br/rrdp/notification.xml",
		"https://rpki-repository.nic.ad.jp/rrdp/ap/notification.xml",
		"https://rpki.roa.net/rrdp/notification.xml",
		"https://rpki-rrdp.idnic.net/rrdp/notification.xml",
		"https://rpki-rrdp.mnihyc.com/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/08c2f264-23f9-49fb-9d43-f8b50bec9261/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/16f1ffee-7461-4674-bb05-fddefa9a02c6/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/20aa329b-fc52-4c61-bf53-09725c042942/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/2df51cd2-e6af-493a-a88a-3221d01f7d90/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/2f059a21-d41b-4846-b7ae-7ea38c32fd4c/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/42582c67-dd3f-4bc5-ba60-e97e552c6e35/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/517f3ed7-58b5-4796-be37-14d62e48f056/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/54602fb0-a9d4-4f9f-b0ca-be2a139ea92b/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/602a26e5-4a9e-4e5e-89f0-ef891490d9c9/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/708aafaf-00b4-485b-854c-0b32ca30f57b/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/71e5236f-c6f1-4928-a1b9-8def09c06085/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/967a255c-d680-42d3-9ec3-ecb3f9da088c/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/a841823c-a10d-477c-bfdf-4086f0b1594c/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/b3f6b688-cff4-402f-97d5-02f6f1886b7e/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/b8a1dd25-c313-4f25-ac21-bf55514d9c7d/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/bd48a1fa-3471-4ab2-8508-ad36b96813e4/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/c3cd7c24-12cb-4abc-8fd2-5e2bcbb85ae6/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/c474d778-43cb-4c30-ad6a-39968cbc94bc/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/db9a372a-09bc-4a32-bfe4-8c48e5dbd219/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/dba8f01c-9669-44a3-ac6e-db2edb099b84/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/dfd7f6d3-e6e9-4987-9ae7-d052c5353898/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/e72d8db0-4728-4fc1-bdd8-471129866362/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/e7518af5-a343-428d-bf78-f982b6e60505/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/f703696e-e47b-4c20-bd93-6f80904e42d2/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/fd36ad64-200f-4064-84a8-3c7cc91cdece/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/fe3737fb-095d-444c-92f4-3f7221fb544c/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/ff9fa84e-9783-4a0b-a58d-6dc8e2433d33/notification.xml",
		"https://rpki-rrdp.warpnet.xyz/notification.xml",
		"https://rpki.sailx.co/rrdp/notification.xml",
		"https://rpki.ssmidge.xyz/rrdp/notification.xml",
		"https://rpki.sunoaki.net/rrdp/notification.xml",
		"https://rpki.tools.westconnect.ca/rrdp/notification.xml",
		"https://rpki.xa.wiki/rrdp/notification.xml",
		"https://rpki.xindi.eu/rrdp/notification.xml",
		"https://rpki.zappiehost.com/rrdp/notification.xml",
		"https://rrdp.afrinic.net/notification.xml",
		"https://rrdp.apnic.net/notification.xml",
		"https://rrdp.arin.net/notification.xml",
		"https://rrdp-as0.apnic.net/notification.xml",
		"https://rrdp.as214749.net/rrdp/notification.xml",
		"https://rrdp.krill.nlnetlabs.nl/notification.xml",
		"https://rrdp.lacnic.net/rrdpas0/notification.xml",
		"https://rrdp.lacnic.net/rrdp/notification.xml",
		"https://rrdp.paas.rpki.ripe.net/notification.xml",
		"https://rrdp.prefixlogic.com/rrdp/notification.xml",
		"https://rrdp.ripe.net/notification.xml",
		"https://rrdp.rp.ki/notification.xml",
		"https://rrdp.rpki.tianhai.link/rrdp/notification.xml",
		"https://rrdp-rps.arin.net/notification.xml",
		"https://rrdp-rps.cnnic.cn/rrdp/notification.xml",
		"https://rrdp.sub.apnic.net/notification.xml",
		"https://rrdp.twnic.tw/rrdp/notification.xml",
		"https://sakuya.nat.moe/rrdp/notification.xml",
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
