package main

import (
	"os"

	belogs "github.com/cpusoft/goutil/belogs"
	_ "github.com/cpusoft/goutil/conf"
	_ "github.com/cpusoft/goutil/logs"
)

func main() {
	t := `server`
	if len(os.Args) > 1 {
		t = os.Args[1]
	}
	belogs.Debug(t)
	if t == "server" {
		belogs.Debug("server")
		CreateTcpServer()
		select {}
	} else if t == "client" {
		belogs.Debug("client")
		CreateTcpClient()
	}

}
func GetData() (buffer []byte) {

	return []byte{0x00, 0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}
