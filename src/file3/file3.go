package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	//cerfile, err := os.Create(cerfileStr)
	f, err := ioutil.TempFile("", "test")
	fmt.Println(f, err)
	fmt.Println(f.Name())
	defer f.Close()
	os.Remove(f.Name())

}
