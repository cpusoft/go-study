package netutil

import (
	"errors"
	"fmt"
	"net"
	"os"
)

func Division(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}

	return a / b, nil
}

type NetUtil struct {
}
type NetUtilInterface interface {
	HandleClient(net.Conn)
}

func (t NetUtil) HandleClient(Conn net.Conn) {
	fmt.Println("do nothing in NetUtil")
	defer Conn.Close()
}
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
func StartTCPServer(addrAndPort string, t NetUtilInterface) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addrAndPort)
	CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go t.HandleClient(conn)
	}
}
