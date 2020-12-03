package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type Impact struct {
	PrefixHijack      bool `json:"prefix-hijack"`
	SubprefixHijacked bool `json:"subprefix-hijacked"`
	SubprefixHijack   bool `json:"subprefix-hijack"`
}
type Relation struct {
	Id       uint64        `json:"id"`
	File     string        `json:"file"`
	Children []interface{} `json:"children"`
}
type PrefixFilter struct {
	Asn             uint64 `json:"asn"`
	Prefix          string `json:"prefix"`
	MaxPrefixLength uint64 `json:"maxPrefixLength"`
	Comment         string `json:"comment"`
}
type PrefixAssertion struct {
	Asn             uint64 `json:"asn"`
	Prefix          string `json:"prefix"`
	MaxPrefixLength uint64 `json:"maxPrefixLength"`
	Comment         string `json:"comment"`
}
type Solution struct {
	PrefixFilters    []PrefixFilter    `json:"prefixFilters"`
	PrefixAssertions []PrefixAssertion `json:"prefixAssertions"`
}

type Competitor struct {
	CompPfx               string         `json:"comp_pfx"`
	CompRoaPfxesAll       []string       `json:"comp_roa_pfxes_all"`
	CompRoaAsn            uint64         `json:"comp_roa_asn"`
	CompRoaValidityPeriod ValidityPeriod `json:"comp_roa_validity_period"`
	InvolvedPfx           string         `json:"involved_pfx"`
	CompRoaAsInfo         string         `json:"comp_roa_as_info"`
	CompRoaName           string         `json:"comp_roa_name"`

	Impact   Impact   `json:"impact"`
	Relation Relation `json:"relation"`
	Solution Solution `json:"solution"`
}

type ValidityPeriod struct {
	Valto   string `json:"valto"`
	Valfrom string `json:"valfrom"`
}
type RoaComp struct {
	AsInfo         string         `json:"as_info"`
	PfxesAll       []string       `json:"pfxes_all"`
	Asn            uint64         `json:"asn"`
	ValidityPeriod ValidityPeriod `json:"validity_period"`
	Competitors    []Competitor   `json:"competitors"`
}

func main() {
	json := `
{
  "as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
  "competitors": [
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8001::/48-48", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8001::/48-48"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:40+08:00", 
        "valfrom": "2020-08-23T03:22:40+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383030313a3a2f34382d3438203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8001::/48", 
            "asn": 264242, 
            "maxPrefixLength": 48
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8002::/48-48", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8002::/48-48"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:44+08:00", 
        "valfrom": "2020-08-23T03:22:44+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383030323a3a2f34382d3438203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8002::/48", 
            "asn": 264242, 
            "maxPrefixLength": 48
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8100::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8100::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:46+08:00", 
        "valfrom": "2020-08-23T03:22:46+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383130303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8100::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8200::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8200::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:49+08:00", 
        "valfrom": "2020-08-23T03:22:49+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383230303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8200::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8300::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8300::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:53+08:00", 
        "valfrom": "2020-08-23T03:22:53+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383330303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8300::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8400::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8400::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:55+08:00", 
        "valfrom": "2020-08-23T03:22:55+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383430303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8400::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8500::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8500::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:58+08:00", 
        "valfrom": "2020-08-23T03:22:58+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383530303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8500::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8600::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8600::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:28:02+08:00", 
        "valfrom": "2020-08-23T03:23:02+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383630303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8600::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8700::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8700::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:28:05+08:00", 
        "valfrom": "2020-08-23T03:23:05+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383730303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8700::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8800::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8800::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:28:08+08:00", 
        "valfrom": "2020-08-23T03:23:08+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383830303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8800::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:9000::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:9000::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:28:11+08:00", 
        "valfrom": "2020-08-23T03:23:11+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a393030303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:9000::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:9100::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:9100::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:28:14+08:00", 
        "valfrom": "2020-08-23T03:23:14+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a393130303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:9100::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:9200::/40-40", 
      "comp_roa_pfxes_all": [
        "2804:24cc:9200::/40-40"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:28:17+08:00", 
        "valfrom": "2020-08-23T03:23:17+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a393230303a3a2f34302d3430203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:9200::/40", 
            "asn": 264242, 
            "maxPrefixLength": 40
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc::/48-48", 
      "comp_roa_pfxes_all": [
        "2804:24cc:0000::/48-48"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:26+08:00", 
        "valfrom": "2020-08-23T03:22:26+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a3a2f34382d3438203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc::/48", 
            "asn": 264242, 
            "maxPrefixLength": 48
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:1::/48-48", 
      "comp_roa_pfxes_all": [
        "2804:24cc:0001::/48-48"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:33+08:00", 
        "valfrom": "2020-08-23T03:22:33+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a313a3a2f34382d3438203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:1::/48", 
            "asn": 264242, 
            "maxPrefixLength": 48
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }, 
    {
      "impact": {
        "prefix-hijack": false, 
        "subprefix-hijacked": true, 
        "subprefix-hijack": false
      }, 
      "comp_pfx": "2804:24cc:8000::/48-48", 
      "comp_roa_pfxes_all": [
        "2804:24cc:8000::/48-48"
      ], 
      "relation": null, 
      "comp_roa_asn": 264242, 
      "comp_roa_validity_period": {
        "valto": "2021-08-23T03:27:37+08:00", 
        "valfrom": "2020-08-23T03:22:37+08:00"
      }, 
      "involved_pfx": "2804:24cc::/32-32", 
      "comp_roa_as_info": " Diogo Cassio Cabral Me, BR; lacnic ", 
      "comp_roa_name": "323830343a323463633a383030303a3a2f34382d3438203d3e20323634323432.roa", 
      "solution": {
        "prefixFilters": [
          {
            "comment": "ROA Competition Monitor", 
            "prefix": "2804:24cc:8000::/48", 
            "asn": 264242, 
            "maxPrefixLength": 48
          }
        ], 
        "prefixAssertions": [
          null
        ]
      }
    }
  ], 
  "pfxes_all": [
    "2804:24cc::/32-32"
  ], 
  "asn": 264242, 
  "validity_period": {
    "valto": "2021-08-23T03:26:59+08:00", 
    "valfrom": "2020-08-23T03:21:59+08:00"
  }
}
`
	roaComp := RoaComp{}
	err := jsonutil.UnmarshalJson(json, &roaComp)
	fmt.Println(roaComp, err)
	fmt.Println(jsonutil.MarshalJson(roaComp))
}
