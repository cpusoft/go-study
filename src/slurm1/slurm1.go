package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

func main() {
	content := `
	{
		"slurmVersion": 2,
		"validationOutputFilters": {
			"prefixFilters": [
				{
					"asn": 398109,
					"prefix": "173.139.139/24",
					"maxPrefixLength": 0,
					"comment": "test prefix 9"
				},
				{
					"asn": 398100,
					"prefix": "173.139.130/24",
					"maxPrefixLength": 0,
					"comment": "test prefix 0"
				}
			],
			"bgpsecFilters": null,
			"aspaFilters": [
				{
					"customerAsid": 64496,
					"comment": "Filter out all VAPs that have 64496 as Customer ASID"
				},
				{
					"customerAsid": 64497,
					"providers": [
						{
							"providerAsid": 64498
						},
						{
							"providerAsid": 64499,
							"afiLimit": "IPv4"
						},
						{
							"providerAsid": 64500,
							"afiLimit": "IPv6"
						}
					],
					"comment": "Filter some providers with 64497 as Customer ASID"
				},
				{
					"providers": [
						{
							"providerAsid": 65001
						}
					],
					"comment": "Never accept 65001 as a valid provider."
				}
			]
		},
		"locallyAddedAssertions": {
			"prefixAssertions": [
				{
					"asn": 64406,
					"prefix": "198.51.100.6/24",
					"maxPrefixLength": 0,
					"comment": "test assertionipv4",
					"treatLevel": ""
				},
				{
					"asn": 64406,
					"prefix": "2001:DB6::/32",
					"maxPrefixLength": 48,
					"comment": "test assertion ipv6",
					"treatLevel": ""
				}
			],
			"bgpsecAssertions": null,
			"aspaAssertions": [
				{
					"customerAsid": 64496,
					"providers": [
						{
							"providerAsid": 64498
						},
						{
							"providerAsid": 64499,
							"afiLimit": "IPv4"
						},
						{
							"providerAsid": 64500,
							"afiLimit": "IPv6"
						}
					],
					"comment": "Authorize additional providers for customer AS 64496"
				}
			]
		}
	}
	
		
	`
	fmt.Println("main():content:", content)
	slurm := Slurm{}
	err := jsonutil.UnmarshalJson(content, &slurm)
	fmt.Println("main(): slurm:", jsonutil.MarshalJson(slurm), err)
}
