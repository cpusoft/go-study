package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cpusoft/goutil/fileutil"
)

func main() {
	if len(os.Args) <= 1 {
		return
	}
	f := os.Args[1]
	fmt.Println(f)
	cmd := exec.Command("python", "analyze_roa_comp.py", f)
	// if success, the len(output) will be zero
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error:", err)
		fmt.Println(string(output))
		return
	}
	fmt.Println(string(output))
	fileutil.WriteBytesToFile("./out.html", output)
	return
}
