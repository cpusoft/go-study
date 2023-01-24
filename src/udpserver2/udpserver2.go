package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/transportutil"
)

type ServerProcess struct {
	transportConnsMutex sync.RWMutex
	TransportConns      map[string]*transportutil.TransportConn
	transportMsg        chan transportutil.TransportMsg
}

func NewServerProcess(transportMsg chan transportutil.TransportMsg) *ServerProcess {
	c := &ServerProcess{}
	c.TransportConns = make(map[string]*transportutil.TransportConn, 16)
	c.transportMsg = transportMsg

	return c
}
func (c *ServerProcess) OnConnectProcess(transportConn *transportutil.TransportConn) {
	fmt.Println("OnConnectProcess")
}
func (c *ServerProcess) OnReceiveAndSendProcess(transportConn *transportutil.TransportConn,
	receiveData []byte) (nextConnectPolicy int, leftData []byte, err error) {
	fmt.Println("OnReceiveAndSendProcess():", convert.PrintBytesOneLine(receiveData))
	// 发送数据
	len, err := transportConn.Write([]byte("from server"))
	fmt.Println("OnReceiveAndSendProcess(): len:", len, err)
	return
}
func (c *ServerProcess) OnCloseProcess(transportConn *transportutil.TransportConn) {
	fmt.Println("OnCloseProcess(): server tcptlsserver tcpTlsConn: ", transportConn.RemoteAddr().String())

}

func main() {
	transportMsg := make(chan transportutil.TransportMsg, 15)
	fmt.Println("main(): transportMsg:", transportMsg)

	// process
	serverProcess := NewServerProcess(transportMsg)
	fmt.Println("main(): serverProcess:", serverProcess)

	ts := transportutil.NewUdpServer(serverProcess, transportMsg)
	go ts.StartUdpServer("9998")
	time.Sleep(5000 * time.Second)

}
