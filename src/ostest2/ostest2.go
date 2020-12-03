package main

import (
	"errors"
	"fmt"

	"os"

	"strings"
)

func GetPwd() string {
	pwd, _ := os.Getwd()
	return pwd
}
func GetPathSeparator() string {
	return string(os.PathSeparator)
}

// judge file is dir or not.
func IsDir(file string) (bool, error) {
	s, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	return s.IsDir(), nil
}

func IsFile(file string) (bool, error) {
	s, err := IsDir(file)
	return !s, err
}

func GetRelativePathInCurrentOrParentAbsolutePath(relativePath string) (absolutePath string, err error) {
	path := GetPwd()
	absolutePath = path + GetPathSeparator() + relativePath
	fmt.Println("path,absolutePath:", path, absolutePath)
	ok, err := IsDir(absolutePath)
	fmt.Println("is dir absolutePath:", absolutePath, ok, err)
	if err == nil && ok {
		return absolutePath, nil
	}
	pos := strings.LastIndex(path, GetPathSeparator())
	path = string([]byte(path)[:pos])
	fmt.Println("pos,path:", pos, path)
	absolutePath = path + GetPathSeparator() + relativePath
	ok, err = IsDir(absolutePath)
	fmt.Println("2 is dir absolutePath:", absolutePath, ok, err)
	if err == nil && ok {
		return absolutePath, nil
	}
	return "", errors.New("cannot found absolutePath of relativePath " + relativePath)

}

func main() {

	f, err := GetRelativePathInCurrentOrParentAbsolutePath("conf")
	fmt.Println(f, err)

}
