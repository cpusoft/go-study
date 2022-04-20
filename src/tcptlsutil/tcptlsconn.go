package main

import (
	"crypto/tls"
	"errors"
	"net"
)

type TcpTlsConn struct {
	// tcp: true
	// tls: false
	isTcpConn bool
	tcpConn   *net.TCPConn
	tlsConn   *tls.Conn

	nextConnectPolicy int
}

func NewFromTcpConn(tcpConn *net.TCPConn) (c *TcpTlsConn) {
	c = &TcpTlsConn{}
	c.tcpConn = tcpConn
	c.isTcpConn = true
	c.nextConnectPolicy = NEXT_CONNECT_POLICE_KEEP
	return c
}

func NewFromTlsConn(tlsConn *tls.Conn) (c *TcpTlsConn) {
	c = &TcpTlsConn{}
	c.tlsConn = tlsConn
	c.isTcpConn = false
	c.nextConnectPolicy = NEXT_CONNECT_POLICE_KEEP
	return c
}

func (c *TcpTlsConn) RemoteAddr() net.Addr {
	if c.isTcpConn && c.tcpConn != nil {
		return c.tcpConn.RemoteAddr()
	} else if c.tlsConn != nil {
		return c.tlsConn.RemoteAddr()
	}
	return nil
}

func (c *TcpTlsConn) Write(b []byte) (n int, err error) {
	if c.isTcpConn && c.tcpConn != nil {
		return c.tcpConn.Write(b)
	} else if c.tlsConn != nil {
		return c.tlsConn.Write(b)
	}
	return -1, errors.New("is not conn")
}

func (c *TcpTlsConn) Read(b []byte) (n int, err error) {
	if c.isTcpConn && c.tcpConn != nil {
		return c.tcpConn.Read(b)
	} else if c.tlsConn != nil {
		return c.tlsConn.Read(b)
	}
	return -1, errors.New("is not conn")
}

func (c *TcpTlsConn) Close() (err error) {
	if c.isTcpConn && c.tcpConn != nil {
		return c.tcpConn.Close()
	} else if c.tlsConn != nil {
		return c.tlsConn.Close()
	}
	return errors.New("is not conn")
}

func (c *TcpTlsConn) SetNil() {
	if c.isTcpConn && c.tcpConn != nil {
		c.tcpConn = nil
	} else if c.tlsConn != nil {
		c.tlsConn = nil
	}
}
