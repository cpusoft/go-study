package main

import (
	"bytes"
	"crypto/tls"
	"net"
	"testing"

	"github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/convert"
	"github.com/onsi/gomega/gstruct/errors"
)

const (
	PDU_TYPE_MIN_LEN = 8
)

type ServerProcessFunc struct {
}

func (spf *ServerProcessFunc) OnConnectProcess(tlsConn *tls.Conn) (err error) {
	peerCerts := tlsConn.ConnectionState().PeerCertificates
	if peerCerts == nil || len(peerCerts) == 0 {
		return errors.New("perrCerts is emtpy")
	}
	// The first element is the leaf certificate that the connection is verified against
	clientCert := peerCerts[0]

	subject := clientCert.Subject.CommonName
	belogs.Debug("OnConnectProcess(): spf: subject:", subject)
	dnsNames := clientCert.DNSNames
	belogs.Info("OnConnectProcess(): spf: dnsNames:", dnsNames)

	// can active send msg to client
	return nil
}
func (spf *ServerProcessFunc) ReceiveAndSendProcess(conn net.Conn, receiveData []byte) (leftData []byte, err error) {
	belogs.Info("ReceiveAndSendProcess(): receiveData:", convert.Bytes2String(receiveData))

	for {
		// check
		// unpack: TCP sticky packet
		if len(receiveData) < PDU_TYPE_MIN_LEN {
			leftData = make([]byte, len(receiveData))
			copy(leftData, receiveData)
			return leftData, nil
		}

		// RTR: length : byte[4:8]
		lengthBuffer := receiveData[4:8]
		length := convert.Bytes2Uint64(lengthBuffer)
		belogs.Info("ReceiveAndSendProcess(): length:", length)

		if length < PDU_TYPE_MIN_LEN {
			leftData = make([]byte, len(receiveData))
			copy(leftData, receiveData)
			return nil, errors.New("length is err")
		}
		if length > len(receiveData) {
			leftData = make([]byte, len(receiveData))
			copy(leftData, receiveData)
			return leftData, nil
		} else if length == len(receiveData) {
			rtrData := receiveData
			err = RtrProcess(rtrData)
			if err != nil {
				return nil, err
			}
			return make([]byte, 0), nil
		} else if length < len(receiveData) {
			rtrData := receiveData[:length]

			err = RtrProcess(rtrData)
			if err != nil {
				return nil, err
			}
			// left may have another rtr data
			leftData = make([]byte, length)
			copy(leftData, receiveData[length:])
			receiveData = leftData
		}

	}
	return leftData, nil
}
func (spf *ServerProcessFunc) OnCloseProcess(conn net.Conn) {

}
func (spf *ServerProcessFunc) ActiveSendProcess(conn net.Conn, sendData []byte) (err error) {
	return
}

func TestCreateTlsServer(t *testing.T) {
	serverProcessFunc := new(ServerProcessFunc)
	rootCrtFileName := `\go-study\data\ca\ca.cer`
	publicCrtFileName := `\go-study\data\server\server.cer`
	privateKeyFileName := `\go-study\data\server\serverkey.pem`

	ts := NewTlsServer(serverProcessFunc, rootCrtFileName, publicCrtFileName, privateKeyFileName, true)
	ts.Start("0.0.0.0:9999")
}

func RtrProcess(receiveData []byte) (err error) {
	buf := bytes.NewReader(receiveData)
	belogs.Info("RtrProcess(): buf:", buf)
	return nil
}
