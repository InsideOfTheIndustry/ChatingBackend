//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: email_test.go
// description: 邮箱服务函数单元测试
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package emailServer

import (
	"chatting/infrastructure/configServer"
	"testing"
)

func TestSendemail(t *testing.T) {
	configServer.ParseConfig("../../config/config.json")
	config := configServer.GetConfig()
	InitEmailEngine(config)
	EmailEngineModel.SendEmail("12138", "webchatting用户注册验证码", "1121883342@qq.com")
}
