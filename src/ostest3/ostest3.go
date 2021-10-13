package main

import (
	"fmt"

	"github.com/cpusoft/goutil/osutil"
)

func main() {
	dir := `F:\share\我的坚果云\Go\rpstir2\source\branches\rpstir2-nodb`
	suffixs := make(map[string]string)
	suffixs[".go"] = ".go"
	files, err := osutil.GetAllFilesBySuffixs(dir, suffixs)
	fmt.Println(files, err)
}
