package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net"
	"time"

	belogs "github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/jsonutil"
	"github.com/cpusoft/goutil/tcpserverclient/util"
)

type TcpTlsClientSendMsg struct {
	NextConnectClosePolicy int //NEXT_CONNECT_CLOSE_POLICE_NO  NEXT_CONNECT_CLOSE_POLICE_GRACEFUL  NEXT_CONNECT_CLOSE_POLICE_FORCIBLE
	NextRwPolice           int //NEXT_RW_POLICE_ALL,NEXT_RW_POLICE_WAIT_READ,NEXT_RW_POLICE_WAIT_WRITE
	SendData               []byte
}

type TcpTlsClient struct {
	// both tcp and tls
	isTcpClient             bool
	tcpTlsClientSendMsg     chan TcpTlsClientSendMsg
	tcpTlsClientProcessFunc TcpTlsClientProcessFunc

	// for tls
	tlsRootCrtFileName    string
	tlsPublicCrtFileName  string
	tlsPrivateKeyFileName string
}

// server: 0.0.0.0:port
func NewTcpClient(tcpTlsClientProcessFunc TcpTlsClientProcessFunc) (tc *TcpTlsClient) {

	belogs.Debug("NewTcpClient():tcpTlsClientProcessFunc:", tcpTlsClientProcessFunc)
	tc = &TcpTlsClient{}
	tc.isTcpClient = true
	tc.tcpTlsClientSendMsg = make(chan TcpTlsClientSendMsg)
	tc.tcpTlsClientProcessFunc = tcpTlsClientProcessFunc
	belogs.Info("NewTcpClient():tc:", tc)
	return tc
}

// server: 0.0.0.0:port
func NewTlsClient(tlsRootCrtFileName, tlsPublicCrtFileName, tlsPrivateKeyFileName string,
	tcpTlsClientProcessFunc TcpTlsClientProcessFunc) (tc *TcpTlsClient) {

	belogs.Debug("NewTlsClient():tcpTlsClientProcessFunc:", tcpTlsClientProcessFunc)
	tc = &TcpTlsClient{}
	tc.isTcpClient = false
	tc.tcpTlsClientSendMsg = make(chan TcpTlsClientSendMsg)
	tc.tcpTlsClientProcessFunc = tcpTlsClientProcessFunc

	tc.tlsRootCrtFileName = tlsRootCrtFileName
	tc.tlsPublicCrtFileName = tlsPublicCrtFileName
	tc.tlsPrivateKeyFileName = tlsPrivateKeyFileName

	belogs.Info("NewTlsClient():tc:", tc)
	return tc
}

// server: **.**.**.**:port
func (tc *TcpTlsClient) StartTcpClient(server string) (err error) {
	belogs.Debug("StartTcpClient(): create client, server is  ", server)

	tcpServer, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		belogs.Error("StartTcpClient():  ResolveTCPAddr fail: ", server, err)
		return err
	}
	belogs.Debug("StartTcpClient(): create client, server is  ", server, "  tcpServer:", tcpServer)

	tcpConn, err := net.DialTCP("tcp4", nil, tcpServer)
	if err != nil {
		belogs.Error("StartTcpClient():  Dial fail, server:", server, "  tcpServer:", tcpServer, err)
		return err
	}
	tcpTlsConn := NewFromTcpConn(tcpConn)
	tc.OnConnect(tcpTlsConn)
	belogs.Info("StartTcpClient(): OnConnect, server is  ", server, "  tcpTlsConn:", tcpTlsConn.RemoteAddr().String())

	//active send to server, and receive from server, loop
	go tc.SendAndReceive(tcpTlsConn)
	belogs.Debug("StartTcpClient(): SendAndReceive, server:", server, "   tcpTlsConn:", tcpTlsConn.RemoteAddr().String())
	return nil
}

// server: **.**.**.**:port
func (tc *TcpTlsClient) StartTlsClient(server string) (err error) {
	belogs.Debug("StartTlsClient(): create client, server is  ", server)

	cert, err := tls.LoadX509KeyPair(tc.tlsPublicCrtFileName, tc.tlsPrivateKeyFileName)
	if err != nil {
		belogs.Error("StartTlsClient(): LoadX509KeyPair fail: server:", server,
			"  tlsPublicCrtFileName, tlsPrivateKeyFileName:", tc.tlsPublicCrtFileName, tc.tlsPrivateKeyFileName, err)
		return err
	}
	rootCrtBytes, err := ioutil.ReadFile(tc.tlsRootCrtFileName)
	if err != nil {
		belogs.Error("StartTlsClient(): ReadFile tlsRootCrtFileName fail, server:", server,
			"  tlsRootCrtFileName:", tc.tlsRootCrtFileName, err)
		return err
	}
	rootCertPool := x509.NewCertPool()
	ok := rootCertPool.AppendCertsFromPEM(rootCrtBytes)
	if !ok {
		belogs.Error("StartTlsClient(): AppendCertsFromPEM tlsRootCrtFileName fail,server:", server,
			"  tlsRootCrtFileName:", tc.tlsRootCrtFileName, "  len(rootCrtBytes):", len(rootCrtBytes), err)
		return err
	}
	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            rootCertPool,
		InsecureSkipVerify: false,
	}

	tlsConn, err := tls.Dial("tcp", server, config)
	if err != nil {
		belogs.Error("StartTlsClient(): Dial fail, server:", server, err)
		return err
	}
	tcpTlsConn := NewFromTlsConn(tlsConn)
	tc.OnConnect(tcpTlsConn)
	belogs.Info("StartTlsClient(): OnConnect, server is  ", server, "  tcpTlsConn:", tcpTlsConn.RemoteAddr().String())

	//active send to server, and receive from server, loop
	go tc.SendAndReceive(tcpTlsConn)
	belogs.Debug("StartTlsClient(): SendAndReceive, server:", server, "   tcpTlsConn:", tcpTlsConn.RemoteAddr().String())
	return nil
}

func (tc *TcpTlsClient) OnConnect(tcpTlsConn *TcpTlsConn) {
	// call process func OnConnect
	tc.tcpTlsClientProcessFunc.OnConnectProcess(tcpTlsConn)
	belogs.Info("OnConnect(): tcptlsclient  after OnConnectProcess, tcpTlsConn: ", tcpTlsConn.RemoteAddr().String())
}

func (tc *TcpTlsClient) OnClose(tcpTlsConn *TcpTlsConn) {
	// close in the end
	belogs.Info("OnClose(): tcptlsclient , tcpTlsConn: ", tcpTlsConn.RemoteAddr().String())
	tcpTlsConn.Close()
	tcpTlsConn.SetNil()
}

func (tc *TcpTlsClient) SendMsg(tcpTlsClientSendMsg *TcpTlsClientSendMsg) {

	belogs.Debug("SendMsg(): tcptlsclient, tcpTlsClientSendMsg:", jsonutil.MarshalJson(*tcpTlsClientSendMsg))
	tc.tcpTlsClientSendMsg <- *tcpTlsClientSendMsg
}

func (tc *TcpTlsClient) SendAndReceive(tcpTlsConn *TcpTlsConn) (err error) {
	belogs.Debug("SendAndReceive(): tcptlsclient , tcpTlsConn:", tcpTlsConn.RemoteAddr().String())
	for {
		// wait next tcpTlsClientSendMsg: only error or NEXT_CONNECT_POLICE_CLOSE_** will end loop
		select {
		case tcpTlsClientSendMsg := <-tc.tcpTlsClientSendMsg:
			nextConnectClosePolicy := tcpTlsClientSendMsg.NextConnectClosePolicy
			nextRwPolice := tcpTlsClientSendMsg.NextRwPolice
			sendData := tcpTlsClientSendMsg.SendData
			belogs.Debug("SendAndReceive(): tcptlsclient , tcpTlsConn:", tcpTlsConn.RemoteAddr().String(),
				"  tcpTlsClientSendMsg: ", jsonutil.MarshalJson(tcpTlsClientSendMsg))

			// if close
			if nextConnectClosePolicy == util.NEXT_CONNECT_POLICE_CLOSE_GRACEFUL ||
				nextConnectClosePolicy == util.NEXT_CONNECT_POLICE_CLOSE_FORCIBLE {
				belogs.Info("SendAndReceive(): tcptlsclient   nextConnectClosePolicy close end client, will end tcpTlsConn: ", tcpTlsConn.RemoteAddr().String(),
					"   nextConnectClosePolicy:", nextConnectClosePolicy)
				tc.OnClose(tcpTlsConn)
				return nil
			}

			// send data
			start := time.Now()
			n, err := tcpTlsConn.Write(sendData)
			if err != nil {
				belogs.Error("SendAndReceive(): tcptlsclient   Write fail:  tcpTlsConn:", tcpTlsConn.RemoteAddr().String(), err)
				return err
			}
			belogs.Debug("SendAndReceive(): tcptlsclient   Write to tcpTlsConn:", tcpTlsConn.RemoteAddr().String(),
				"  len(sendData):", len(sendData), "  write n:", n, "   nextRwPolice:", nextRwPolice,
				"  time(s):", time.Now().Sub(start).Seconds())

			// if wait receive, then wait next tcpTlsClientSendMsg
			if nextRwPolice == util.NEXT_RW_POLICE_WAIT_READ {
				// if server tell client: end this loop, or end conn
				err := tc.OnReceive(tcpTlsConn)
				if err != nil {
					belogs.Error("SendAndReceive(): tcptlsclient   Write fail:  tcpTlsConn:", tcpTlsConn.RemoteAddr().String(), err)
					return err
				}
				belogs.Info("SendAndReceive(): tcptlsclient  shouldWaitReceive yes, tcpTlsConn:", tcpTlsConn.RemoteAddr().String(),
					"  len(sendData):", len(sendData), "  write n:", n,
					"  time(s):", time.Now().Sub(start).Seconds())
				continue
			} else {
				belogs.Info("SendAndReceive(): tcptlsclient  OnReceive, shouldWaitReceive no, will return: ", tcpTlsConn.RemoteAddr().String())
				continue
			}
		}
	}

}

func (tc *TcpTlsClient) OnReceive(tcpTlsConn *TcpTlsConn) (err error) {
	belogs.Debug("OnReceive(): tcptlsclient  wait for OnReceive, tcpTlsConn:", tcpTlsConn.RemoteAddr().String())
	var leftData []byte
	// one packet
	buffer := make([]byte, 2048)
	// wait for new packet to read

	for {
		n, err := tcpTlsConn.Read(buffer)
		start := time.Now()
		belogs.Debug("OnReceive(): tcptlsclient  client read: Read n: ", tcpTlsConn.RemoteAddr().String(), n)
		if err != nil {
			if err == io.EOF {
				// is not error, just client close
				belogs.Debug("OnReceive(): tcptlsclient   io.EOF, client close: ", tcpTlsConn.RemoteAddr().String(), err)
				return nil
			}
			belogs.Error("OnReceive(): tcptlsclient   Read fail, err ", tcpTlsConn.RemoteAddr().String(), err)
			return err
		}
		if n == 0 {
			continue
		}

		belogs.Debug("OnReceive(): tcptlsclient  client tcpTlsConn: ", tcpTlsConn.RemoteAddr().String(), "  n:", n,
			" , will call process func: OnReceiveAndSend,  time(s):", time.Now().Sub(start))
		nextRwPolicy, leftData, err := tc.tcpTlsClientProcessFunc.OnReceiveProcess(tcpTlsConn, append(leftData, buffer[:n]...))
		belogs.Info("OnReceive(): tcptlsclient  tcpTlsClientProcessFunc.OnReceiveProcess, tcpTlsConn: ", tcpTlsConn.RemoteAddr().String(), " receive n: ", n,
			"  len(leftData):", len(leftData), "  nextRwPolicy:", nextRwPolicy, "  time(s):", time.Now().Sub(start))
		if err != nil {
			belogs.Error("OnReceive(): tcptlsclient  tcpTlsClientProcessFunc.OnReceiveProcess  fail ,will close this tcpTlsConn : ", tcpTlsConn.RemoteAddr().String(), err)
			return err
		}
		if nextRwPolicy == util.NEXT_RW_POLICE_END_READ {
			belogs.Debug("OnReceive(): tcptlsclient  nextRwPolicy, will end this write/read loop: ", tcpTlsConn.RemoteAddr().String())
			return nil
		}
	}

}
