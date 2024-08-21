package main

import (
	"fmt"
	"time"

	"github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/jsonutil"
	dns "labscm.zdns.cn/dns-mod/dns-library"
)

func main() {
	secret := map[string]string{"test.": "pRZgBrBvI4NAHZYhxmhs/Q=="}

	// F:\share\我的坚果云\Go\dns\research\tdns\tdns\childsync_utils.go
	m := new(dns.Msg)
	m.SetUpdate(dns.Fqdn("example.com"))
	insertRR, err := dns.NewRR("test1.example.com. 300 A 192.0.2.1")
	m.Insert([]dns.RR{insertRR})
	m.SetTsig("test.", dns.HmacSHA256, 300, time.Now().Unix())
	//	var adds, removes []dns.RR
	//	m.Remove(removes)
	//	m.Insert(adds)

	// dns-library
	c := new(dns.Client)
	c.TsigSecret = secret
	belogs.Debug("TestServerRoundtripTsig(): client tsig m:", jsonutil.MarshalJson(m))
	_, _, err = c.Exchange(m, ":1053")
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(5 * time.Minute)
}
