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
	xormgroup "chatting/database/xormdatabase/repository/group"
	xormuser "chatting/database/xormdatabase/repository/user"
	commonservice "chatting/domain/common/service"
	groupservice "chatting/domain/group/service"
	userservice "chatting/domain/user/service"
	"chatting/infrastructure/emailServer"
)

var UserService userservice.UserService
var CommonService commonservice.CommonService
var GroupService groupservice.GroupService

// InitDomainService 将领域服务全部初始化
func InitDomainService() {

	// 新建用户领域服务
	var newuserservice = userservice.UserService{
		UserRepository:      xormuser.UserRepository{XormEngine: xormdatabase.DBEngine},
		UserCacheRepository: redisuser.UserCacheRepository{RedisEngine: redisdatabase.RedisClient},
	}

	// 新建通用服务
	var newcommonservice = commonservice.CommonService{
		CommonCacheRepository: redisuser.UserCacheRepository{RedisEngine: redisdatabase.RedisClient},
		CommonEmailRepository: emailServer.EmailEngineModel,
	}

	// 新建群聊服务
	var newgroupservice = groupservice.GroupService{
		GroupRepository: xormgroup.GroupRepository{XormEngine: xormdatabase.DBEngine},
	}

	UserService = newuserservice
	CommonService = newcommonservice
	GroupService = newgroupservice
}
