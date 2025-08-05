package main

import (
	"fmt"
	"net"
	"net/url"
)

func main() {
	urls := []string{
		"https://rrdp.apnic.net/notification.xml",
		"https://rrdp.ripe.net/notification.xml",
		"https://rrdp.lacnic.net/rrdpas0/notification.xml",
		"https://rrdp.lacnic.net/rrdp/notification.xml",
		"https://rrdp.arin.net/notification.xml",
		"https://rrdp.afrinic.net/notification.xml",
		"https://rrdp-as0.apnic.net/notification.xml",
		//	"https://rpki.fzca.com/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/e7518af5-a343-428d-bf78-f982b6e60505/notification.xml",
		"https://rrdp.paas.rpki.ripe.net/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/dfd7f6d3-e6e9-4987-9ae7-d052c5353898/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/71e5236f-c6f1-4928-a1b9-8def09c06085/notification.xml",
		"https://rrdp.sub.apnic.net/notification.xml",
		"https://rpki.tools.westconnect.ca/rrdp/notification.xml",
		"https://rpki.multacom.com/rrdp/notification.xml",
		"https://krill.accuristechnologies.ca/rrdp/notification.xml",
		"https://repodepot.wildtky.com/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/20aa329b-fc52-4c61-bf53-09725c042942/notification.xml",
		"https://cloudie.rpki.app/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/602a26e5-4a9e-4e5e-89f0-ef891490d9c9/notification.xml",
		"https://repo.kagl.me/rpki/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/e72d8db0-4728-4fc1-bdd8-471129866362/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/b8a1dd25-c313-4f25-ac21-bf55514d9c7d/notification.xml",
		"https://rpki.roa.net/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/42582c67-dd3f-4bc5-ba60-e97e552c6e35/notification.xml",
		"https://rrdp-rps.arin.net/notification.xml",
		"https://dev.tw/rpki/notification.xml",
		"https://rpki.sailx.co/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/ff9fa84e-9783-4a0b-a58d-6dc8e2433d33/notification.xml",
		"https://oto.wakuwaku.ne.jp/pki/oshirase.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/b3f6b688-cff4-402f-97d5-02f6f1886b7e/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/517f3ed7-58b5-4796-be37-14d62e48f056/notification.xml",
		"https://rpki.admin.freerangecloud.com/rrdp/notification.xml",
		"https://rpki.miralium.net/rrdp/notification.xml",
		"https://rpki.zappiehost.com/rrdp/notification.xml",
		"https://rpki.xa.wiki/rrdp/notification.xml",
		"https://rpki.komorebi.network:3030/rrdp/notification.xml",
		"https://rpki.cc/rrdp/notification.xml",
		"https://repo.rpki.space/rrdp/notification.xml",
		"https://rpki.sn-p.io/rrdp/notification.xml",
		"https://rpki-publication.haruue.net/rrdp/notification.xml",
		"https://rpki-01.pdxnet.uk/rrdp/notification.xml",
		"https://krill.immarket.space/rrdp/notification.xml",
		"https://magellan.ipxo.com/rrdp/notification.xml",
		"https://rpki.as207960.net/rrdp/notification.xml",
		"https://rpki.pudu.be/rrdp/notification.xml",
		"https://rpki.uz/rrdp/notification.xml",
		"https://rpki.co/rrdp/notification.xml",
		"https://rpki.qs.nu/rrdp/notification.xml",
		"https://krill.stonham.uk/rrdp/notification.xml",
		"https://rpki.ssmidge.xyz/rrdp/notification.xml",
		"https://krill.ca-bc-01.ssmidge.xyz/rrdp/notification.xml",
		"https://krill.uta.ng:3030/rrdp/notification.xml",
		"https://rpki-repo.registro.br/rrdp/notification.xml",
		"https://repo-rpki.idnic.net/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/a841823c-a10d-477c-bfdf-4086f0b1594c/notification.xml",
		"https://rpki.netiface.net/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/08c2f264-23f9-49fb-9d43-f8b50bec9261/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/54602fb0-a9d4-4f9f-b0ca-be2a139ea92b/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/dba8f01c-9669-44a3-ac6e-db2edb099b84/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/967a255c-d680-42d3-9ec3-ecb3f9da088c/notification.xml",
		"https://ca.nat.moe/rrdp/notification.xml",
		"https://rrdp.rpki.co/rrdp/notification.xml",
		"https://rpki.folf.systems/rrdp/notification.xml",
		"https://rpki.sunoaki.net/rrdp/notification.xml",
		"https://rrdp.rpki.tianhai.link/rrdp/notification.xml",
		"https://rpki.athene-center.net/rrdp/notification.xml",
		"https://rpki.owl.net/rrdp/notification.xml",
		"https://0.sb/rrdp/notification.xml",
		"https://rrdp.rp.ki/notification.xml",
		"https://rpki-repository.nic.ad.jp/rrdp/ap/notification.xml",
		"https://rpki.nellicus.net/rrdp/notification.xml",
		"https://rrdp.twnic.tw/rrdp/notify.xml",
		"https://chloe.sobornost.net/rpki/news.xml",
		"https://rrdp.krill.nlnetlabs.nl/notification.xml",
		"https://rov-measurements.nlnetlabs.net/rrdp/notification.xml",
		"https://rpki01.hel-fi.rpki.win/rrdp/notification.xml",
		"https://rpki.cnnic.cn/rrdp/notify.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/db9a372a-09bc-4a32-bfe4-8c48e5dbd219/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/f703696e-e47b-4c20-bd93-6f80904e42d2/notification.xml",
		"https://pub.krill.ausra.cloud/rrdp/notification.xml",
		"https://ca.rg.net/rrdp/notify.xml",
		"https://rpkica.mckay.com/rrdp/notify.xml",
		"https://sakuya.nat.moe/rrdp/notification.xml",
		"https://rpki.xindi.eu/rrdp/notification.xml",
		"https://rpki.leitecastro.com/notification.xml",
		"https://rpki.rand.apnic.net/rrdp/notification.xml",
		"https://rpki-rrdp.mnihyc.com/rrdp/notification.xml",
		"https://krill.rg.net/rrdp/notification.xml",
		"https://rpki.luys.cloud/rrdp/notification.xml",
		"https://rpki.apernet.io/rrdp/notification.xml",
		"https://orca.rg.net/rrdp/notification.xml",
		"https://feo.tla.org/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/bd48a1fa-3471-4ab2-8508-ad36b96813e4/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/2f059a21-d41b-4846-b7ae-7ea38c32fd4c/notification.xml",
		"https://subrepo.wildtky.com/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/c3cd7c24-12cb-4abc-8fd2-5e2bcbb85ae6/notification.xml",
		"https://rpki-pp.com/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/708aafaf-00b4-485b-854c-0b32ca30f57b/notification.xml",
		"https://rpki-rrdp.warpnet.xyz/notification.xml",
		"https://rpki.0i1.eu/rrdp/notification.xml",
		"https://rrdp.apnic.net/notification.xml",
		"https://rrdp.ripe.net/notification.xml",
		"https://rrdp.afrinic.net/notification.xml",
		"https://rrdp.lacnic.net/rrdpas0/notification.xml",
		"https://rrdp.lacnic.net/rrdp/notification.xml",
		"https://rrdp.arin.net/notification.xml",
		"https://rrdp-as0.apnic.net/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/602a26e5-4a9e-4e5e-89f0-ef891490d9c9/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/42582c67-dd3f-4bc5-ba60-e97e552c6e35/notification.xml",
		"https://cloudie.rpki.app/rrdp/notification.xml",
		"https://rrdp.paas.rpki.ripe.net/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/517f3ed7-58b5-4796-be37-14d62e48f056/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/b8a1dd25-c313-4f25-ac21-bf55514d9c7d/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/ff9fa84e-9783-4a0b-a58d-6dc8e2433d33/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/dfd7f6d3-e6e9-4987-9ae7-d052c5353898/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/e7518af5-a343-428d-bf78-f982b6e60505/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/71e5236f-c6f1-4928-a1b9-8def09c06085/notification.xml",
		"https://krill.accuristechnologies.ca/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/b3f6b688-cff4-402f-97d5-02f6f1886b7e/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/20aa329b-fc52-4c61-bf53-09725c042942/notification.xml",
		"https://rpki.tools.westconnect.ca/rrdp/notification.xml",
		"https://rrdp-rps.arin.net/notification.xml",
		"https://rpki.roa.net/rrdp/notification.xml",
		"https://repodepot.wildtky.com/rrdp/notification.xml",
		"https://repo.kagl.me/rpki/notification.xml",
		"https://rpki.admin.freerangecloud.com/rrdp/notification.xml",
		"https://oto.wakuwaku.ne.jp/pki/oshirase.xml",
		"https://rpki.sailx.co/rrdp/notification.xml",
		"https://rpki.multacom.com/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/e72d8db0-4728-4fc1-bdd8-471129866362/notification.xml",
		"https://rrdp.sub.apnic.net/notification.xml",
		"https://repo.rpki.space/rrdp/notification.xml",
		"https://rpki-publication.haruue.net/rrdp/notification.xml",
		"https://rpki.zappiehost.com/rrdp/notification.xml",
		"https://rpki.sn-p.io/rrdp/notification.xml",
		"https://rpki.komorebi.network:3030/rrdp/notification.xml",
		"https://rpki.cc/rrdp/notification.xml",
		"https://magellan.ipxo.com/rrdp/notification.xml",
		"https://rpki.qs.nu/rrdp/notification.xml",
		"https://rpki.pudu.be/rrdp/notification.xml",
		"https://rpki.as207960.net/rrdp/notification.xml",
		"https://rpki.uz/rrdp/notification.xml",
		"https://rpki-01.pdxnet.uk/rrdp/notification.xml",
		"https://krill.uta.ng:3030/rrdp/notification.xml",
		"https://rpki-repo.registro.br/rrdp/notification.xml",
		"https://repo-rpki.idnic.net/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/08c2f264-23f9-49fb-9d43-f8b50bec9261/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/54602fb0-a9d4-4f9f-b0ca-be2a139ea92b/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/dba8f01c-9669-44a3-ac6e-db2edb099b84/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/a841823c-a10d-477c-bfdf-4086f0b1594c/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/967a255c-d680-42d3-9ec3-ecb3f9da088c/notification.xml",
		"https://rpki.sunoaki.net/rrdp/notification.xml",
		"https://rrdp.twnic.tw/rrdp/notify.xml",
		"https://rrdp.rp.ki/notification.xml",
		"https://rpki.folf.systems/rrdp/notification.xml",
		"https://rpki.athene-center.net/rrdp/notification.xml",
		"https://ca.nat.moe/rrdp/notification.xml",
		"https://rpki.nellicus.net/rrdp/notification.xml",
		"https://rpki-repository.nic.ad.jp/rrdp/ap/notification.xml",
		"https://0.sb/rrdp/notification.xml",
		"https://rpki.owl.net/rrdp/notification.xml",
		"https://chloe.sobornost.net/rpki/news.xml",
		"https://rrdp.krill.nlnetlabs.nl/notification.xml",
		"https://rov-measurements.nlnetlabs.net/rrdp/notification.xml",
		"https://rpki.cnnic.cn/rrdp/notify.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/db9a372a-09bc-4a32-bfe4-8c48e5dbd219/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/f703696e-e47b-4c20-bd93-6f80904e42d2/notification.xml",
		"https://sakuya.nat.moe/rrdp/notification.xml",
		"https://ca.rg.net/rrdp/notify.xml",
		"https://rpki.xindi.eu/rrdp/notification.xml",
		"https://rpki.leitecastro.com/notification.xml",
		"https://rpki.rand.apnic.net/rrdp/notification.xml",
		"https://rpki-rrdp.mnihyc.com/rrdp/notification.xml",
		"https://krill.rg.net/rrdp/notification.xml",
		"https://rpki.luys.cloud/rrdp/notification.xml",
		"https://rpki.apernet.io/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/bd48a1fa-3471-4ab2-8508-ad36b96813e4/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/708aafaf-00b4-485b-854c-0b32ca30f57b/notification.xml",
		"https://subrepo.wildtky.com/rrdp/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/2f059a21-d41b-4846-b7ae-7ea38c32fd4c/notification.xml",
		"https://rpki-rrdp.warpnet.xyz/notification.xml",
		"https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/c3cd7c24-12cb-4abc-8fd2-5e2bcbb85ae6/notification.xml",
		"https://rrdp.rpki.tianhai.link/rrdp/notification.xml",
		"https://rpki-pp.com/rrdp/notification.xml",
	}
	m := make(map[string]string)
	for _, u := range urls {
		host, err := url.Parse(u)
		if err != nil {
			fmt.Println("err:", u, err)
			continue
		}
		domain := host.Hostname()
		if _, ok := m[domain]; !ok {
			m[domain] = domain
		}
	}
	fmt.Println("总数:", len(m))
	onlyIpv4 := 0
	onlyIpv6 := 0
	both := 0
	none := 0
	for _, dns := range m {
		// 解析ip地址
		fmt.Println(dns)
		ns, err := net.LookupIP(dns)
		if err != nil {
			fmt.Println("lookupIp fail", dns, err)
			continue
		}
		haveIpv4 := false
		haveIpv6 := false
		for _, n := range ns {
			ip := net.ParseIP(n.String())
			if ip.To4() != nil {
				fmt.Println("IPv4:", n.String())
				haveIpv4 = true
			} else if ip.To16() != nil {
				fmt.Println("IPv6:", n.String())
				haveIpv6 = true
			}

		}
		if haveIpv4 && haveIpv6 {
			both++
		} else if haveIpv4 {
			onlyIpv4++
		} else if haveIpv6 {
			onlyIpv6++
		} else {
			none++
		}
		fmt.Println("------------")
	}
	fmt.Println("总数:", len(m))
	fmt.Println("onlyIpv4:", onlyIpv4)
	fmt.Println("onlyIpv6:", onlyIpv6)
	fmt.Println("both:", both)
	fmt.Println("none:", none)
}
