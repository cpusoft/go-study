package main

import "fmt"

var Directives = []string{
	"root",
	"metadata",
	"geoip",
	"cancel",
	"tls",
	"timeouts",
	"reload",
	"nsid",
	"bufsize",
	"bind",
	"debug",
	"trace",
	"ready",
	"health",
	"pprof",
	"prometheus",
	"errors",
	"log",
	"dnstap",
	"local",
	"dns64",
	"acl",
	"any",
	"chaos",
	"loadbalance",
	"tsig",
	"cache",
	"rewrite",
	"header",
	"dnssec",
	"autopath",
	"minimal",
	"template",
	"transfer",
	"hosts",
	"route53",
	"azure",
	"clouddns",
	"k8s_external",
	"kubernetes",
	"file",
	"auto",
	"secondary",
	"etcd",
	"loop",
	"forward",
	"grpc",
	"erratic",
	"whoami",
	"on",
	"sign",
	"view",
}

type ServerType struct {
	// Function that returns the list of directives, in
	// execution order, that are valid for this server
	// type. Directives should be one word if possible
	// and lower-cased.
	Directives  func() []string
	Directives2 []string
}

func main() {
	st := ServerType{
		Directives: func() []string { return Directives },
	}
	st2 := ServerType{
		Directives2: Directives,
	}

	fmt.Println("st:", st)
	dd := st.Directives()
	fmt.Println("dd:", dd)
	fmt.Println("st2:", st2)
	for _, d2 := range st2.Directives2 {
		fmt.Println(d2)
	}
}
