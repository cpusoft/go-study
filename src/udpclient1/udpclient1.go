package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	udpAddr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:9998")

	//连接udpAddr，返回 udpConn
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("udp dial ok ")
	bufchan := make(chan []byte, 2048)
	go func() {
		for {
			buf := make([]byte, 1024)
			length, _ := udpConn.Read(buf)
			fmt.Println("go client read data:", string(buf[:length]))
			str := string(buf[:length])

			if strings.HasPrefix(str, "1") {
				bufchan <- buf[:length]
			} else {
				fmt.Println("go client self process buf:", str)
			}
		}
	}()

	for i := 0; i < 3; i++ {
		// 发送数据
		len, err := udpConn.Write([]byte("AAAAAA"))
		if err != nil {
			return
		}

		buf := <-bufchan
		fmt.Println("client read data from go :", len, string(buf))
	}
	time.Sleep(1000 * time.Second)
	/*
		//读取数据
		buf := make([]byte, 1024)
		len, _ = udpConn.Read(buf)
		fmt.Println("client read len:", len)
		fmt.Println("client read data:", convert.PrintBytesOneLine(buf[:len]))

			// 发送数据
			len, err = udpConn.Write([]byte("BBBBB"))
			if err != nil {
				return
			}
			fmt.Println("client write len:", len)

			//读取数据
			buf = make([]byte, 1024)
			len, _ = udpConn.Read(buf)

			fmt.Println("client read len:", len)
			fmt.Println("client read data:", convert.PrintBytesOneLine(buf[:len]))
	*/
}
