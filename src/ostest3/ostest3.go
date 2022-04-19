package main

import (
	"fmt"

	"github.com/cpusoft/goutil/osutil"
)

func main() {
	dir := `rpstir2\source\branches\rpstir2-nodb`
	suffixs := make(map[string]string)
	suffixs[".go"] = ".go"
	files, err := osutil.GetAllFilesBySuffixs(dir, suffixs)
	fmt.Println(files, err)
}
