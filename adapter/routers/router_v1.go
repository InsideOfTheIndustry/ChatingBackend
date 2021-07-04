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
	groupApi "chatting/adapter/controller/group"
	userApi "chatting/adapter/controller/user"
	Middleware "chatting/adapter/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	// 路径管理
	v1 := r.Group("/v1")

	v1.GET("/chatting", chattingApi.ConnectToTcpServer) // 用户聊天时的连接

	user := v1.Group("/user")
	user.POST("/register", userApi.RegisterAccount)              // 用户注册
	user.POST("/verificationcode", userApi.SendVerificationCode) // 发送验证码

	user.POST("/login", userApi.UserLogin)                                                                // 用户登录
	user.POST("/userinfo", Middleware.TokenVerifyWithoutTImeLimit, userApi.GetUserInfo)                   // 获取用户信息
	user.PUT("/useravatar", Middleware.TokenVerify, userApi.UpdateUserAvatar)                             // 修改用户的头像信息
	user.PUT("/userinfo", Middleware.TokenVerify, userApi.UpdateUserInfo)                                 // 修改用户的信息
	user.GET("/usergroup/:useraccount", Middleware.TokenVerifyWithoutTImeLimit, userApi.GetUserGroupInfo) // 获取用户群聊信息
	user.GET("/userfriend/:useraccount", Middleware.TokenVerifyWithoutTImeLimit, userApi.GetFriendInfo)   // 获取用户好友信息

	group := v1.Group("/group")
	group.POST("/newgroupverificationcode", Middleware.TokenVerify, groupApi.SendNewGroupVerificationCode) // 发送创建群聊验证码
	group.POST("/group", Middleware.TokenVerify, groupApi.CreateNewGroup)                                  // 新增一个群
	group.GET("/groupinfo/:groupid", groupApi.QueryGroupInfo)                                              // 查询一个群信息
	group.PUT("/groupinfo", Middleware.TokenVerify, groupApi.UpdateGroupInfo)                              // 更新群聊信息
	group.PUT("/groupavatar", Middleware.TokenVerify, groupApi.UpdateGroupAvatar)                          // 更新群头像
	group.GET("/groupmemberinfo/:groupid", groupApi.QueryGroupMember)                                      // 查询群成员信息
	group.PUT("/usernameingroup", Middleware.TokenVerify, groupApi.UpdateUserNameInGroup)                  // 更新群内昵称

}
