package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {

	//https://blog.csdn.net/Zhymax/article/details/7683925
	// 比较底层的，能识别ROA各个键值
	//openssl asn1parse -in privatekey.der -inform DER

	cmd := exec.Command("openssl", "asn1parse", "-in", "G:\\Download\\cert\\029d506f4cfb1a1e4eae7d68b5ebbc15c8b52c93.crl", "--inform", "der")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	result := string(output)
	results := strings.Split(result, "\r\n")
	fmt.Println(len(results))

	for i, one := range results {
		if strings.Contains(one, ":X509v3 Authority Key Identifier") {
			akis := results[i+1]
			index := strings.Index(akis, "[HEX DUMP]:")
			aki := string([]byte(akis)[index+len("[HEX DUMP]:"):])
			fmt.Println(aki)
		}
	}
}
