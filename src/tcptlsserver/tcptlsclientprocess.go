package main

type TlsClientProcessFunc interface {
	OnConnectProcess(tcpTlsConn *TcpTlsConn)
	OnCloseProcess(tcpTlsConn *TcpTlsConn)
	OnReceiveProcess(tcpTlsConn *TcpTlsConn, sendData []byte) (nextRwPolicy int, leftData []byte, err error)
}
