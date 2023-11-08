package main

import (
	"fmt"
	"strings"
)

const (
	LICENSE_TIME_FORMAT = "2006-01-02T15:04:05Z07" // RFC3339
)

func main() {
	//s := "2006-01-02T15:04:05+07:00"
	//t, err := time.Parse(time.RFC3339, s)
	//fmt.Println(s, " rfc:", time.RFC3339)
	//fmt.Println(t, "fail:", err)

	p := "-----BEGIN PUBLIC KEY-----MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwh2VooRzDT5pT58O5BHbn1PxPnXzGRlgft5CsiC9vTMAmaaMiBXg8di7eUa07dEY9MWqlhzejKWDTRY+Acda1i50wzBKh/eTUzCM/bosdePs2x6sUkNoOpRsQqD+4DknFrDwt2VlMjeYyEAKxZirVB/y3kMX8g+Amj0veQ5Pm5rJWppHjtQJxfox+7FIlTJ8BMyFi6XJgXxXU5oc+zWZuZEl1zldfi/PBU7otnMMw+GoaKTPLTeQwRtnWxFCSRdWNUUnf5wXlAcN5Qx8MVklQI3DJrVjQk79kWhg3U7uxAX85273VlQcHCkRxp5ZMDkpFcfb7XqryspV7eMQfCt2vwIDAQAB-----END PUBLIC KEY-----"
	n := strings.ContainsAny(p, "\n")
	fmt.Println(n)

	p = strings.Replace(p, "-----BEGIN PUBLIC KEY-----", "-----BEGIN PUBLIC KEY-----\n", -1)
	p = strings.Replace(p, "-----END PUBLIC KEY-----", "\n-----END PUBLIC KEY-----", -1)

	n = strings.ContainsAny(p, "\n")
	fmt.Println(n)

}
