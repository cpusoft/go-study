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

func main() {
	v := []T{{1, 3}, {0, 6}, {3, 2}, {8, 7}}
	fmt.Println(v)
	sort.Sort(TVector(v))
	fmt.Println(v)
}
