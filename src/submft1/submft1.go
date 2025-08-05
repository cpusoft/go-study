package main

import (
	"fmt"

	"github.com/cpusoft/goutil/fileutil"
)

func main() {

	//files := []string{`G:\Download\cert\C5zKkN0Neoo3ZmsZIX_g2EA3t6I.mft`}
	files := []string{`G:\Download\cert\WpIYL0DvGR69MNqQHSOeKTTXrTw.roa`}
	for _, file := range files {
		fmt.Println(file)
		b, err := fileutil.ReadFileToBytes(file)
		if err != nil {
			fmt.Println(file, err)
			continue
		}
		eeCertStart := 138
		eeCertEnd := 1452
		nb := b[eeCertStart:eeCertEnd]
		fileutil.WriteBytesToFile(`G:\Download\cert\WpIYL0DvGR69MNqQHSOeKTTXrTw-roaee.cer`, nb)
	}
}
