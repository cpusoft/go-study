package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("test main")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println("your input:", input.Text())
	cmd := exec.Command("nohup", "./fork1", ">", "./nohup.log", "2>&1", "&")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("fork1:err: ", err)
		return
	}
	if len(output) != 0 {
		fmt.Println("fork1:", output)
		return
	}

}
