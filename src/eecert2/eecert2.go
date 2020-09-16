package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	/*
		path := `E:\Go\go-study\src\eecert2\`
		oldFile := path + `db42e932-926a-42bd-afdb-63320fa7ec40.roa`
		eeStart := 838969
		eeEnd := 1019659
		newFile := path + "db42e932-926a-42bd-afdb-63320fa7ec40.ee.cer"
	*/
	path := `G:\Download\cert\`
	oldFile := path + `f69cda529c78bcae844cfa0cd9ed17830658c7dc.roa`
	eeStart := 139
	eeEnd := 1505
	newFile := path + "f69cda529c78bcae844cfa0cd9ed17830658c7dc.ee.cer"

	oldFileByte, err := ioutil.ReadFile(oldFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	newFileByte := oldFileByte[eeStart:eeEnd]

	newF, err := os.Create(newFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer newF.Close()
	newF.Write(newFileByte)
}
