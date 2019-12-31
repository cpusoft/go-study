package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func main() {

	//https://blog.csdn.net/Zhymax/article/details/7683925
	// 比较底层的，能识别ROA各个键值
	//openssl asn1parse -in privatekey.der -inform DER
	//no_subject,no_header,no_version,no_serial,no_signame,no_validity,no_subject,no_issuer,no_pubkey,no_sigdump,no_aux

	//	   	ca_default
	//	          the value used by the ca utility, equivalent to no_issuer, no_pubkey,
	//	          no_header, no_version, no_sigdump and no_signame.
	//no_subject,no_header,no_version,no_serial,no_signame,no_validity,no_subject,no_issuer,no_pubkey,no_sigdump,no_aux
	/*  static const NAME_EX_TBL cert_tbl[] = {
	    {"compatible", X509_FLAG_COMPAT, 0xffffffffl},
	    {"ca_default", X509_FLAG_CA, 0xffffffffl},
	    {"no_header", X509_FLAG_NO_HEADER, 0},
	    {"no_version", X509_FLAG_NO_VERSION, 0},
	    {"no_serial", X509_FLAG_NO_SERIAL, 0},
	    {"no_signame", X509_FLAG_NO_SIGNAME, 0},
	    {"no_validity", X509_FLAG_NO_VALIDITY, 0},
	    {"no_subject", X509_FLAG_NO_SUBJECT, 0},
	    {"no_issuer", X509_FLAG_NO_ISSUER, 0},
	    {"no_pubkey", X509_FLAG_NO_PUBKEY, 0},
	    {"no_extensions", X509_FLAG_NO_EXTENSIONS, 0},
	    {"no_sigdump", X509_FLAG_NO_SIGDUMP, 0},
	    {"no_aux", X509_FLAG_NO_AUX, 0},
	    {"no_attributes", X509_FLAG_NO_ATTRIBUTES, 0},
	    {"ext_default", X509V3_EXT_DEFAULT, X509V3_EXT_UNKNOWN_MASK},
	    {"ext_error", X509V3_EXT_ERROR_UNKNOWN, X509V3_EXT_UNKNOWN_MASK},
	    {"ext_parse", X509V3_EXT_PARSE_UNKNOWN, X509V3_EXT_UNKNOWN_MASK},
	    {"ext_dump", X509V3_EXT_DUMP_UNKNOWN, X509V3_EXT_UNKNOWN_MASK},
	*/
	// 比较上层，只能cer、crl
	//openssl x509 -noout -text -in 2.cer --inform der

	cmd := exec.Command("openssl", "x509", "-noout", "-text", "-in", "G:\\Download\\cert\\H.cer", "--inform", "der")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	//fmt.Println("Result: " + out.String())

	cmd = exec.Command("openssl", "x509", "-noout", "-text", "-in", "G:\\Download\\cert\\H.cer", "--inform", "der")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	result := string(output)
	results := strings.Split(result, "\r\n")
	fmt.Println(len(results))

	for i, one := range results {
		if strings.Contains(one, "Version:") {
			version := strings.TrimSpace(strings.Split(one, ":")[1])
			fmt.Println("Version:", version)
		} else if strings.Contains(one, "Serial Number:") {
			SerialNumber := strings.TrimSpace(results[i+1])
			fmt.Println("SerialNumber:", SerialNumber)
		} else if strings.Contains(one, "Issuer:") {
			Issuer := strings.TrimSpace(strings.Split(one, ":")[1])
			fmt.Println("Issuer:", Issuer)
		} else if strings.Contains(one, "Issuer:") {
			Issuer := strings.TrimSpace(strings.Split(one, ":")[1])
			fmt.Println("Issuer:", Issuer)
		} else if strings.Contains(one, "Not Before:") {
			NotBefore := strings.TrimSpace(strings.Replace(one, "Not Before:", "", -1))
			fmt.Println("NotBefore:", NotBefore)
		} else if strings.Contains(one, "Not After :") {
			NotAfter := strings.TrimSpace(strings.Replace(one, "Not After :", "", -1))
			fmt.Println("NotAfter:", NotAfter)
		} else if strings.Contains(one, "sbgp-ipAddrBlock:") {
			ips := results[i+1:]
			end := 0
			for j, ip := range ips {
				if strings.Contains(ip, "IPv4:") || strings.Contains(ip, "IPv6:") {
					//ignore
					continue
				}
				if len(strings.TrimSpace(ip)) == 0 {
					// end
					end = j
					fmt.Println(end)
					break
				}
			}
			ips = ips[:end]
			fmt.Println(ips)
		} else if strings.Contains(one, "Autonomous System Numbers:") {
			ass := results[i+1:]
			end := 0
			for j, as := range ass {
				if strings.Contains(as, "IPv4:") || strings.Contains(as, "IPv6:") {
					//ignore
					continue
				}
				if len(strings.TrimSpace(as)) == 0 {
					// end
					end = j
					fmt.Println(end)
					break
				}
			}
			ass = ass[:end]
			fmt.Println(ass)
		}
	}
}
