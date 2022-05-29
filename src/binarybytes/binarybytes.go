package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cpusoft/goutil/convert"
)

func main() {
	b := []byte{0x01, 0x02, 0x03, 0x04}
	c := NewDsoEncryptionPaddingTlvModel(b)
	fmt.Println(c.PrintBytes())
	var q uint16
	q = (1 << 15)
	fmt.Printf("%0x\n", q)
	q = 6 << 11
	fmt.Printf("%0x\n", q)
	q |= uint16(2)
	fmt.Printf("%0x\n", q)
}

const (
	DSO_TYPE_RESERVED           = 0
	DSO_TYPE_KEEP_ALIVE         = 1
	DSO_TYPE_RETRY_DELAY        = 2
	DSO_TYPE_ENCRYPTION_PADDING = 3
)

type DsoModel struct {
	MessageId uint16 `json:"messageId"`

	QrOpCodeZRCodeQdCount uint16 `json:"qrOpCodeZRCodeQdCount"`
	Qr                    uint8  `json:"qr"`
	OpCode                uint8  `json:"opCode"`
	Z                     uint8  `json:"z"`
	RCode                 uint8  `json:"rCode"`

	QdCount uint16 `json:"qdCount"`
	AnCount uint16 `json:"anCount"`
	NsCount uint16 `json:"nsCount"`
	ArCount uint16 `json:"arCount"`

	DsoTlvModels []DsoTlvModel `json:"dsoTlvModels"`
}

func (c *DsoModel) Bytes() []byte {
	wr := bytes.NewBuffer([]byte{})
	binary.Write(wr, binary.BigEndian, c.MessageId)
	binary.Write(wr, binary.BigEndian, c.QrOpCodeZRCodeQdCount)
	binary.Write(wr, binary.BigEndian, c.QdCount)
	binary.Write(wr, binary.BigEndian, c.NsCount)
	binary.Write(wr, binary.BigEndian, c.ArCount)
	for i := range c.DsoTlvModels {
		binary.Write(wr, binary.BigEndian, c.DsoTlvModels[i].Bytes())
	}
	return wr.Bytes()
}
func (c *DsoModel) PrintBytes() string {
	return convert.PrintBytes(c.Bytes(), 8)
}

type DsoTlvModel interface {
	Bytes() []byte
	PrintBytes() string
}

type DsoKeepaliveTlvModel struct {
	DsoType           uint16 `json:"dsoType"`
	DsoLength         uint16 `json:"dsoLength"`
	InactivityTimeout uint16 `json:"inactivityTimeout"`
	KeepaliveInterval uint16 `json:"keepaliveInterval"`
}

func NewDsoKeepaliveTlvModel(inactivityTimeout, keepaliveInterval uint16) *DsoKeepaliveTlvModel {
	c := &DsoKeepaliveTlvModel{
		DsoType:           DSO_TYPE_KEEP_ALIVE,
		DsoLength:         2,
		InactivityTimeout: inactivityTimeout,
		KeepaliveInterval: keepaliveInterval,
	}
	return c
}
func (c *DsoKeepaliveTlvModel) Bytes() []byte {
	wr := bytes.NewBuffer([]byte{})
	binary.Write(wr, binary.BigEndian, c.DsoType)
	binary.Write(wr, binary.BigEndian, c.DsoLength)
	binary.Write(wr, binary.BigEndian, c.InactivityTimeout)
	binary.Write(wr, binary.BigEndian, c.KeepaliveInterval)
	return wr.Bytes()
}
func (c *DsoKeepaliveTlvModel) PrintBytes() string {
	return convert.PrintBytes(c.Bytes(), 8)
}

type DsoRetryDelayTlvModel struct {
	DsoType    uint16 `json:"dsoType"`
	DsoLength  uint16 `json:"dsoLength"`
	RetryDelay uint16 `json:"retryDelay"`
}

func (c *DsoRetryDelayTlvModel) Bytes() []byte {
	wr := bytes.NewBuffer([]byte{})
	binary.Write(wr, binary.BigEndian, c.DsoType)
	binary.Write(wr, binary.BigEndian, c.DsoLength)
	binary.Write(wr, binary.BigEndian, c.RetryDelay)
	return wr.Bytes()
}
func (c *DsoRetryDelayTlvModel) PrintBytes() string {
	return convert.PrintBytes(c.Bytes(), 8)
}

type DsoEncryptionPaddingTlvModel struct {
	DsoType           uint16 `json:"dsoType"`
	DsoLength         uint16 `json:"dsoLength"`
	EncryptionPadding []byte `json:"encryptionPadding"`
}

func NewDsoEncryptionPaddingTlvModel(encryptionPadding []byte) *DsoEncryptionPaddingTlvModel {
	c := &DsoEncryptionPaddingTlvModel{
		DsoType:           DSO_TYPE_ENCRYPTION_PADDING,
		DsoLength:         uint16(len(encryptionPadding)),
		EncryptionPadding: encryptionPadding,
	}
	return c
}
func (c *DsoEncryptionPaddingTlvModel) Bytes() []byte {
	wr := bytes.NewBuffer([]byte{})
	binary.Write(wr, binary.BigEndian, c.DsoType)
	binary.Write(wr, binary.BigEndian, c.DsoLength)
	binary.Write(wr, binary.BigEndian, c.EncryptionPadding)
	return wr.Bytes()

}
func (c *DsoEncryptionPaddingTlvModel) PrintBytes() string {
	return convert.PrintBytes(c.Bytes(), 8)
}
