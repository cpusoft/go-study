package main

import (
	"net"
	"time"

	belogs "github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/convert"
	tcpclient "github.com/cpusoft/goutil/tcpserverclient/tcpclient"
	tcputil "github.com/cpusoft/goutil/tcpserverclient/util"
)

type ClientProcessFunc struct {
}

func (cp *ClientProcessFunc) OnConnectProcess(tcpConn *net.TCPConn) {
	// can active read here

}
func (cp *ClientProcessFunc) OnCloseProcess(tcpConn *net.TCPConn) {
}

func (sq *ClientProcessFunc) OnReceiveProcess(tcpConn *net.TCPConn, receiveData []byte) (nextRwPolicy int, leftData []byte, err error) {

	belogs.Debug("OnReceiveProcess :", tcpConn, convert.Bytes2String(receiveData))
	return tcputil.NEXT_RW_POLICE_END_READ, make([]byte, 0), nil
}

func CreateTcpClient() {
	clientProcessFunc := new(ClientProcessFunc)

	//CreateTcpClient("127.0.0.1:9999", ClientProcess1)
	tc := tcpclient.NewTcpClient(clientProcessFunc)
	err := tc.Start("192.168.83.139:9999")
	belogs.Debug("tc:", tc, err)
	if err != nil {
		return
	}
	belogs.Debug("will SendData")
	tcpClientMsg := &tcpclient.TcpClientMsg{NextConnectClosePolicy: tcputil.NEXT_CONNECT_POLICE_KEEP,
		NextRwPolice: tcputil.NEXT_RW_POLICE_WAIT_READ,
		SendData:     GetData(),
	}
	tc.SendMsg(tcpClientMsg)
	time.Sleep(60 * time.Second)

	belogs.Debug("will stop")
	tcpClientMsg.NextConnectClosePolicy = tcputil.NEXT_CONNECT_POLICE_CLOSE_GRACEFUL
	tcpClientMsg.SendData = nil
	tc.SendMsg(tcpClientMsg)
	time.Sleep(60 * time.Second)
}
