package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"sync/atomic"
	"text/template"
	"time"
	"unicode"

	"github.com/cpusoft/goutil/jsonutil"
)

type Person struct {
	Name string
	Age  int
}

func personSum(ps []Person, p chan Person) {
	one := Person{"", 0}
	for _, v := range ps {
		one.Age += v.Age
		one.Name += (v.Name + ";")
	}
	p <- one
	fmt.Printf("%+v    ", p)
}
func fib(n int, c chan int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}
func randSeq(n int) string {
	//letters := []rune("abcdefghijklmnopqrstuvwxyz")
	letters := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fi(c, quit chan int, expire chan bool) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			fmt.Println("before x:", x, "  y:", y)
			x, y = y, x+y
			fmt.Println("after x:", x, "  y:", y)
			fmt.Println(c)
		case <-time.After(5 * time.Second):
			fmt.Println("timeout")
			expire <- true
			break

		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
func goRun(c, quit chan int) {
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	quit <- 0
}

type RtrPdu interface {
	String() string
	Bytes() []byte
}

type RtrIpv4Prefix struct {
	ProtocolVersion uint8  `json:"protocolVersion"`
	PduType         uint8  `json:"pduType"`
	Zero0           uint16 `json:"zero0"`
	Length          uint32 `json:"length"`
	Flags           uint8  `json:"flags"`
	PrefixLength    uint8  `json:"prefixLength"`
	MaxLength       uint8  `json:"maxLength"`
	Zero1           uint8  `json:"zero1"`
	Ipv4Prefix      uint32 `json:"ipv4Prefix"`
	Asn             uint32 `json:"asn"`
}

func (c *RtrIpv4Prefix) String() string {
	return jsonutil.MarshalJson(*c)
}

func (p *RtrIpv4Prefix) Bytes() []byte {
	wr := bytes.NewBuffer([]byte{})
	binary.Write(wr, binary.BigEndian, p.ProtocolVersion)
	binary.Write(wr, binary.BigEndian, p.PduType)
	binary.Write(wr, binary.BigEndian, p.Zero0)
	binary.Write(wr, binary.BigEndian, p.Length)
	binary.Write(wr, binary.BigEndian, p.Flags)
	binary.Write(wr, binary.BigEndian, p.PrefixLength)
	binary.Write(wr, binary.BigEndian, p.MaxLength)
	binary.Write(wr, binary.BigEndian, p.Zero1)
	binary.Write(wr, binary.BigEndian, p.Ipv4Prefix)
	binary.Write(wr, binary.BigEndian, p.Asn)
	return wr.Bytes()
}

func test() (RtrPdu, error) {
	rt := RtrIpv4Prefix{}
	rt.ProtocolVersion = 1
	rt.PduType = 1
	rt.PrefixLength = 12
	rt.Length = 12
	rt.MaxLength = 13
	rt.Ipv4Prefix = 0xA1A2A3A4
	rt.Asn = 12
	return &rt, nil

}

func main() {
	sss1 := []byte{0x01, 0x02, 0x03, 0x04, 0xaa, 0xb1, 0xc8}
	fmt.Println(hex.Dump(sss1))

	rt, _ := test()
	fmt.Println(rt)
	fmt.Println(rt.Bytes())
	switch rt.(type) {
	case *RtrIpv4Prefix:

		fmt.Println("case   ", rt.String())

	}

	rtStr := jsonutil.MarshalJson(rt)
	fmt.Println(rtStr)
	fmt.Println(rt.Bytes())
	rt2 := RtrIpv4Prefix{}
	jsonutil.UnmarshalJson(rtStr, &rt2)
	fmt.Println(rt2)

	ssss := `10d0c9f4328576d51cc73c042cfc15e9b3d6378`
	sn, err := strconv.ParseUint(ssss, 16, 0)
	fmt.Println(sn, err)

	//  [:xdigit:]
	//reg := regexp.MustCompile(`[:xdigit:]`)
	b, err := regexp.MatchString(`^[0-9a-fA-F]+$`, ssss+"a111")
	fmt.Println(b, err)

	var as int64
	as = 0
	as1 := atomic.AddInt64(&as, 1)
	fmt.Println(as, as1)

	s := "sssss<script/>"
	ss1 := template.HTMLEscapeString(s)
	fmt.Println(ss1)

	ss1 = template.JSEscapeString(s)
	fmt.Println(ss1)

	h := md5.New()
	io.WriteString(h, strconv.FormatInt(int64(time.Now().Nanosecond()), 10))
	io.WriteString(h, "test")
	token := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(token)

	h = md5.New()
	h.Write([]byte("sfdsafdsafdsafds"))
	hmd5 := h.Sum(nil)
	fmt.Printf("%x", hmd5)

	ss := make([]string, 13)
	sss := []string{}
	ss[0] = "aa"
	sss = append(sss, "bbb")
	fmt.Printf("%q\n", ss)
	fmt.Printf("%q\n", sss)

	for _, r := range "Hello 世界！" {
		// 判断字符是否为汉字
		if unicode.Is(unicode.Scripts["Han"], r) {
			fmt.Printf("%c", r) // 世界
		}
	}

	c := make(chan int)
	quit := make(chan int)
	expire := make(chan bool)
	go goRun(c, quit)
	fi(c, quit, expire)
	<-expire

	var ps [20]Person
	for i := 0; i < len(ps); i++ {
		ps[i].Age = rand.Intn(100)
		ps[i].Name = randSeq(5)
		fmt.Printf("%+v\n", ps[i])
	}
	fmt.Printf("%+v\n", ps)

	pm := make(chan Person, 3)
	go personSum(ps[len(ps)/2:], pm)
	go personSum(ps[:len(ps)/2], pm)
	pSum := <-pm
	fmt.Printf("%+v\n", pSum)

	cfb := make(chan int, 10)
	go fib(cap(cfb), cfb)
	for i := range cfb {
		fmt.Println(i)
	}
}
