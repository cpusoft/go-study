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
	oldFile := path + `C1A58E56F2D411EAAA2A3738C4F9AE02.roa`
	eeStart := 420
	eeEnd := 1837
	newFile := path + "C1A58E56F2D411EAAA2A3738C4F9AE02.ee.cer"

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
	_, err = newF.Write(newFileByte)
	fmt.Println(err)
}
