package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

func Echo(ws *websocket.Conn) {
	for {
		var reply string
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Cannot receive", err)
			break
		}
		fmt.Println("Recevie :" + reply)
		msg := "Received:" + reply + "  on " + time.Now().String()
		if err := websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Cannot Send", err)
			break
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo))
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe fail,", err)
	}
}
