//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: configinit.go
// description: 配置信息初始化
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-06-01
//

package configServer

import "github.com/spf13/viper"

//InitTcpserver 初始化tcp服务器配置
func InitApplicationcfg(appcfg *viper.Viper) *Application {
	return &Application{
		Port:     appcfg.GetString("port"),
		Ip:       appcfg.GetString("ip"),
		ServerIp: appcfg.GetString("serverip"),
		JwtKey:   appcfg.GetString("jwtkey"),
	}
}

//InitTcpserver 初始化tcp服务器配置
func InitTcpserver(tcpcfg *viper.Viper) *TcpServer {
	return &TcpServer{
		Port: tcpcfg.GetString("port"),
		Ip:   tcpcfg.GetString("ip"),
	}
}

//InitDatabase 初始化数据库配置
func InitDatabase(databasecfg *viper.Viper) *Database {
	return &Database{
		Type:     databasecfg.GetString("type"),
		User:     databasecfg.GetString("user"),
		Password: databasecfg.GetString("password"),
		Host:     databasecfg.GetString("host"),
		Port:     databasecfg.GetString("port"),
		DBName:   databasecfg.GetString("dbName"),
		Charset:  databasecfg.GetString("charset"),
		Showsql:  databasecfg.GetBool("showsql"),
	}
}

// InitRedis 初始化redis
func InitRedis(rediscfg *viper.Viper) *Redis {
	return &Redis{
		Addr:     rediscfg.GetString("addr") + ":" + rediscfg.GetString("port"),
		Password: rediscfg.GetString("password"),
		Db:       rediscfg.GetInt("db"),
	}
}

// InitEmailcfg 初始化redis
func InitEmailcfg(emailcfg *viper.Viper) *Email {
	return &Email{
		User:     emailcfg.GetString("user"),
		Password: emailcfg.GetString("password"),
		Host:     emailcfg.GetString("host"),
	}
}

// InitResourceStore 初始化资源文件存储位置配置信息
func InitResourceStore(resourcestorecfg *viper.Viper) *ResourceStore {
	return &ResourceStore{
		UserAvatar:  resourcestorecfg.GetString("useravatar"),
		GroupAvatar: resourcestorecfg.GetString("groupavatar"),
	}
}
