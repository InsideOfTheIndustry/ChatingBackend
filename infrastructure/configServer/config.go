//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: config.go
// description: 配置文件读取
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package configServer

import (
	"bufio"
	"chatting/infrastructure/logServer"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

var _config *ConfigStruct

// ParseConfig 解析日志文件
func ParseConfig(path string) (*ConfigStruct, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer jsonFile.Close()

	fileReader := bufio.NewReader(jsonFile)
	byteJsonFile, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byteJsonFile, &_config)
	if err != nil {
		return nil, err
	}
	return _config, nil
}

// GetConfig 获取配置结构
func GetConfig() *ConfigStruct {
	return _config
}

var Applicationcfg = new(Application)     // 基础配置信息
var Emailcfg = new(Email)                 // 邮箱配置信息
var Tcpservercfg = new(TcpServer)         // tcp服务器配置信息
var Databasecfg = new(Database)           // 数据库配置信息
var Rediscfg = new(Redis)                 // redis缓存数据库配置信息
var ResourceStorecfg = new(ResourceStore) // 资源文件存储位置配置信息

// InitConfig 初始化全局的配置信息
func InitConfig(path string) {
	settingCfg := viper.New()
	settingCfg.SetConfigFile(path)
	if err := settingCfg.ReadInConfig(); err != nil {
		logServer.Error("配置文件读取失败:%s", err.Error())
		return
	}

	// 初始化基础配置
	applicationcfg := settingCfg.Sub("settings.application")
	if applicationcfg == nil {
		logServer.Error("找不到application设置")
		return
	}
	Applicationcfg = InitApplicationcfg(applicationcfg)

	// 初始化tcpserver配置
	tcpservercfg := settingCfg.Sub("settings.tcpserver")
	if tcpservercfg == nil {
		logServer.Error("找不到tcpserver设置")
		return
	}
	Tcpservercfg = InitTcpserver(tcpservercfg)

	// 初始化database配置
	databasecfg := settingCfg.Sub("settings.database")
	if databasecfg == nil {
		logServer.Error("找不到database设置")
		return
	}
	Databasecfg = InitDatabase(databasecfg)

	// 初始化redis配置
	rediscfg := settingCfg.Sub("settings.redis")
	if rediscfg == nil {
		logServer.Error("找不到redis设置")
		return
	}
	Rediscfg = InitRedis(rediscfg)

	// 初始化email配置
	emailcfg := settingCfg.Sub("settings.email")
	if emailcfg == nil {
		logServer.Error("找不到email设置")
		return
	}
	Emailcfg = InitEmailcfg(emailcfg)

	// 初始化资源存储位置配置
	resourcestorecfg := settingCfg.Sub("settings.resourcestore")
	if resourcestorecfg == nil {
		logServer.Error("找不到resourcestore设置")
		return
	}
	ResourceStorecfg = InitResourceStore(resourcestorecfg)

}
