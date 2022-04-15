package main

import (
	"io"
	"net"
	"sync"
	"time"

	belogs "github.com/cpusoft/goutil/belogs"
)

//https://github.com/SmarkSeven/socket/blob/master/protocol.go

type TcpServer struct {
	tcpConns             map[string]*net.TCPConn // map[addr]*net.TCPConn
	tcpConnsMutex        sync.RWMutex
	tcpServerProcessFunc TcpServerProcessFunc
}
type TcpServerProcessFunc interface {
	OnConnectProcess(tcpConn *net.TCPConn)
	ReceiveAndSendProcess(tcpConn *net.TCPConn, receiveData []byte) (nextConnectPolicy int, leftData []byte, err error)
	OnCloseProcess(tcpConn *net.TCPConn)
	ActiveSendProcess(tcpConn *net.TCPConn, sendData []byte) (err error)
}

// server: 0.0.0.0:port
func NewTcpServer(tcpServerProcessFunc TcpServerProcessFunc) (ts *TcpServer) {

	belogs.Debug("NewTcpServer():tcpProcessFunc:", tcpServerProcessFunc)
	ts = &TcpServer{}
	ts.tcpConns = make(map[string]*net.TCPConn, 16)
	ts.tcpServerProcessFunc = tcpServerProcessFunc
	belogs.Debug("NewTcpServer():ts:", ts, "   ts.tcpConnsMutex:", ts.tcpConnsMutex)
	return ts
}

// server: 0.0.0.0:port
func (ts *TcpServer) Start(server string) (err error) {
	tcpServer, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		belogs.Error("Start(): tcp  ResolveTCPAddr fail: ", server, err)
		return err
	}

	listen, err := net.ListenTCP("tcp", tcpServer)
	if err != nil {
		belogs.Error("Start(): tcp  ListenTCP fail: ", server, err)
		return err
	}
	defer listen.Close()
	belogs.Debug("Start(): tcp  create server ok, server is ", server, "  will accept client")

	for {
		tcpConn, err := listen.AcceptTCP()
		belogs.Info("Start(): tcp  Accept remote: ", tcpConn.RemoteAddr().String())
		if err != nil {
			belogs.Error("Start(): tcp  Accept remote fail: ", server, tcpConn.RemoteAddr().String(), err)
			continue
		}
		if tcpConn == nil {
			continue
		}

		ts.OnConnect(tcpConn)

		// call func to process tcpConn
		go ts.ReceiveAndSend(tcpConn)

	}
	return nil
}

func (ts *TcpServer) OnConnect(tcpConn *net.TCPConn) {
	start := time.Now()
	belogs.Debug("OnConnect(): tcp new tcpConn: ", tcpConn)

	// add new tcpConn to tcpconns
	ts.tcpConnsMutex.Lock()
	defer ts.tcpConnsMutex.Unlock()
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(time.Second * 300)
	connKey := tcpConn.RemoteAddr().String()
	ts.tcpConns[connKey] = tcpConn
	belogs.Debug("OnConnect():tcp tcpConn: ", tcpConn.RemoteAddr().String(), ", connKey:", connKey, "  new len(tcpConns): ", len(ts.tcpConns))

	// call process func OnConnect
	belogs.Debug("OnConnect():tcp tcpConn: ", tcpConn.RemoteAddr().String(), "   call process func: OnConnect ")
	ts.tcpServerProcessFunc.OnConnectProcess(tcpConn)
	belogs.Info("OnConnect(): tcp add tcpConn: ", tcpConn.RemoteAddr().String(), "   len(tcpConns): ", len(ts.tcpConns), "   time(s):", time.Now().Sub(start).Seconds())
}

func (ts *TcpServer) ReceiveAndSend(tcpConn *net.TCPConn) {

	defer ts.OnClose(tcpConn)

	var leftData []byte
	// one packet
	buffer := make([]byte, 2048)
	// wait for new packet to read
	for {
		n, err := tcpConn.Read(buffer)
		start := time.Now()
		belogs.Debug("ReceiveAndSend(): tcp server read: Read n: ", tcpConn.RemoteAddr().String(), n)
		if err != nil {
			if err == io.EOF {
				// is not error, just client close
				belogs.Debug("ReceiveAndSend(): tcp server Read io.EOF, client close: ", tcpConn.RemoteAddr().String(), err)
				return
			}
			belogs.Error("ReceiveAndSend(): tcp server Read fail, err ", tcpConn.RemoteAddr().String(), err)
			return
		}
		if n == 0 {
			continue
		}

		// call process func OnReceiveAndSend
		// copy to leftData
		belogs.Debug("ReceiveAndSend(): tcp  will ReceiveAndSendProcess, server tcpConn: ", tcpConn.RemoteAddr().String(), "  n:", n,
			" , will call process func: OnReceiveAndSend,  time(s):", time.Now().Sub(start))
		nextConnectPolicy, leftData, err := ts.tcpServerProcessFunc.ReceiveAndSendProcess(tcpConn, append(leftData, buffer[:n]...))
		belogs.Debug("ReceiveAndSend(): tcp  after ReceiveAndSendProcess,server tcpConn: ", tcpConn.RemoteAddr().String(), " receive n: ", n,
			"  len(leftData):", len(leftData), "  time(s):", time.Now().Sub(start))
		if err != nil {
			belogs.Error("OnReceiveAndSend():tcp server fail ,will remove this tcpConn : ", tcpConn.RemoteAddr().String(), err)
			return
		}
		if nextConnectPolicy == NEXT_CONNECT_POLICE_CLOSE_GRACEFUL || nextConnectPolicy == NEXT_CONNECT_POLICE_CLOSE_FORCIBLE {
			belogs.Debug("OnReceiveAndSend():nextConnectPolicy return : ", tcpConn.RemoteAddr().String(), nextConnectPolicy)
			return
		}
	}
}

func (ts *TcpServer) OnClose(tcpConn *net.TCPConn) {
	// close in the end
	defer tcpConn.Close()
	start := time.Now()

	// call process func OnClose
	belogs.Debug("OnClose(): tcp server,tcpConn: ", tcpConn.RemoteAddr().String(), "   call process func: OnClose ")
	ts.tcpServerProcessFunc.OnCloseProcess(tcpConn)

	// remove tcpConn from tcpConns
	ts.tcpConnsMutex.Lock()
	defer ts.tcpConnsMutex.Unlock()
	belogs.Debug("OnClose(): tcp server,tcpConn: ", tcpConn.RemoteAddr().String(), "   old len(tcpConns): ", len(ts.tcpConns))
	newTcpConns := make(map[string]*net.TCPConn, len(ts.tcpConns))
	for i := range ts.tcpConns {
		if ts.tcpConns[i] != tcpConn {
			connKey := tcpConn.RemoteAddr().String()
			newTcpConns[connKey] = tcpConn
		}
	}
	ts.tcpConns = newTcpConns
	tcpConn = nil
	belogs.Info("OnClose(): tcp server,new len(tcpConns): ", len(ts.tcpConns), "  time(s):", time.Now().Sub(start).Seconds())
}

// connKey is "": send to all clients
// connKey is net.Conn.Address.String(): send this client
func (ts *TcpServer) ActiveSend(sendData []byte, connKey string) (err error) {
	ts.tcpConnsMutex.RLock()
	defer ts.tcpConnsMutex.RUnlock()
	start := time.Now()

	belogs.Debug("ActiveSend(): tcp server,len(sendData):", len(sendData), "   len(tcpConns): ", len(ts.tcpConns), "  connKey:", connKey)
	if len(connKey) == 0 {
		belogs.Debug("ActiveSend(): tcp server, to all, len(sendData):", len(sendData), "   len(tcpConns): ", len(ts.tcpConns))
		for i := range ts.tcpConns {
			belogs.Debug("ActiveSend(): tcp   to all, client: ", i, "    ts.tcpConns[i]:", ts.tcpConns[i], "   call process func: ActiveSend ")
			err = ts.tcpServerProcessFunc.ActiveSendProcess(ts.tcpConns[i], sendData)
			if err != nil {
				// just logs, not return or break
				belogs.Error("ActiveSend(): tcp  fail, to all, client: ", i, "    ts.tcpConns[i]:", ts.tcpConns[i], err)
			}
		}
		belogs.Info("ActiveSend(): tcp  send to all clients ok,  len(sendData):", len(sendData), "   len(tcpConns): ", len(ts.tcpConns),
			"  time(s):", time.Now().Sub(start).Seconds())
		return
	} else {
		belogs.Debug("ActiveSend(): tcp server, to connKey:", connKey)
		if tcpConn, ok := ts.tcpConns[connKey]; ok {
			err = ts.tcpServerProcessFunc.ActiveSendProcess(tcpConn, sendData)
			if err != nil {
				// just logs, not return or break
				belogs.Error("ActiveSend(): tcp  fail, to connKey: ", connKey, "   tcpConn:", tcpConn.RemoteAddr().String(), err)
			}
		}
		belogs.Info("ActiveSend(): tcp  send to connKey ok,  len(sendData):", len(sendData), "   connKey: ", connKey,
			"  time(s):", time.Now().Sub(start).Seconds())
		return
	}

}

// gracefulClose: close gracefully or forcibly abort
// connKey is "": send to all clients
// connKey is net.Conn.Address.String(): send this client
func (ts *TcpServer) ActiveClose(gracefulClose bool, connKey string) {

}
