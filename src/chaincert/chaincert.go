package main

import (
	"fmt"
	"github.com/cpusoft/goutil/jsonutil"
)

type ChainCer struct {
	Id uint64 `json:"id" xorm:"id int"`
	//all parent cer, trace back to root
	ParentChainCers []ChainCer `json:"parentChainCers"`
}
type ChainCertCer struct {
	Id uint64 `json:"id" xorm:"id int"`

	ParentChainCers []ChainCertCer `json:"parentChainCers,omitempty"`

	ChildChainCers []ChainCertCer `json:"childChainCers,omitempty"`
}

func main() {
	chainCer := ChainCer{}
	chainCer.ParentChainCers = make([]ChainCer, 0)
	chainCer1 := ChainCer{Id: 1}
	chainCer2 := ChainCer{Id: 2}
	chainCer.ParentChainCers = append(chainCer.ParentChainCers, chainCer1)
	chainCer.ParentChainCers = append(chainCer.ParentChainCers, chainCer2)

	chainCertCer := ChainCertCer{}
	chainCertCer.ParentChainCers = make([]ChainCertCer, 0)
	for i, _ := range chainCer.ParentChainCers {
		chainCertCer := ChainCertCer{Id: chainCer.ParentChainCers[i].Id}
		fmt.Println("NewChainCertCer():i chainCertCer:", i, jsonutil.MarshalJson(chainCertCer))
		chainCertCer.ParentChainCers = append(chainCertCer.ParentChainCers, chainCertCer)
		fmt.Println("NewChainCertCer():i, chainCertCer.ParentChainCers:",
			i, jsonutil.MarshalJson(chainCertCer.ParentChainCers), len(chainCertCer.ParentChainCers))
	}
	fmt.Println("NewChainCertCer():chainCertCer.Id:", chainCertCer.Id,
		"     len(chainCer.ParentChainCers)", len(chainCer.ParentChainCers),
		"     len(chainCertCer.ParentChainCers):", len(chainCertCer.ParentChainCers))

	parentChainCers := make([]ChainCertCer, 0)
	for i, _ := range chainCer.ParentChainCers {
		chainCertCer := ChainCertCer{Id: chainCer.ParentChainCers[i].Id}
		fmt.Println("NewChainCertCer():i chainCertCer:", i, jsonutil.MarshalJson(chainCertCer))
		parentChainCers = append(parentChainCers, chainCertCer)
		fmt.Println("NewChainCertCer():i, parentChainCers:",
			i, jsonutil.MarshalJson(parentChainCers), len(parentChainCers))
	}
	fmt.Println("NewChainCertCer():chainCertCer.Id:", chainCertCer.Id,
		"     len(chainCer.ParentChainCers)", len(chainCer.ParentChainCers),
		"     len(parentChainCers):", len(parentChainCers))
}
