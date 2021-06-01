//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: main.go
// description: 主文件
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-11
//

package main

import (
	"chatting/adapter/controller"
	"chatting/adapter/middleware/cors"
	"chatting/adapter/routers"
	"chatting/database/redisdatabase"
	"chatting/database/xormdatabase"
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/emailServer"
	"chatting/infrastructure/logServer"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	// 读取配置文件
	configServer.ParseConfig("./config/config.json")
	config := configServer.GetConfig()

	// 初始化
	engine.Use(cors.Cors())             // 解决跨域问题
	logServer.SetFileLevel("info")      // 设置日志等级
	routers.InitRouter(engine)          // 初始化路由
	xormdatabase.InitXormEngine(config) // 初始化mysql数据库连接引擎
	redisdatabase.InitRedis(config)     // 初始化redis数据库
	emailServer.InitEmailEngine(config) // 初始化邮件服务
	controller.InitDomainService()      // 初始化领域服务

	err := engine.Run(config.Ip + ":" + config.Port)
	if err != nil {
		logServer.Error("服务启动失败(%s)...", err.Error())
	}
}
