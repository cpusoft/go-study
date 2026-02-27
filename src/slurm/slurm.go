package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	config "github.com/beego/beego/v2/core/config"
)

// slurm head , 这里特殊处理，相当于每组只会有一个
type SlurmTarget struct {
	Asn      int64  `json:"asn,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// filter
type PrefixFilters struct {
	Asn     int64  `json:"asn"`
	Prefix  string `json:"prefix"`
	Comment string `json:"comment"`
}

type BgpsecFilters struct {
	Asn       int64  `json:"asn"`
	RouterSKI string `json:"routerSKI"`
	Comment   string `json:"comment"`
}

type ValidationOutputFilters struct {
	PrefixFilters []PrefixFilters `json:"prefixFilters"`
	BgpsecFilters []BgpsecFilters `json:"bgpsecFilters"`
}

// assertion
type PrefixAssertions struct {
	Asn             int64  `json:"asn"`
	Prefix          string `json:"prefix"`
	MaxPrefixLength int    `json:"maxPrefixLength"`
	Comment         string `json:"comment"`
}

type BgpsecAssertions struct {
	Asn       int64  `json:"asn"`
	Comment   string `json:"comment"`
	SKI       string `json:"SKI"`
	PublicKey string `json:"publicKey"`
}

type LocallyAddedAssertions struct {
	PrefixAssertions []PrefixAssertions `json:"prefixAssertions"`
	BgpsecAssertions []BgpsecAssertions `json:"bgpsecAssertions"`
}

type Slurm struct {
	SlurmVersion            int                     `json:"slurmVersion"`
	SlurmTarget             []SlurmTarget           `json:"slurmTarget"`
	ValidationOutputFilters ValidationOutputFilters `json:"validationOutputFilters"`
	LocallyAddedAssertions  LocallyAddedAssertions  `json:"locallyAddedAssertions"`
}

func main() {

	// load config
	conf, err := config.NewConfig("ini", "E:\\Go\\test1\\src\\main\\slurm.conf")
	if err != nil {
		fmt.Println("load config failed, err:", err)
		return
	}
	t, _ := conf.Strings("target::asn")
	asnInConfs := strings.Split(t[0], ",")
	t, _ = conf.Strings("target::hostname")
	hostnameInConfs := strings.Split(t[0], ",")
	fmt.Println("asnInConfs:", asnInConfs)
	fmt.Println("hostnameInConfs:", hostnameInConfs)

	// slurm init
	slurm := Slurm{}

	errMsg := ""

	// load slurm json file
	f, err := ioutil.ReadFile("E:\\Go\\test1\\src\\main\\slurm.json")
	if err != nil {
		fmt.Println("Load config fail: ", err)
		errMsg += (err.Error() + "; ")
	}

	err = json.Unmarshal(f, &slurm)
	if err != nil {
		fmt.Println("Para json failed: ", err)
		errMsg += (err.Error() + "; ")
	}
	fmt.Print(slurm)

	// check slurm file is right ?
	if slurm.SlurmVersion != 1 {
		fmt.Println("slurm version is not 1 ")
		errMsg += ("slurm version is not 1; ")
	}

	// check target, and get valid asn and hostname
	// if no target is ok, if target is not in slurm.conf ,is error;
	var targetValid bool = false
	var asnValid []int64
	var hostnameValid []string
	if slurm.SlurmTarget != nil {
		for _, slurmTarget := range slurm.SlurmTarget {
			fmt.Println("slurmTarget: ", slurmTarget)
			if slurmTarget.Asn > 0 {
				for _, asnInConf := range asnInConfs {
					fmt.Println("asnInConf: ", asnInConf, strconv.FormatInt(slurmTarget.Asn, 10))
					if strings.Compare(asnInConf, strconv.FormatInt(slurmTarget.Asn, 10)) == 0 {
						asnValid = append(asnValid, slurmTarget.Asn)
					}
				}
			}
			fmt.Println("asnValid:", asnValid)

			if slurmTarget.Hostname != "" {
				for _, hostnameInConf := range hostnameInConfs {
					fmt.Println("hostnameInConf: ", hostnameInConf)
					if strings.Compare(hostnameInConf, slurmTarget.Hostname) == 0 {
						hostnameValid = append(hostnameValid, slurmTarget.Hostname)
					}
				}
			}
			fmt.Println("hostnameValid:", hostnameValid)
		}
		if len(asnValid) > 0 || len(hostnameValid) > 0 {
			targetValid = true
		}
	} else {
		targetValid = true
	}
	fmt.Println("targetValid:", targetValid)
	fmt.Println("asnValid:", asnValid)
	fmt.Println("hostnameValid:", hostnameValid)
	if !targetValid {
		errMsg += ("target is not valid; ")
	}

	//check filter list
	if len(slurm.ValidationOutputFilters.PrefixFilters) > 0 ||
		len(slurm.ValidationOutputFilters.BgpsecFilters) > 0 {
		var prefixFilters = slurm.ValidationOutputFilters.PrefixFilters
		var bgpsecFilters = slurm.ValidationOutputFilters.BgpsecFilters

		if len(prefixFilters) > 0 {
			for _, prefixFilter := range prefixFilters {
				fmt.Println("prefixFilter:", prefixFilter)
				if len(prefixFilter.Prefix) > 0 {
					if err := checkPrefix(prefixFilter.Prefix); err != nil {
						errMsg += (err.Error() + "; ")
					}
				}

			}
		}

		if len(bgpsecFilters) > 0 {

		}

	}
	if len(errMsg) > 0 {
		fmt.Println(errMsg)
	}
}

func checkPrefix(prefix string) error {
	if !strings.Contains(prefix, "/") {
		return errors.New("prefix is not contains '/' ")
	}
	return nil
}
func checkAsn(asn int) {
	if asn == 0 {

	}
}
