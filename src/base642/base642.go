package main

import (
	"fmt"

	"github.com/cpusoft/goutil/base64util"
	"github.com/cpusoft/goutil/convert"
)

func main() {
	s := `MIGcMBShEjAQMA4EAQIwCQMHACABBnwgjDALBglghkgBZQMEAgEwdzA0FhBiNDJfaXB2Nl9sb2EucG5nBCCVFt1kvnwXJbn8oRcSDljo2EKlIGhzOZs93/yRxLas8DA/FhtiNDJfc2VydmljZV9kZWZpbml0aW9uLmpzb24EIArhOUciAFzZL0xqoCTV1rPi5n1inxFyDZR4pjOhF6HH`

	b, err := base64util.DecodeBase64(s)
	s = convert.PrintBytes(b, 8)
	fmt.Println(b, err)
	fmt.Println(s)

	s = `RHVwbGljYXRlIGFubm91bmNlbWVudCByZWNlaXZlZA==`
	b, err = base64util.DecodeBase64(s)
	s = string(b)
	fmt.Println(b, err)
	fmt.Println(s)

}
