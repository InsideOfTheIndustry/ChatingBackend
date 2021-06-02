//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: user_test.go
// description: 用户表单元测试
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package user

import (
	"chatting/database/xormdatabase"
	"chatting/domain/user/entity"
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/logServer"
	"testing"
)

func TestCreate(t *testing.T) {
	// 读取配置文件
	configServer.InitConfig("../../../../config/config.json")
	// 测试数据库
	var userdao = UserRepository{}
	xormdatabase.InitXormEngine()
	userdao.XormEngine = xormdatabase.DBEngine
	var userinfo = &entity.UserInfo{
		Avatar:      "http://image13.m1905.cn/mdb/uploadfile/2018/1121/thumb_1_300_410_20181121095100733988.jpg",
		UserAccount: 100009,
	}
	// account, _ := userdao.Create(&userinfo)
	// users, err := userdao.Query(12313131)
	// ifre, err := userdao.QueryEmailIfAlreadyUse("12138")
	friends, err := userdao.QueryFriends(10009)
	err = userdao.UpdateAvatar(userinfo)
	// userdao.Update(&userinfo)

	//logServer.Info("用户账号为:%v", account)
	//logServer.Info("用户是否存在:%v,%v", *users, err)
	// logServer.Info("邮箱是否已被注册:%v", ifre)
	logServer.Info("错误信息:%v", err)
	logServer.Info("好友信息:%v", friends)
}
