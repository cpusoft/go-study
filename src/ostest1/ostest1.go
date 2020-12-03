package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("os.Args[0]:", os.Args[0])
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("filepath.Dir:", dir)

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dirs := strings.Split(path, string(os.PathSeparator))
	index := len(dirs)
	if len(dirs) > 2 {
		index = len(dirs) - 2
	}
	ret := strings.Join(dirs[:index], string(os.PathSeparator))
	fmt.Println("parent ret:", ret)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pwd:", pwd)

	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println("executable  filepath.Dir:", ex, exPath)

	_, filename, _, _ := runtime.Caller(1)
	fmt.Println("Caller:", filename)
}
