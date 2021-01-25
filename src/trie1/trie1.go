package main

import (
	"fmt"

	"github.com/derekparker/trie"
)

type NodeInfo struct {
	Address string
}

func main() {

	t := trie.New()
	t.Add("1.1.1.0", NodeInfo{Address: "1.1.1.0/24"})
	t.Add("1.1.1.0", NodeInfo{Address: "1.1.1.0/22"})
	t.Add("1.1.1.0", NodeInfo{Address: "1.1.1.0/21"})
	t.Add("1.1.2.0", NodeInfo{Address: "1.1.2.0/24"})
	t.Add("1.2.1.0", NodeInfo{Address: "1.2.1.0/24"})
	t.Add("1.2.2.0", NodeInfo{Address: "1.2.2.0/24"})

	keys1 := t.PrefixSearch("1.1")
	fmt.Println(keys1)
	for _, key := range keys1 {
		node, ok := t.Find(key)
		meta := node.Meta()
		fmt.Println(key, ok, meta)
	}
}
