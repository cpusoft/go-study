package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/miekg/dns"
)

// https://help.aliyun.com/zh/dns/dns-over-https?spm=a2c4g.11186623.0.i0
// https://alidns.com/knowledge?type=SETTING_DOCS#company_json
func main() {
	query := dns.Msg{}
	query.SetQuestion("www.taobao.com.", dns.TypeA)
	msg, _ := query.Pack()
	b64 := base64.RawURLEncoding.EncodeToString(msg)
	resp, err := http.Get("https://9**9.alidns.com/dns-query?dns=" + b64)
	if err != nil {
		fmt.Printf("Send query error, err:%v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	response := dns.Msg{}
	response.Unpack(bodyBytes)
	fmt.Printf("Dns answer is :%v\n", response.String())
}
