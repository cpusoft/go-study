package main

import (
	"bytes"
	"net"
	"time"

	belogs "github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/convert"
	tcpserver "github.com/cpusoft/goutil/tcpserverclient/tcpserver"
	tcputil "github.com/cpusoft/goutil/tcpserverclient/util"
)

func CreateTcpServer() {
	serverProcessFunc := new(ServerProcessFunc)
	ts := tcpserver.NewTcpServer(serverProcessFunc)
	ts.Start("0.0.0.0:9999")
	time.Sleep(2 * time.Second)
	ts.ActiveSend(GetData(), "")
}

func RtrProcess(receiveData []byte) (sendData []byte, err error) {
	buf := bytes.NewReader(receiveData)
	belogs.Debug("RtrProcess(): buf:", buf)
	return nil, nil
}

type ServerProcessFunc struct {
}

func (spf *ServerProcessFunc) OnConnectProcess(tcpConn *net.TCPConn) {

}
func (spf *ServerProcessFunc) ReceiveAndSendProcess(tcpConn *net.TCPConn, receiveData []byte) (nextConnectPolicy int, leftData []byte, err error) {
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

	_, err = tcpConn.Write(GetData())
	if err != nil {
		belogs.Error("ReceiveAndSendProcess(): tcp  Write fail:  tcpConn:", tcpConn.RemoteAddr().String(), err)
		return tcputil.NEXT_CONNECT_POLICE_CLOSE_FORCIBLE, nil, err
	}
	// continue to receive next receiveData
	return tcputil.NEXT_CONNECT_POLICE_KEEP, leftData, nil
}
func (spf *ServerProcessFunc) OnCloseProcess(tcpConn *net.TCPConn) {

}
func (spf *ServerProcessFunc) ActiveSendProcess(tcpConn *net.TCPConn, sendData []byte) (err error) {
	return
}

const (
	PDU_TYPE_MIN_LEN      = 8
	PDU_TYPE_LENGTH_START = 4
	PDU_TYPE_LENGTH_END   = 8
)
