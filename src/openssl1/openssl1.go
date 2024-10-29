package main

import (
	"fmt"
	"os/exec"
)

func main() {
	opensslCmd := `/home/openssl/openssl/bin/openssl`
	cmd := exec.Command(opensslCmd, "version")

	output, err := cmd.CombinedOutput()
	fmt.Println(string(output), err)
}
