package main

import (
	"fmt"
)

type I interface {
	M(name string)
	M2() string
}
type T struct {
	name string
}

func (t *T) M(name string) {
	t.name = name
}
func (t *T) M2() string {
	return t.name
}
func main() {
	var i I = &T{"foo"}
	fmt.Println(i.(*T).name)
	f := I.M
	f(i, "bar")
	fmt.Println(i.(*T).name)
	fmt.Println(i, i.(*T))

	i.M("aaaa")
	fmt.Println(i.(*T).name)
	fmt.Println(i.M2())
}
