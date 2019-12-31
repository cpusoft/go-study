package main

import (
	"fmt"
	asn1go "github.com/chemikadze/asn1go"
)

func TestDefinitiveIdentifier() {
	content := `
	KerberosV5Spec2 {
        iso(1) identified-organization(3) dod(6)
        nameform
        42 --numberform
        mixedform(88)
	} DEFINITIONS EXPLICIT TAGS ::= BEGIN
	END
	`
	r, err := asn1go.ParseString(content)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("%v", r)
}

func main() {
	TestDefinitiveIdentifier()
}
