package main

import (
	"fmt"
	"sort"
)

type T struct {
	Foo int
	Bar int
}

// TVector is our basic vector type.
type TVector []T

func (v TVector) Len() int {
	return len(v)
}

func (v TVector) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// default comparison.
func (v TVector) Less(i, j int) bool {
	return v[i].Foo < v[j].Foo
}

type NotificationDelta struct {
	Serial uint64 `xml:"serial,attr"`
	Uri    string `xml:"uri,attr"`
	Hash   string `xml:"hash,attr"`
}
type NotificationDeltas []NotificationDelta

func (v NotificationDeltas) Len() int {
	return len(v)
}

func (v NotificationDeltas) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// default comparison.
func (v NotificationDeltas) Less(i, j int) bool {
	return v[i].Serial < v[j].Serial
}
func main() {
	v := []T{{1, 3}, {0, 6}, {3, 2}, {8, 7}}
	fmt.Println(v)
	sort.Sort(TVector(v))
	fmt.Println(v)

	v2 := []NotificationDelta{{Serial: 1, Uri: "3", Hash: "333"},
		{Serial: 0, Uri: "6", Hash: "666"},
		{Serial: 3, Uri: "2", Hash: "2222"},
		{Serial: 8, Uri: "7", Hash: "7777"}}
	fmt.Println(v2)
	sort.Sort(NotificationDeltas(v2))
	fmt.Println(v2)
}
