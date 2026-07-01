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
