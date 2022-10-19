package main

import "fmt"

type dnsError struct {
	msg string // from error.Error()

	id                uint16 // from dns id/messageId
	opCode            uint8
	rCode             uint8 // response DSO_RCODE_***
	nextConnectPolicy int   //	tcptlsutil.NEXT_CONNECT_POLICY_***

}

func (c dnsError) Error() string {
	//return fmt.Sprintf(`{"code":%d,"msg":"%v"}`, e.code, e.msg)
	return fmt.Sprintf(`{"msg":"%s","id":%d,"opCode":%d,"rCode":%d,"nextConnectPolicy":%d}`,
		c.msg, c.id, c.opCode, c.rCode, c.nextConnectPolicy)
}

func NewDnsError(msg string, id uint16, opCode uint8, rCode uint8, nextConnectPolicy int) error {
	return dnsError{
		msg:               msg,
		id:                id,
		opCode:            opCode,
		rCode:             rCode,
		nextConnectPolicy: nextConnectPolicy,
	}
}

func GetMsg(err error) string {
	if e, ok := err.(dnsError); ok {
		return e.msg
	}
	return ""
}

func main() {
	e := NewMyErr("ss", 0, 0, 0, 0)
	exx, ok := e.(dnsError)

	fmt.Println(e)
	fmt.Println(exx, ok)
}
