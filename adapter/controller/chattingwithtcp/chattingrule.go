//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: chattingrule.go
// description: 通信协议 即数据 与tcp服务器中同
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-11
//

package chattingwithtcp

// Message 信息传递结构
type Message struct {
	MessageType int `json:"messagetype"`// 消息的类型 比如心跳、发送信息等
	Token string `json:"token"` // token用于验证用户是否登录
	Message string `json:"message"` // 发送的信息
	Sender string `json:"sender"` // 发送者账号
	Receiver string `json:"receiver"` // 接收者账号
}

/*
具体的协议类型定义： 初版
MessageType：0，含义：心跳
MessageType: 1, 含义：初次连接 
MessageType：2，含义：发送信息
MessageType：4，含义：断开连接
*/