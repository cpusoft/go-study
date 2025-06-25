package main

import (
	"fmt"

	"github.com/cpusoft/goutil/fileutil"
)

func main() {

	//files := []string{`G:\Download\cert\C5zKkN0Neoo3ZmsZIX_g2EA3t6I.mft`}
	files := []string{`G:\Download\cert\006fd0a7-9256-4454-93a2-2c0167d518bc.mft`}
	for _, file := range files {
		fmt.Println(file)
		b, err := fileutil.ReadFileToBytes(file)
		if err != nil {
			fmt.Println(file, err)
			continue
		}
		eeCertStart := 461
		eeCertEnd := 2093
		nb := b[eeCertStart:eeCertEnd]
		fileutil.WriteBytesToFile(`G:\Download\cert\006fd0a7-9256-4454-93a2-2c0167d518bc-EE.mft`, nb)
	}
}
