package main

import (
	"bytes"
	"net"
	"os"
	"time"

	belogs "github.com/cpusoft/goutil/belogs"
	_ "github.com/cpusoft/goutil/conf"
	"github.com/cpusoft/goutil/convert"
	_ "github.com/cpusoft/goutil/logs"
)

const (
	PDU_TYPE_MIN_LEN      = 8
	PDU_TYPE_LENGTH_START = 4
	PDU_TYPE_LENGTH_END   = 8
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

type ServerProcessFunc struct {
}

func (spf *ServerProcessFunc) OnConnectProcess(tcpConn *net.TCPConn) {

}
func (spf *ServerProcessFunc) ReceiveAndSendProcess(tcpConn *net.TCPConn, receiveData []byte) (nextConnectPolicy int, leftData []byte, err error) {
	belogs.Debug("ReceiveAndSendProcess(): len(receiveData):", len(receiveData), "   receiveData:", convert.Bytes2String(receiveData))
	// need recombine
	packets, leftData, err := RecombineReceiveData(receiveData, PDU_TYPE_MIN_LEN, PDU_TYPE_LENGTH_START, PDU_TYPE_LENGTH_END)
	if err != nil {
		belogs.Error("ReceiveAndSendProcess(): RecombineReceiveData fail:", err)
		return NEXT_CONNECT_POLICE_CLOSE_FORCIBLE, nil, err
	}
	belogs.Debug("ReceiveAndSendProcess(): RecombineReceiveData packets.Len():", packets.Len())

	if packets == nil || packets.Len() == 0 {
		belogs.Debug("ReceiveAndSendProcess(): RecombineReceiveData packets is empty:  len(leftData):", len(leftData))
		return NEXT_CONNECT_POLICE_CLOSE_GRACEFUL, leftData, nil
	}
	for e := packets.Front(); e != nil; e = e.Next() {
		packet, ok := e.Value.([]byte)
		if !ok || packet == nil || len(packet) == 0 {
			belogs.Debug("ReceiveAndSendProcess(): for packets fail:", convert.ToString(e.Value))
			break
		}
		_, err := RtrProcess(packet)
		if err != nil {
			belogs.Error("ReceiveAndSendProcess(): RtrProcess fail:", err)
			return NEXT_CONNECT_POLICE_CLOSE_FORCIBLE, nil, err
		}

	}

	_, err = tcpConn.Write(GetData())
	if err != nil {
		belogs.Error("ReceiveAndSendProcess(): tcp  Write fail:  tcpConn:", tcpConn.RemoteAddr().String(), err)
		return NEXT_CONNECT_POLICE_CLOSE_FORCIBLE, nil, err
	}
	// continue to receive next receiveData
	return NEXT_CONNECT_POLICE_KEEP, leftData, nil
}
func (spf *ServerProcessFunc) OnCloseProcess(tcpConn *net.TCPConn) {

}
func (spf *ServerProcessFunc) ActiveSendProcess(tcpConn *net.TCPConn, sendData []byte) (err error) {
	return
}

func CreateTcpServer() {
	serverProcessFunc := new(ServerProcessFunc)
	ts := NewTcpServer(serverProcessFunc)
	ts.Start("0.0.0.0:9999")
	time.Sleep(2 * time.Second)
	ts.ActiveSend()
}

func RtrProcess(receiveData []byte) (sendData []byte, err error) {
	buf := bytes.NewReader(receiveData)
	belogs.Debug("RtrProcess(): buf:", buf)
	return nil, nil
}

type ClientProcessFunc struct {
}

func (cp *ClientProcessFunc) OnConnectProcess(tcpConn *net.TCPConn) {
	// can active read here

}
func (cp *ClientProcessFunc) OnCloseProcess(tcpConn *net.TCPConn) {
}

func (sq *ClientProcessFunc) OnReceiveProcess(tcpConn *net.TCPConn, receiveData []byte) (nextRwPolicy int, leftData []byte, err error) {

	belogs.Debug("OnReceiveProcess :", tcpConn, convert.Bytes2String(receiveData))
	return NEXT_RW_POLICE_END_READ, make([]byte, 0), nil
}

func CreateTcpClient() {
	clientProcessFunc := new(ClientProcessFunc)

	//CreateTcpClient("127.0.0.1:9999", ClientProcess1)
	tc := NewTcpClient("stop", clientProcessFunc)
	err := tc.Start("192.168.83.139:9999")
	belogs.Debug("tc:", tc, err)
	if err != nil {
		return
	}
	belogs.Debug("will SendData")
	tcpClientMsg := &TcpClientMsg{NextConnectClosePolicy: NEXT_CONNECT_POLICE_KEEP,
		NextRwPolice: NEXT_RW_POLICE_WAIT_READ,
		SendData:     GetData(),
	}
	tc.SendMsg(tcpClientMsg)
	time.Sleep(60 * time.Second)

	belogs.Debug("will stop")
	tcpClientMsg.NextConnectClosePolicy = NEXT_CONNECT_POLICE_CLOSE_GRACEFUL
	tcpClientMsg.SendData = nil
	tc.SendMsg(tcpClientMsg)
	time.Sleep(60 * time.Second)
}
