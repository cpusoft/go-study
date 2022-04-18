package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"sync"
	"time"

	belogs "github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/jsonutil"
)

//https://github.com/SmarkSeven/socket/blob/master/protocol.go
//https://zhuanlan.zhihu.com/p/338688506

// https://stackoverflow.com/questions/13110713/upgrade-a-connection-to-tls-in-go
// https://groups.google.com/g/golang-nuts/c/x8xOwJ6i8vg ok
type TlsServer struct {
	tlsConns      map[string]*tls.Conn // map[addr]*net.TCPConn
	tlsConnsMutex sync.RWMutex

	rootCrtFileName    string
	publicCrtFileName  string
	privateKeyFileName string
	verifyClient       bool

	tlsServerProcessFunc TlsServerProcessFunc
}

type TlsServerProcessFunc interface {
	OnConnectProcess(tlsConn *tls.Conn) (err error)
	ReceiveAndSendProcess(tlsConn *tls.Conn, receiveData []byte) (leftData []byte, err error)
	OnCloseProcess(tlsConn *tls.Conn)
	ActiveSendProcess(tlsConn *tls.Conn, sendData []byte) (err error)
}

func NewTlsServer(rootCrtFileName, publicCrtFileName, privateKeyFileName string, verifyClient bool,
	tlsServerProcessFunc TlsServerProcessFunc) (ts *TlsServer) {

	belogs.Debug("NewTcpServer():tlsServerProcessFunc:", tlsServerProcessFunc)
	ts = &TlsServer{}
	ts.tlsConns = make(map[string]*tls.Conn, 16)
	ts.rootCrtFileName = rootCrtFileName
	ts.publicCrtFileName = publicCrtFileName
	ts.privateKeyFileName = privateKeyFileName
	ts.tlsServerProcessFunc = tlsServerProcessFunc
	ts.verifyClient = verifyClient
	belogs.Debug("NewTlsServer():ts:", ts, "   ts.tlsConnsMutex:", ts.tlsConnsMutex)
	return ts
}

// port `:8888`
func (ts *TlsServer) Start(port string) (err error) {

	cert, err := tls.LoadX509KeyPair(ts.publicCrtFileName, ts.privateKeyFileName)
	if err != nil {
		belogs.Error("Start(): tls  LoadX509KeyPair fail: port:", port,
			"  publicCrtFileName, privateKeyFileName:", ts.publicCrtFileName, ts.privateKeyFileName, err)
		return err
	}
	rootCrtBytes, err := ioutil.ReadFile(ts.rootCrtFileName)
	if err != nil {
		belogs.Error("Start(): tls  ReadFile rootCrtFileName fail, port:", port,
			"  rootCrtFileName:", ts.rootCrtFileName, err)
		return err
	}
	rootCertPool := x509.NewCertPool()
	ok := rootCertPool.AppendCertsFromPEM(rootCrtBytes)
	if !ok {
		belogs.Error("Start(): tls  AppendCertsFromPEM rootCrtFileName fail,port:", port,
			"  rootCrtFileName:", ts.rootCrtFileName, "  len(rootCrtBytes):", len(rootCrtBytes), err)
		return err
	}
	clientAuthType := tls.NoClientCert
	if ts.verifyClient {
		clientAuthType = tls.RequireAndVerifyClientCert
	}

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         clientAuthType,
		RootCAs:            rootCertPool,
		InsecureSkipVerify: false,
	}
	listen, err := tls.Listen("tcp", port, config)
	if err != nil {
		belogs.Error("Start(): tls  Listen fail: ", port, err)
		return err
	}
	defer listen.Close()
	belogs.Debug("Start(): tls  create server ok, server is ", port, "  will accept client")
	//var tlsConn tls.Conn
	for {
		conn, err := listen.Accept()
		belogs.Info("Start(): tls  Accept remote: ", conn.RemoteAddr().String())
		if err != nil {
			belogs.Error("Start(): tls  Accept remote fail: ", port, conn.RemoteAddr().String(), err)
			continue
		}
		if conn == nil {
			continue
		}

		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			belogs.Error("Start(): tls  conn cannot conver to tls: ", port, conn.RemoteAddr().String(), err)
			continue
		}
		ts.OnConnect(tlsConn)
		// call func to process tlsConn
		go ts.ReceiveAndSend(tlsConn)

	}
	return nil
}

func (ts *TlsServer) OnConnect(tlsConn *tls.Conn) {
	start := time.Now()
	belogs.Debug("OnConnect(): tls   new tlsConn: ", tlsConn)

	pcs := tlsConn.ConnectionState().PeerCertificates
	belogs.Debug("OnConnect(): tls perr certs: ", jsonutil.MarshalJson(pcs))

	// add new tlsConn to tcpconns
	ts.tlsConnsMutex.Lock()
	defer ts.tlsConnsMutex.Unlock()
	connKey := tlsConn.RemoteAddr().String()
	ts.tlsConns[connKey] = tlsConn
	belogs.Debug("OnConnect(): tls tlsConn: ", tlsConn.RemoteAddr().String(), ", new len(tlsConns): ", len(ts.tlsConns))

	// call process func OnConnect
	belogs.Debug("OnConnect(): tls tlsConn: ", tlsConn.RemoteAddr().String(), "   call process func: OnConnect ")
	err := ts.tlsServerProcessFunc.OnConnectProcess(tlsConn)
	if err != nil {
		belogs.Error("OnConnect(): tls tlsServerProcessFunc.OnConnect fail, will Close: ", tlsConn.RemoteAddr().String(), err)
		ts.OnClose(tlsConn)
		return
	}
	belogs.Info("OnConnect(): tls add tlsConn: ", tlsConn.RemoteAddr().String(), "   len(tlsConns): ", len(ts.tlsConns), "   time(s):", time.Now().Sub(start).Seconds())
}

func (ts *TlsServer) ReceiveAndSend(tlsConn *tls.Conn) {

	defer ts.OnClose(tlsConn)

	var leftData []byte
	// one packet
	buffer := make([]byte, 2048)
	// wait for new packet to read
	var err error
	for {
		n, err := tlsConn.Read(buffer)
		start := time.Now()
		belogs.Debug("ReceiveAndSend(): tls server read: Read n: ", tlsConn.RemoteAddr().String(), n)
		if err != nil {
			if err == io.EOF {
				// is not error, just client close
				belogs.Debug("ReceiveAndSend(): tls server Read io.EOF, client close: ", tlsConn.RemoteAddr().String(), err)
				return
			}
			belogs.Error("ReceiveAndSend(): tls server Read fail, err ", tlsConn.RemoteAddr().String(), err)
			return
		}
		if n == 0 {
			continue
		}

		// call process func OnReceiveAndSend
		// copy to leftData
		belogs.Debug("ReceiveAndSend(): tls server tlsConn: ", tlsConn.RemoteAddr().String(), "  n:", n,
			" , will call process func: OnReceiveAndSend,  time(s):", time.Now().Sub(start))
		leftData, err = ts.tlsServerProcessFunc.ReceiveAndSendProcess(tlsConn, append(leftData, buffer[:n]...))
		belogs.Debug("ReceiveAndSend(): tls server tlsConn: ", tlsConn.RemoteAddr().String(), " receive n: ", n,
			"  len(leftData):", len(leftData), "  time(s):", time.Now().Sub(start))
		if err != nil {
			belogs.Error("OnReceiveAndSend():server fail ,will remove this tlsConn : ", tlsConn.RemoteAddr().String(), err)
			break
		}
	}
	belogs.Info("OnReceiveAndSend():break for, will remove this tlsConn: ", tlsConn.RemoteAddr().String(), "  err:", err)
}

func (ts *TlsServer) OnClose(tlsConn *tls.Conn) {
	// close in the end
	defer tlsConn.Close()
	start := time.Now()

	// call process func OnClose
	belogs.Debug("OnClose(): tls server,tlsConn: ", tlsConn.RemoteAddr().String(), "   call process func: OnClose ")
	ts.tlsServerProcessFunc.OnCloseProcess(tlsConn)

	// remove tlsConn from tlsConns
	ts.tlsConnsMutex.Lock()
	defer ts.tlsConnsMutex.Unlock()
	belogs.Debug("OnClose(): tls server,tlsConn: ", tlsConn.RemoteAddr().String(), "   old len(tlsConns): ", len(ts.tlsConns))
	newTlsConns := make(map[string]*tls.Conn, len(ts.tlsConns))
	for i := range ts.tlsConns {
		if ts.tlsConns[i] != tlsConn {
			connKey := tlsConn.RemoteAddr().String()
			newTlsConns[connKey] = tlsConn
		}
	}
	ts.tlsConns = newTlsConns
	belogs.Info("OnClose(): tls server,new len(tlsConns): ", len(ts.tlsConns), "  time(s):", time.Now().Sub(start).Seconds())
}

func (ts *TlsServer) ActiveSend(sendData []byte, connKey string) (err error) {
	ts.tlsConnsMutex.RLock()
	defer ts.tlsConnsMutex.RUnlock()
	start := time.Now()

	belogs.Debug("ActiveSend(): tls server,len(sendData):", len(sendData), "   len(tlsConns): ", len(ts.tlsConns), "  connKey:", connKey)
	if len(connKey) == 0 {
		belogs.Debug("ActiveSend(): tls server, to all, len(sendData):", len(sendData), "   len(tlsConns): ", len(ts.tlsConns))
		for i := range ts.tlsConns {
			belogs.Debug("ActiveSend(): tls   to all, client: ", i, "    ts.tlsConns[i]:", ts.tlsConns[i], "   call process func: ActiveSend ")
			err = ts.tlsServerProcessFunc.ActiveSendProcess(ts.tlsConns[i], sendData)
			if err != nil {
				// just logs, not return or break
				belogs.Error("ActiveSend(): tls  fail, to all, client: ", i, "    ts.tlsConns[i]:", ts.tlsConns[i], err)
			}
		}
		belogs.Info("ActiveSend(): tls  send to all clients ok,  len(sendData):", len(sendData), "   len(tlsConns): ", len(ts.tlsConns),
			"  time(s):", time.Now().Sub(start).Seconds())
		return
	} else {
		belogs.Debug("ActiveSend(): tls server, to connKey:", connKey)
		if tcpConn, ok := ts.tlsConns[connKey]; ok {
			err = ts.tlsServerProcessFunc.ActiveSendProcess(tcpConn, sendData)
			if err != nil {
				// just logs, not return or break
				belogs.Error("ActiveSend(): tls  fail, to connKey: ", connKey, "   tcpConn:", tcpConn.RemoteAddr().String(), err)
			}
		}
		belogs.Info("ActiveSend(): tls  send to connKey ok,  len(sendData):", len(sendData), "   connKey: ", connKey,
			"  time(s):", time.Now().Sub(start).Seconds())
		return
	}
	return
}
