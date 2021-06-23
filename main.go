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
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	// 读取配置文件
	configServer.InitConfig("./config/config.yaml")

	logServer.Info("%s", configServer.ResourceStorecfg.UserAvatar)
	// 资源文件夹开放
	engine.StaticFS("/useravatar", http.Dir(configServer.ResourceStorecfg.UserAvatar)) // 用户头像位置

	// 初始化
	engine.Use(cors.Cors())        // 解决跨域问题
	logServer.SetFileLevel("info") // 设置日志等级
	routers.InitRouter(engine)     // 初始化路由
	if err := xormdatabase.InitXormEngine(); err != nil {
		logServer.Error("数据库引擎启动失败:%s", err.Error())
		panic(err.Error())
	} // 初始化mysql数据库连接引擎
	redisdatabase.InitRedis()      // 初始化redis数据库
	emailServer.InitEmailEngine()  // 初始化邮件服务
	controller.InitDomainService() // 初始化领域服务

	err := engine.Run(configServer.Applicationcfg.Ip + ":" + configServer.Applicationcfg.Port)
	if err != nil {
		logServer.Error("服务启动失败(%s)...", err.Error())
	}
}
