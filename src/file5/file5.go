package main

import (
	"fmt"
	"os"
)

func main() {
	f := `G:\Download\cert\5zp5RARE_jTDEPeEiUxKsh6sqgQ.mft`
	b, err := os.ReadFile(f)
	fmt.Println(len(b), err)
}
