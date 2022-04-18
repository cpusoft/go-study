package main

import (
	"bytes"
	"os"
	"time"

	belogs "github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/convert"
	tcputil "github.com/cpusoft/goutil/tcpserverclient/util"
)

func main() {
	t := `server`
	if len(os.Args) > 1 {
		t = os.Args[1]
	}
	belogs.Debug(t)
	if t == "tcpServer" {
		belogs.Debug("tcpServer")
		CreateTcpServer()
		select {}
	} else if t == "tlsServer" {
		belogs.Debug("tlsServer")
		CreateTlsServer()
		select {}
	}

}

func CreateTcpServer() {
	serverProcessFunc := new(ServerProcessFunc)
	ts := NewTcpServer(serverProcessFunc)
	ts.StartTcpServer("9999")
	time.Sleep(2 * time.Second)
	ts.ActiveSend(GetData(), "")
}
func CreateTlsServer() {
	serverProcessFunc := new(ServerProcessFunc)
	rootCrtFileName := `./ca/ca.cer`
	publicCrtFileName := `./server/server.cer`
	privateKeyFileName := `./server/serverkey.pem`

	ts := NewTlsServer(rootCrtFileName, publicCrtFileName, privateKeyFileName, true, serverProcessFunc)
	ts.StartTlsServer("9999")
	time.Sleep(2 * time.Second)
	ts.ActiveSend(GetData(), "")
}
func RtrProcess(receiveData []byte) (sendData []byte, err error) {
	buf := bytes.NewReader(receiveData)
	belogs.Debug("RtrProcess(): buf:", buf)
	return nil, nil
}
func GetData() (buffer []byte) {

	return []byte{0x00, 0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}

type ServerProcessFunc struct {
}

func (spf *ServerProcessFunc) OnConnectProcess(tcpTlsConn *TcpTlsConn) {

}
func (spf *ServerProcessFunc) ReceiveAndSendProcess(tcpTlsConn *TcpTlsConn, receiveData []byte) (nextConnectPolicy int, leftData []byte, err error) {
	belogs.Debug("ReceiveAndSendProcess(): len(receiveData):", len(receiveData), "   receiveData:", convert.Bytes2String(receiveData))
	// need recombine
	packets, leftData, err := tcputil.RecombineReceiveData(receiveData, PDU_TYPE_MIN_LEN, PDU_TYPE_LENGTH_START, PDU_TYPE_LENGTH_END)
	if err != nil {
		belogs.Error("ReceiveAndSendProcess(): RecombineReceiveData fail:", err)
		return tcputil.NEXT_CONNECT_POLICE_CLOSE_FORCIBLE, nil, err
	}
	belogs.Debug("ReceiveAndSendProcess(): RecombineReceiveData packets.Len():", packets.Len())

	if packets == nil || packets.Len() == 0 {
		belogs.Debug("ReceiveAndSendProcess(): RecombineReceiveData packets is empty:  len(leftData):", len(leftData))
		return tcputil.NEXT_CONNECT_POLICE_CLOSE_GRACEFUL, leftData, nil
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
			return tcputil.NEXT_CONNECT_POLICE_CLOSE_FORCIBLE, nil, err
		}

	}

	_, err = tcpTlsConn.Write(GetData())
	if err != nil {
		belogs.Error("ReceiveAndSendProcess(): tcp  Write fail:  tcpConn:", tcpTlsConn.RemoteAddr().String(), err)
		return tcputil.NEXT_CONNECT_POLICE_CLOSE_FORCIBLE, nil, err
	}
	// continue to receive next receiveData
	return tcputil.NEXT_CONNECT_POLICE_KEEP, leftData, nil
}
func (spf *ServerProcessFunc) OnCloseProcess(tcpTlsConn *TcpTlsConn) {

}
func (spf *ServerProcessFunc) ActiveSendProcess(tcpTlsConn *TcpTlsConn, sendData []byte) (err error) {
	return
}

const (
	PDU_TYPE_MIN_LEN      = 8
	PDU_TYPE_LENGTH_START = 4
	PDU_TYPE_LENGTH_END   = 8
)
