package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	service := "www.baidu.com:80"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)
	fmt.Println("tcpAddr :")
	typeof(tcpAddr)

	myConn, err1 := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err1)
	fmt.Println("myConn :")
	typeof(myConn)
	fmt.Println(myConn.RemoteAddr())
	fmt.Println(myConn.LocalAddr())

	_, err = myConn.Write([]byte("HEAD / HTTP/1.1\r\n\r\n"))
	checkError(err)

	result, err := ioutil.ReadAll(myConn)
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}

func typeof(v interface{}) {
	fmt.Printf("type is:%T\n", v)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
