//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: serverInit.go
// description: 初始化领域服务
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-23
//

package controller

import (
	"chatting/database/redisdatabase"
	redisuser "chatting/database/redisdatabase/user"
	"chatting/database/xormdatabase"
	xormuser "chatting/database/xormdatabase/repository/user"
	"chatting/domain/user/service"
	"chatting/infrastructure/emailServer"
)

var UserService service.UserService

// InitDomainService 将领域服务全部初始化
func InitDomainService() {

	// 新建用户领域服务
	var newuserservice = service.UserService{
		UserRepository: xormuser.UserRepository{XormEngine: xormdatabase.DBEngine},
		UserCacheRepository: redisuser.UserCacheRepository{RedisEngine: redisdatabase.RedisClient},
		UserEmailServer: emailServer.EmailEngineModel,
	}

	UserService = newuserservice
}