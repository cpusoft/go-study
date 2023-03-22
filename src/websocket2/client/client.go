package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:12312/ws" //服务器地址
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			err := ws.WriteMessage(websocket.BinaryMessage, []byte("ping"))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("receive: ", string(data))
	}
}
