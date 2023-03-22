package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// https://github.com/516134941/websocket-gin-demo/blob/master/clients/client.go
func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8999", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
	time.Sleep(10000 * time.Second)

}
