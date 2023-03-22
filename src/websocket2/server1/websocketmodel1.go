package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// https://blog.csdn.net/qq_42887507/article/details/120230212
// https://blog.csdn.net/ALakers/article/details/111713405
// https://github.com/lwnmengjing/pushMessage
// https://zhuanlan.zhihu.com/p/489023088
// https://blog.csdn.net/takujo/article/details/104083799
// https://github.com/516134941/websocket-gin-demo
var webSocketServer *WebSocketServer

func init() {
	webSocketServer = &WebSocketServer{}
	webSocketServer.clientConns = make(map[string]*websocket.Conn)
	webSocketServer.clientChans = make(map[string]chan []byte)
}

type WebSocketServer struct {
	// 互斥锁，防止程序对统一资源同时进行读写
	mux sync.Mutex
	// websocket客户端链接池
	clientConns map[string]*websocket.Conn
	// 消息通道
	clientChans map[string]chan []byte
}

// 将客户端添加到客户端链接池
func (c *WebSocketServer) addClient(clientIpPort string, conn *websocket.Conn) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.clientConns[clientIpPort] = conn
}

// 获取指定客户端链接
func (c *WebSocketServer) getClient(clientIpPort string) (conn *websocket.Conn, exist bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	conn, exist = c.clientConns[clientIpPort]
	return
}

// 删除客户端链接
func (c *WebSocketServer) deleteClient(clientIpPort string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.clientConns, clientIpPort)
}

// 添加用户消息通道
func (c *WebSocketServer) addChannel(clientIpPort string, m chan []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.clientChans[clientIpPort] = m

}

// 获取指定用户消息通道
func (c *WebSocketServer) getChannel(clientIpPort string) (m chan []byte, exist bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	m, exist = c.clientChans[clientIpPort]

	return
}

// 删除指定消息通道
func (c *WebSocketServer) deleteChannel(clientIpPort string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if m, ok := c.clientChans[clientIpPort]; ok {
		close(m)
		delete(c.clientChans, clientIpPort)
	}
}

// websocket Upgrader
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
