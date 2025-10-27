package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// go run netcat.go -l -a localhost:8080
// go run netcat.go -a localhost:8080
func main() {
	var listen bool
	var address string

	flag.BoolVar(&listen, "l", false, "Listen mode")
	flag.StringVar(&address, "a", "localhost:8080", "Address to connect/listen to")
	flag.Parse()

	if listen {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listening: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Listening on %s...\n", address)
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting: %v\n", err)
			os.Exit(1)
		}
		handleConnection(conn)
	} else {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting: %v\n", err)
			os.Exit(1)
		}
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Copy stdin to conn and conn to stdout concurrently
	go func() {
		io.Copy(conn, os.Stdin)
		conn.Close()
	}()
	io.Copy(os.Stdout, conn)
}
