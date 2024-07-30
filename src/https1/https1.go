package main

import (
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// https://github.com/golang/go/issues/21900
func verifyPeerCerts(serverName string, rawCerts [][]byte, verifiedChains [][]*x509.Certificate) (err error) {
	fmt.Println("verifyPeerCerts():serverName,", serverName)

	// some dummy code to check all certs available (not very useful and indeed a security issue if
	// InsecureSkipVerify is set to true and the server supplies arbitrary certs)
	for i := 0; i < len(rawCerts); i++ {
		fmt.Println("verifyPeerCerts():i,", i)
		cert, err := x509.ParseCertificate(rawCerts[i])

		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		hash := sha1.Sum(rawCerts[i])
		fmt.Printf("Fingerprint: %x\n\n", hash)

		fmt.Println("DNSNames:", cert.DNSNames, "Subject:", cert.Subject)
		err = cert.VerifyHostname(serverName)
		if err == nil {
			fmt.Println("VerifyHostname(): VerifyHostname  ok,", serverName)
			//return nil
			return errors.New("shaodebug fail")
		} else {
			fmt.Println("VerifyHostname(): serverName, fail,", serverName, err)
		}
	}
	return err
}

func tlsHost(targetAddr string) string {
	fmt.Println("tlsHost():targetAddr,", targetAddr)
	if strings.LastIndex(targetAddr, ":") > strings.LastIndex(targetAddr, "]") {
		targetAddr = targetAddr[:strings.LastIndex(targetAddr, ":")]
	}
	fmt.Println("tlsHost():then targetAddr,", targetAddr)
	return targetAddr
}

func main() {
	url := "https://tal.apnic.net/apnic.tal"

	tr := &http.Transport{
		DialTLS: func(network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}

			conn, err := tls.DialWithDialer(dialer, network, addr, &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         addr,
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					return verifyPeerCerts(tlsHost(addr), rawCerts, verifiedChains)
				},
			})
			if err != nil {
				return conn, err
			}
			return conn, nil
		},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(url)

	if err != nil {
		fmt.Println("client.Get(), fail,", url, err)
		return
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	fmt.Println(string(body), err)
}
