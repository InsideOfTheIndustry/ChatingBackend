//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: connecttotcp.go
// description: 连接tcp服务
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-11
//

package chattingwithtcp

import (
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/logServer"
	"net"
)

// newTcpConnecton 建立一个新的tcp连接
func newTcpConnecton() (net.Conn, error) {
	tcpconnect, err := net.Dial("tcp", configServer.Tcpservercfg.Ip+":"+configServer.Tcpservercfg.Port)
	if err != nil {
		logServer.Error("连接至Tcp服务器失败:(%s)", err.Error())
		return nil, err
	}
	return tcpconnect, nil

}
