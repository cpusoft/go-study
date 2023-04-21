package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

func main() {
	content := `
	{
		"slurmVersion": 1,
		"validationOutputFilters": {
			"prefixFilters": [
				{
					"asn": 398109,
					"prefix": "173.139.139/24",
					"comment": "test prefix 9"
				},
				{
					"asn": 398100,
					"prefix": "173.139.130/24",
					"comment": "test prefix 0"
				}
			],
			"aspaFilters": [
				{
					"customerAsn": 64495,
					"comment": "Ignore ASPA(s) that have 64496 as Customer ASID"
				},
				{
					"customerASID": 64696,
					"providerAsns": [
						{
							"providerAsn": 64697
						},
						{
							"addressFamily": 1,
							"providerAsn": 64698
						},
						{
							"addressFamily": 2,
							"providerAsn": 64699
						}
					],
					"comment": "Ignore ASPA(s) that have 64696 as Customer ASID, and have 64697 or 64698 or 64699 as Provider ASID"
				}
			]
		},
		"locallyAddedAssertions": {
			"prefixAssertions": [
				{
					"asn": 64406,
					"prefix": "198.51.100.6/24",
					"comment": "test assertion ipv4"
				},
				{
					"asn": 64406,
					"prefix": "2001:DB6::/32",
					"maxPrefixLength": 48,
					"comment": "test assertion ipv6"
				}
			],
			"aspaAssertions": [
				{
					"customerAsn": 64416,
					"providerAsns": [
						{
							"providerAsn": 64617
						},
						{
							"addressFamily": 1,
							"providerAsn": 64618
						},
						{
							"addressFamily": 2,
							"providerAsn": 64619
						}
					],
					"comment": "Pretend 64617,64618 and 64619 are upstream for 64416"
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
