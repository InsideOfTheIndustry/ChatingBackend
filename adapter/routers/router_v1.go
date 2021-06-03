//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: router_v1.go
// description: api路由
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-11
//

package routers

import (
	chattingApi "chatting/adapter/controller/chattingwithtcp"
	userApi "chatting/adapter/controller/user"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	// 路径管理
	v1 := r.Group("/v1")

	v1.GET("/chatting", chattingApi.ConnectToTcpServer) // 用户聊天时的连接

	user := v1.Group("/user")
	user.POST("/register", userApi.RegisterAccount)                        // 用户注册
	user.POST("/verificationcode", userApi.SendVerificationCode)           // 发送验证码
	user.POST("/login", userApi.UserLogin)                                 // 用户登录
	user.POST("/userinfo", userApi.TokenVerify, userApi.GetUserInfo)       // 获取用户信息
	user.POST("/userfriend", userApi.TokenVerify, userApi.GetFriendInfo)   // 获取用户好友信息
	user.PUT("/useravatar", userApi.TokenVerify, userApi.UpdateUserAvatar) // 修改用户的头像信息
	user.PUT("/userinfo", userApi.TokenVerify, userApi.UpdateUserInfo)     // 修改用户的信息

}
