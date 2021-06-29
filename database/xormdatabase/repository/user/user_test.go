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
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/logServer"
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	// 读取配置文件
	// configServer.InitConfig("../../../../config/config.json")
	// 测试数据库
	configServer.InitConfig("../../../../config/config.yaml")
	logServer.SetFileLevel("info") // 设置日志等级)
	xormdatabase.InitXormEngine()
	// t, _ := userdao.Create(&userinfo)
	// users, err := userdao.Query(12313131)
	// ifre, err := userdao.QueryEmailIfAlreadyUse("12138")

	// userdao.Update(&userinfo)

	//logServer.Info("用户账号为:%v", account)
	//logServer.Info("用户是否存在:%v,%v", *users, err)
	// logServer.Info("邮箱是否已被注册:%v", ifre)
	ur := UserRepository{xormdatabase.DBEngine}
	h, _ := ur.QueryFriends(100009)
	fmt.Printf("内容为:%v", h)

}
