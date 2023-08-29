package main

import (
	"fmt"
	"time"

	"github.com/cpusoft/goutil/ginserver"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	r := gin.Default()
	r.GET("/ws", WsConnect)
	r.Run("localhost:8080")
}
func WsConnect(c *gin.Context) {
	clientIpPort, _ := ginserver.GetClientIpPort(c)
	fmt.Println(clientIpPort)
	// 升级为websocket长链接
	WsHandler(c, clientIpPort)
}
func DeleteClient(c *gin.Context) {
	clientIpPort, _ := ginserver.GetClientIpPort(c)
	fmt.Println(clientIpPort)
	// 关闭websocket链接
	conn, exist := webSocketServer.getClient(clientIpPort)
	if exist {
		conn.Close()
		webSocketServer.deleteClient(clientIpPort)
	} else {
		//ginserver.
	}
	// 关闭其消息通道
	_, exist = webSocketServer.getChannel(clientIpPort)
	if exist {
		webSocketServer.deleteChannel(clientIpPort)
	}
}

func WsHandler(c *gin.Context, clientIpPort string) {
	var conn *websocket.Conn
	var err error
	var exist bool

	// 创建一个定时器用于服务端心跳
	pingTicker := time.NewTicker(time.Second * 10)
	conn, err = wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 把与客户端的链接添加到客户端链接池中
	webSocketServer.addClient(clientIpPort, conn)

	// 获取该客户端的消息通道
	m, exist := webSocketServer.getChannel(clientIpPort)
	if !exist {
		m = make(chan []byte)
		webSocketServer.addChannel(clientIpPort, m)
	}

	// 设置客户端关闭ws链接回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		webSocketServer.deleteClient(clientIpPort)
		fmt.Println(code)
		return nil
	})

	for {
		select {
		case content, _ := <-m:
			// 从消息通道接收消息，然后推送给前端
			err = conn.WriteJSON(content)
			if err != nil {
				fmt.Println(err)
				conn.Close()
				webSocketServer.deleteClient(clientIpPort)
				return
			}
		case <-pingTicker.C:
			// 服务端心跳:每20秒ping一次客户端，查看其是否在线
			conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
			err = conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				fmt.Println("send ping err:", err)
				conn.Close()
				webSocketServer.deleteClient(clientIpPort)
				return
			}
		}
	}
}
