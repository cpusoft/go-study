package main

import (
	"bytes"
	"os"
	"time"

	belogs "github.com/cpusoft/goutil/belogs"
	_ "github.com/cpusoft/goutil/conf"
	"github.com/cpusoft/goutil/convert"
	_ "github.com/cpusoft/goutil/logs"
)

func main() {
	t := `tcpServer`
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
	} else if t == "tcpClient" {
		belogs.Debug("tcpClient")
		CreateTcpClient()
		select {}
	} else if t == "tlsClient" {
		belogs.Debug("tlsClient")
		CreateTlsClient()
		select {}
	}

}

func CreateTcpServer() {
	serverProcessFunc := new(ServerProcessFunc)
	ts := NewTcpServer(serverProcessFunc)
	belogs.Debug("CreateTcpServer():", 9999)
	err := ts.StartTcpServer("9999")
	if err != nil {
		belogs.Error("CreateTcpServer(): StartTcpServer ts fail: ", &ts, err)
		return
	}
	time.Sleep(2 * time.Second)
	ts.ActiveSend(GetData(), "")

	time.Sleep(5 * time.Second)
	ts.CloseGraceful()
}
func CreateTlsServer() {
	serverProcessFunc := new(ServerProcessFunc)
	tlsRootCrtFileName := `ca.cer`
	tlsPublicCrtFileName := `server.cer`
	tlsPrivateKeyFileName := `serverkey.pem`
	belogs.Debug("CreateTlsServer(): tlsRootCrtFileName:", tlsRootCrtFileName,
		"tlsPublicCrtFileName:", tlsPublicCrtFileName,
		"tlsPrivateKeyFileName:", tlsPrivateKeyFileName)

	ts, err := NewTlsServer(tlsRootCrtFileName, tlsPublicCrtFileName, tlsPrivateKeyFileName, true, serverProcessFunc)
	if err != nil {
		belogs.Error("CreateTlsServer(): NewTlsServer ts fail: ", &ts, err)
		return
	}
	go ts.StartTlsServer("9999")

	time.Sleep(5 * time.Second)
	ts.ActiveSend(GetData(), "")
	time.Sleep(8 * time.Second)
	ts.CloseGraceful()
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

	_, err = tcpTlsConn.Write(GetData())
	if err != nil {
		belogs.Error("ReceiveAndSendProcess(): tcp  Write fail:  tcpTlsConn:", tcpTlsConn.RemoteAddr().String(), err)
		return NEXT_CONNECT_POLICE_CLOSE_FORCIBLE, nil, err
	}
	// continue to receive next receiveData
	return NEXT_CONNECT_POLICE_KEEP, leftData, nil
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

func CreateTcpClient() {
	clientProcessFunc := new(ClientProcessFunc)
	belogs.Debug("CreateTcpClient():", "192.168.83.139:9999")
	//CreateTcpClient("127.0.0.1:9999", ClientProcess1)
	tc := NewTcpClient(clientProcessFunc)
	err := tc.StartTcpClient("192.168.83.139:9999")
	if err != nil {
		belogs.Error("CreateTcpClient(): StartTcpClient tc fail: ", &tc, err)
		return
	}
	belogs.Debug("CreateTcpClient(): tcpclient will SendData")
	tcpClientSendMsg := &TcpTlsClientSendMsg{NextConnectClosePolicy: NEXT_CONNECT_POLICE_KEEP,
		NextRwPolice: NEXT_RW_POLICE_WAIT_READ,
		SendData:     GetTcpClientData(),
	}
	tc.SendMsg(tcpClientSendMsg)
	time.Sleep(60 * time.Second)

	belogs.Debug("CreateTcpClient(): tcpclient will stop")
	tc.CloseGraceful()

}

func CreateTlsClient() {
	clientProcessFunc := new(ClientProcessFunc)
	tlsRootCrtFileName := `ca.cer`
	tlsPublicCrtFileName := `client.cer`
	tlsPrivateKeyFileName := `clientkey.pem`
	belogs.Debug("CreateTlsClient(): tlsRootCrtFileName:", tlsRootCrtFileName,
		"tlsPublicCrtFileName:", tlsPublicCrtFileName,
		"tlsPrivateKeyFileName:", tlsPrivateKeyFileName)
	//CreateTcpClient("192.168.83.139:9999", ClientProcess1)
	tc, err := NewTlsClient(tlsRootCrtFileName, tlsPublicCrtFileName, tlsPrivateKeyFileName, clientProcessFunc)
	if err != nil {
		belogs.Error("CreateTcpClient(): NewTlsClient tc fail: ", &tc, err)
		return
	}
	err = tc.StartTlsClient("192.168.83.139:9999")
	if err != nil {
		belogs.Error("CreateTcpClient(): StartTlsClient tc fail: ", &tc, err)
		return
	}
	belogs.Debug("CreateTcpClient(): tcpclient will SendData")
	tcpClientSendMsg := &TcpTlsClientSendMsg{NextConnectClosePolicy: NEXT_CONNECT_POLICE_KEEP,
		NextRwPolice: NEXT_RW_POLICE_WAIT_READ,
		SendData:     GetTcpClientData(),
	}
	tc.SendMsg(tcpClientSendMsg)
	time.Sleep(60 * time.Second)

	belogs.Debug("CreateTcpClient(): tcpclient will stop")
	tcpClientSendMsg.NextConnectClosePolicy = NEXT_CONNECT_POLICE_CLOSE_GRACEFUL
	tcpClientSendMsg.SendData = nil
	tc.SendMsg(tcpClientSendMsg)

}

type ClientProcessFunc struct {
}

func (cp *ClientProcessFunc) OnConnectProcess(tcpTlsConn *TcpTlsConn) {

	belogs.Info("OnConnectProcess(): tcpclient tcpTlsConn:", tcpTlsConn.RemoteAddr().String())

}
func (cp *ClientProcessFunc) OnCloseProcess(tcpTlsConn *TcpTlsConn) {
	if tcpTlsConn != nil {
		belogs.Info("OnCloseProcess(): tcpclient tcpTlsConn:", tcpTlsConn.RemoteAddr().String())
	}
}

func (sq *ClientProcessFunc) OnReceiveProcess(tcpTlsConn *TcpTlsConn, receiveData []byte) (nextRwPolicy int, leftData []byte, err error) {

	belogs.Debug("OnReceiveProcess() tcpclient  :", tcpTlsConn, convert.Bytes2String(receiveData))
	return NEXT_RW_POLICE_END_READ, make([]byte, 0), nil
}

func GetTcpClientData() (buffer []byte) {

	return []byte{0x00, 0x0b, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x0b, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}
