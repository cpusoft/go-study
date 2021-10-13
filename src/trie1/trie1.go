package main

import (
	"fmt"

	"github.com/derekparker/trie"
)

type NodeInfo struct {
	Address string
	Key     string
}

func main() {

	t := trie.New()
	ns := make([]NodeInfo, 0)
	ns = append(ns, NodeInfo{Address: "1.1.1.0/24", Key: "a"})
	ns = append(ns, NodeInfo{Address: "1.1.1.0/22", Key: "a"})
	ns = append(ns, NodeInfo{Address: "1.1.1.0/21", Key: "a"})
	ns = append(ns, NodeInfo{Address: "1.1.2.0/24", Key: "a"})
	ns = append(ns, NodeInfo{Address: "1.2.1.0/24", Key: "b"})
	ns = append(ns, NodeInfo{Address: "1.2.1.0/21", Key: "b"})
	ns = append(ns, NodeInfo{Address: "1.2.2.0/24", Key: "c"})
	/*
		t.Add("1.1.1.0", NodeInfo{Address: "1.1.1.0/24"})
		t.Add("1.1.1.0", NodeInfo{Address: "1.1.1.0/22"})
		t.Add("1.1.1.0", NodeInfo{Address: "1.1.1.0/21"})
		t.Add("1.1.2.0", NodeInfo{Address: "1.1.2.0/24"})
		t.Add("1.2.1.0", NodeInfo{Address: "1.2.1.0/24"})
		t.Add("1.2.2.0", NodeInfo{Address: "1.2.2.0/24"})
	*/
	for i := range ns {
		key := ns[i].Key
		node, ok := t.Find(key)
		if ok {
			meta := node.Meta()
			sames, ok2 := meta.([]NodeInfo)
			if !ok2 {
				fmt.Println("error1 ", sames, ok2)
				continue
			}
			sames = append(sames, ns[i])
			//belogs.Debug("buildTrie():find key, len(sameKeyRoaModels):", key, len(sameKeyRoaModels))
			t.Add(key, sames)
		} else {
			sames := make([]NodeInfo, 0)
			sames = append(sames, ns[i])
			t.Add(key, sames)
		}
	}
	keys1 := t.PrefixSearch("")
	fmt.Println(keys1)
	for _, key := range keys1 {
		node, ok := t.Find(key)
		meta := node.Meta()
		fmt.Println(key, ok, meta)
	}

}
