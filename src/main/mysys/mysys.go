package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	shell         = syscall.MustLoadDLL("Shell32.dll")
	getFolderPath = shell.MustFindProc("SHGetFolderPathW")
)

const (
	CSIDL_DESKTOP = 0
	CSIDL_APPDATA = 26
)

func main() {
	b := make([]uint16, syscall.MAX_PATH)
	fmt.Print(b)
	// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762181%28v=vs.85%29.aspx
	r, _, err := getFolderPath.Call(0, CSIDL_DESKTOP, 0, 0, uintptr(unsafe.Pointer(&b[0])))
	if uint32(r) != 0 {
		fmt.Sprintf("", err)
	}
	a_dir := syscall.UTF16ToString(b)

	r, _, err = getFolderPath.Call(0, CSIDL_APPDATA, 0, 0, uintptr(unsafe.Pointer(&b[0])))
	if uint32(r) != 0 {
		fmt.Sprintf(err.Error())
	}
	b_dir := syscall.UTF16ToString(b)

	fmt.Printf("%d %s\n", CSIDL_DESKTOP, a_dir)
	fmt.Printf("%d  %s\n", CSIDL_APPDATA, b_dir)
}
