package main

import (
	"bytes"
	"fmt"
	"io"

	"os"
)

// get eecert from roa/mft
func main() {
	//openssl asn1parse -in E:\Go\parse-cert\data\-0AU6cJZAl7QHJeNhN9vE3zUBr4.roa -inform DER

	path := `E:\Go\go-study\src\eecert1\`
	file := path + `db42e932-926a-42bd-afdb-63320fa7ec40.roa`
	eeStart := 835579
	eeEnd := 1015426
	eeFile := path + "db42e932-926a-42bd-afdb-63320fa7ec40.ee.cer"

	//cerfileStr := path + `-ard.mft.cer`
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = f.Seek(int64(eeStart), 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	var buffer bytes.Buffer
	io.CopyN(&buffer, f, int64(eeEnd-eeStart))
	_bytes := buffer.Bytes()

	ee, err := os.Create(eeFile)

	defer ee.Close()
	ee.Write(_bytes)
}
