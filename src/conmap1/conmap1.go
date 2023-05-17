package main

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type Animal struct {
	Name string
}

func main() {
	cs := cmap.New[*Animal]()
	a1 := &Animal{Name: "a1"}
	cs.Set("a1", a1)
	b1 := &Animal{Name: "b1"}
	cs.Set("b1", b1)

	add(cs)
	for item := range cs.Iter(){
		fmt.Println( item.Val)
	}

}
func add(cs cmap.ConcurrentMap[string,*Animal]){
	a2 := &Animal{Name: "a2"}
	cs.Set("a2", a2)
	b2 := &Animal{Name: "b2"}
	cs.Set("b2", b2)
}