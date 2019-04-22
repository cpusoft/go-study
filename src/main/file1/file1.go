package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	//openssl asn1parse -in E:\Go\parse-cert\data\-0AU6cJZAl7QHJeNhN9vE3zUBr4.roa -inform DER

	path := `E:\Go\parse-cert\data\`
	/*
		roafileStr := path + `-0AU6cJZAl7QHJeNhN9vE3zUBr4.roa`
		cerfileStr := path + `-0AU6cJZAl7QHJeNhN9vE3zUBr4.roa_create.cer`
		roafile, err := os.OpenFile(roafileStr, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer roafile.Close()

		roafile.Seek(89, 0)
		fmt.Println("Success Open ROA File")
		var buffer bytes.Buffer
		io.CopyN(&buffer, roafile, 1329-89)
		_bytes := buffer.Bytes()

		for _, by := range _bytes {
			fmt.Printf("%02X ", by)
		}

		cerfile, err := os.Create(cerfileStr)
		defer cerfile.Close()
		cerfile.Write(_bytes)
		buffer.Reset()
	*/

	mftfileStr := path + `-ard.mft`
	//cerfileStr := path + `-ard.mft.cer`
	mftfile, err := os.OpenFile(mftfileStr, os.O_RDONLY, os.ModePerm)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer mftfile.Close()

	mftfile.Seek(188, 0)
	fmt.Println("Success Open MFT File")
	var mftBuffer bytes.Buffer
	io.CopyN(&mftBuffer, mftfile, 1567-188)
	_bytes := mftBuffer.Bytes()
	for _, by := range _bytes {
		fmt.Printf("%02X ", by)
	}

	//cerfile, err := os.Create(cerfileStr)
	cerfile, err := ioutil.TempFile("", "-ard.mft.cer")
	defer cerfile.Close()
	cerfile.Write(_bytes)
}
