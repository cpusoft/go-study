package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	_ "golang.org/x/net/websocket"
)

type MyNetUtil struct {
	NetUtil
}

func (t MyNetUtil) HandleClient(Conn net.Conn) {
	fmt.Println("MyNetUtil handle the conn")
	Conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	Conn.SetWriteDeadline(time.Now().Add(2 * time.Minute))

	defer Conn.Close()

	request := make([]byte, 128)
	for {
		readLine, err := Conn.Read(request)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(readLine)
		fmt.Println(string(request))

		if readLine == 0 {
			break
		} else if strings.TrimSpace(string(request[:readLine])) == "t" {
			daytime := strconv.FormatInt(time.Now().Unix(), 10) + "\r\n"
			Conn.Write([]byte(daytime))
		} else {
			daytime := time.Now().String() + "\r\n"
			Conn.Write([]byte(daytime))
		}
		request = make([]byte, 128)
	}
}

func main() {
	myNetUtil := MyNetUtil{}
	var netUtilInterface NetUtilInterface
	netUtilInterface = myNetUtil
	StartTCPServer(":7777", netUtilInterface)
}
