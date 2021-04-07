package main

import (
	"fmt"

	conf1 "github.com/cpusoft/goutil/conf"
)

func main() {
	value := conf1.VariableString("rrdp::destpath")
	fmt.Println(value)

	err := conf1.Set("name", "astaxie")
	fmt.Println(err)

}
