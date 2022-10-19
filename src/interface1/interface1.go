package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type DnsConnect interface {
	Close()
	SendTcpTlsMsg(msg string)
}

type UndefinedConnect struct {
	IsOpen bool
	// close outside
	Msg string
}

func NewUndefinedConnect() DnsConnect {
	c := &UndefinedConnect{}
	c.IsOpen = true
	c.Msg = "undefined"
	return c
}

func (c UndefinedConnect) SendTcpTlsMsg(msg string) {
	c.Msg = msg
}
func (c UndefinedConnect) Close() {
	c.IsOpen = false
}

func main() {
	dnsConnect := NewUndefinedConnect()
	fmt.Println(jsonutil.MarshalJson(dnsConnect))

	undefinedConnect, ok := (dnsConnect).(*UndefinedConnect)
	fmt.Println(jsonutil.MarshalJson(undefinedConnect), ok)
}
