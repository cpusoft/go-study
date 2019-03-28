package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "os"
	"strconv"
	"strings"
)

type Mft struct {
	Roas    []string `json:"roas"`
	Valto   string   `json:"valto"`
	Name    string   `json:"name"`
	Valfrom string   `json:"valfrom"`
	CaCerts []string `json:"cacerts"`
}

type EE struct {
	Prefix  []string `json:"prefix"`
	Valto   string   `json:"valto"`
	Name    string   `json:"name"`
	Valfrom string   `json:"valfrom"`
}

type CRL struct {
	Valto   string `json:"valto"`
	Name    string `json:"name"`
	Valfrom string `json:"valfrom"`
}

type Roa struct {
	Prefix  []string `json:"prefix"`
	Warning []string `json:"warning"`
	Mft     Mft      `json:"mft"`
	Name    string   `json:"name"`
	Error   []string `json:"error"`
	EE      EE       `json:"ee"`
	ASN     int      `json:"asn"`
	CRL     CRL      `json:"crl"`
}

type Chain struct {
	Cert        string   `json:"cert"`
	Valto       string   `json:"valto"`
	Mft         Mft      `json:"mft"`
	Valfrom     string   `json:"valfrom"`
	Error       []string `json:"error"`
	Prefix      []string `json:"prefix"`
	Warning     []string `json:"warning"`
	ASN         []int    `json:"asn"`
	CRL         CRL      `json:"crl"`
	PrefixOther []string `json:"prefix_other"`
}
type TA struct {
	Prefix      []string `json:"prefix"`
	Warning     []string `json:"warning"`
	Name        string   `json:"name"`
	Valfrom     string   `json:"valfrom"`
	Error       []string `json:"error"`
	Valto       string   `json:"valto"`
	ASN         []int    `json:"asn"`
	PrefixOther []string `json:"prefix_other"`
}
type StatResult struct {
	Roa   Roa     `json:"roa"`
	Chain []Chain `json:"chain"`
	TA    TA      `json:"ta"`
}

/////////////////////////////
//////   new data struct
///////////////////////////

type TAL_new struct {
	IpPrefix string `json:"IPPrefix"`
	Name     string `json:"name"`
	Asns     string `json:"ASNs"`
}

type EE_new struct {
	Name     string `json:"name"`
	IpPrefix string `json:"IPPrefix"`
}
type ROA_new struct {
	Name     string   `json:"name"`
	EE_new   EE_new   `json:"EE"`
	Asns     string   `json:"ASNs"`
	IpPrefix string   `json:"IPPrefix"`
	Warn     string   `json:"warn"`
	Error    []string `json:"error"`
}

type ISP_new struct {
	IpPrefix string   `json:"IPPrefix"`
	MFTWarn  []string `json:"MFTWarn"`
	Name     string   `json:"name"`
	Asns     string   `json:"ASNs"`
}

type StatResult_new1 struct {
	TAL_new  TAL_new `json:"TAL"`
	ROA_new  ROA_new `json:"ROA"`
	ISP_new1 ISP_new `json:"ISP1"`
}
type StatResult_new2 struct {
	TAL_new  TAL_new `json:"TAL"`
	ROA_new  ROA_new `json:"ROA"`
	ISP_new1 ISP_new `json:"ISP1"`
	ISP_new2 ISP_new `json:"ISP2"`
}

type StatResult_new3 struct {
	TAL_new  TAL_new `json:"TAL"`
	ROA_new  ROA_new `json:"ROA"`
	ISP_new1 ISP_new `json:"ISP1"`
	ISP_new2 ISP_new `json:"ISP2"`
	ISP_new3 ISP_new `json:"ISP3"`
}

type StatResult_new4 struct {
	TAL_new  TAL_new `json:"TAL"`
	ROA_new  ROA_new `json:"ROA"`
	ISP_new1 ISP_new `json:"ISP1"`
	ISP_new2 ISP_new `json:"ISP2"`
	ISP_new3 ISP_new `json:"ISP3"`
	ISP_new4 ISP_new `json:"ISP4"`
}

type StatResult_new5 struct {
	TAL_new  TAL_new `json:"TAL"`
	ROA_new  ROA_new `json:"ROA"`
	ISP_new1 ISP_new `json:"ISP1"`
	ISP_new2 ISP_new `json:"ISP2"`
	ISP_new3 ISP_new `json:"ISP3"`
	ISP_new4 ISP_new `json:"ISP4"`
	ISP_new5 ISP_new `json:"ISP5"`
}

type StatResult_new6 struct {
	TAL_new  TAL_new `json:"TAL"`
	ROA_new  ROA_new `json:"ROA"`
	ISP_new1 ISP_new `json:"ISP1"`
	ISP_new2 ISP_new `json:"ISP2"`
	ISP_new3 ISP_new `json:"ISP3"`
	ISP_new4 ISP_new `json:"ISP4"`
	ISP_new5 ISP_new `json:"ISP5"`
}

func main() {
	path := `G:\Download\cert\data_20190327\`
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("num:", len(fileInfos))
	for _, fileInfo := range fileInfos {
		fmt.Println("readFile :", fileInfo.Name())
		contents, err := ioutil.ReadFile(path + fileInfo.Name())
		if err != nil {
			fmt.Println("readFile err:", fileInfo.Name(), err)
			continue
		}
		statResult := StatResult{}
		err = json.Unmarshal(contents, &statResult)
		if err != nil {
			fmt.Println("Unmarshal json failed: ", fileInfo.Name(), err)
			continue
		}

		////////////// convert //////////////

		/////// TAL ////////////////
		ipPrefix := "IP Prefix:" + strings.Join(statResult.TA.Prefix, ",")
		asns := IntJoin(statResult.TA.ASN)
		tal_new := TAL_new{
			IpPrefix: ipPrefix,
			Name:     statResult.TA.Name,
			Asns:     asns,
		}

		//////// ROA ////////////
		ee_new := EE_new{
			Name:     statResult.Roa.EE.Name,
			IpPrefix: strings.Join(statResult.Roa.EE.Prefix, ","),
		}
		name := "ROA" //statResult.Roa.Name,
		asns = "ASNs: " + strconv.Itoa(statResult.Roa.ASN)
		ipPrefix = "IP Prefix:" + strings.Join(statResult.Roa.Prefix, ",")
		roa_new := ROA_new{
			Name:     name,
			EE_new:   ee_new,
			Asns:     asns,
			IpPrefix: ipPrefix,
			Warn:     strings.Join(statResult.Roa.Warning, ","),
			Error:    statResult.Roa.Error,
		}

		//////// ISP //////////////////
		isp_news := make([]ISP_new, 0)
		for _, chain := range statResult.Chain {
			ipPrefix = "IP Prefix:" + strings.Join(chain.Prefix, ",")
			asns = "ASNs: " + IntJoin(chain.ASN)
			isp_new := ISP_new{
				IpPrefix: ipPrefix,
				MFTWarn:  chain.Warning,
				Name:     chain.Cert,
				Asns:     asns,
			}
			isp_news = append(isp_news, isp_new)
		}

		/////// Restul /////////////

		var datanew []byte

		//注意 chain和ispnew 是反着的
		if len(isp_news) == 1 {
			statResult_new := StatResult_new1{
				TAL_new:  tal_new,
				ROA_new:  roa_new,
				ISP_new1: isp_news[0],
			}
			datanew, _ = json.Marshal(statResult_new)

		} else if len(isp_news) == 2 {
			statResult_new := StatResult_new2{
				TAL_new:  tal_new,
				ROA_new:  roa_new,
				ISP_new1: isp_news[1],
				ISP_new2: isp_news[0],
			}
			datanew, _ = json.Marshal(statResult_new)

		} else if len(isp_news) == 3 {
			statResult_new := StatResult_new3{
				TAL_new:  tal_new,
				ROA_new:  roa_new,
				ISP_new1: isp_news[2],
				ISP_new2: isp_news[1],
				ISP_new3: isp_news[0],
			}
			datanew, _ = json.Marshal(statResult_new)

		} else if len(isp_news) == 4 {
			statResult_new := StatResult_new4{
				TAL_new:  tal_new,
				ROA_new:  roa_new,
				ISP_new1: isp_news[3],
				ISP_new2: isp_news[2],
				ISP_new3: isp_news[1],
				ISP_new4: isp_news[0],
			}
			datanew, _ = json.Marshal(statResult_new)

		} else if len(isp_news) == 5 {
			statResult_new := StatResult_new5{
				TAL_new:  tal_new,
				ROA_new:  roa_new,
				ISP_new1: isp_news[4],
				ISP_new2: isp_news[3],
				ISP_new3: isp_news[2],
				ISP_new4: isp_news[1],
				ISP_new5: isp_news[0],
			}
			datanew, _ = json.Marshal(statResult_new)

		}
		newName := strings.Replace(fileInfo.Name(), ".0", "", -1)
		path_new := `G:\Download\cert\datanew\`
		ioutil.WriteFile(path_new+newName, datanew, 0644)
	}

	fmt.Println("ok")
}

func IntJoin(asn []int) string {
	l := len(asn)
	if l == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, a := range asn {
		buffer.WriteString(strconv.Itoa(a) + ",")
	}

	return buffer.String()
}
