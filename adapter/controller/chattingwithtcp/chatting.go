//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: chatting.go
// description: 将http连接升级为websocket连接 直接与TCP服务器通信的控制器
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-11
//

package chattingwithtcp

import (
	"net"
	"net/http"

	"chatting/infrastructure/logServer"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ConnectToTcpServer 将http服务升级为websocket 进行通信
func ConnectToTcpServer(c *gin.Context) {

	// 与前端建立的websocket连接
	webclient, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logServer.Error("升级为websocket服务失败：(%s)", err.Error())
		return
	}

	// 建立一个与tcp的连接
	tcpconn, err := newTcpConnecton()
	if err != nil {
		return
	}

	go fronttoend(webclient, tcpconn)

	endtofront(webclient, tcpconn)

}

//fronttoend 从前端发送至后端
func fronttoend(webclient *websocket.Conn, tcpconn net.Conn) {
	for {
		_, message, err := webclient.ReadMessage()
		if err != nil {
			logServer.Error("从前端接收数据失败：（%s）", err.Error())
			webclient.Close()
			tcpconn.Close()
			return
		}

		_, err = tcpconn.Write(message)
		if err != nil {
			logServer.Error("数据转发至TCP服务器失败：（%s）", err.Error())
			webclient.Close()
			tcpconn.Close()
			return
		}
	}
}

//endtofront 从后端发送至前端
func endtofront(webclient *websocket.Conn, tcpconn net.Conn) {
	for {
		var message = make([]byte, 1024*2)
		count, err := tcpconn.Read(message)
		if err != nil {
			logServer.Error("从后端获取数据失败：(%s)", err.Error())
			webclient.Close()
			tcpconn.Close()
			return
		}

		err = webclient.WriteMessage(websocket.TextMessage, message[:count])
		if err != nil {
			logServer.Error("数据转发至前端失败：(%s)", err.Error())
			webclient.Close()
			tcpconn.Close()
			return
		}

	}
}
