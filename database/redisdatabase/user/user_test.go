//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: user_test.go
// description: redis数据库测试
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//
package user

import (
	"chatting/database/redisdatabase"
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/logServer"
	"testing"
)

func TestRedis(t *testing.T) {
	// 读取配置文件
	configServer.InitConfig("../../../config/config.json")

	redisdatabase.InitRedis()

	var userdao = UserCacheRepository{
		redisdatabase.RedisClient,
	}

	userdao.SetVerificationCode("1121883342@qq.com", "你看看")
	msg, _ := userdao.GetVerificationCode("1121883342@qq.com")
	logServer.Info("数据为:%s", msg)

}
