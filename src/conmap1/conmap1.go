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
	cs.Set("1", a1)
	b1 := &Animal{Name: "b1"}
	cs.Set("1", b1)

	add(cs)
	for item := range cs.Iter(){
		fmt.Println( item.Key)
		fmt.Println( item.Val)
	}

}
func add(cs cmap.ConcurrentMap[string,*Animal]){
	a2 := &Animal{Name: "a2"}
	cs.Set("2", a2)
	b2 := &Animal{Name: "b2"}
	cs.Set("2", b2)
}