//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: config_struct.go
// description: 日数据结构定义
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package configServer

// 配置文件结构
type ConfigStruct struct {
	Ip         string    `json:"ip"`         // ip地址
	Port       string    `json:"port"`       // 端口号
	ConfigPath string    `json:"configpath"` // 日志路由
	Tcpserver  TcpServer `json:"tcpserver"`  // tcp服务器配置
	Database   Database  `json:"database"`   // 数据库配置文件
	Email      Email     `json:"email"`      // 邮箱设置
	Redis      Redis     `json:"redis"`      // Redis缓存数据库
}

// Applicationcfg 基础的配置信息，包括本项目运行的端口和ip applicationcfg
type Application struct {
	Ip       string `json:"ip"`   // ip地址
	Port     string `json:"port"` // 端口号
	ServerIp string // 程序运行所处服务的ip
}

// tcp服务器ip和port  tcpserver
type TcpServer struct {
	Ip   string // tcp服务器的ip地址 ip
	Port string // tcp服务器的端口 port
}

//Database 数据库配置文件 database
type Database struct {
	Type     string // 数据库类型 mysql.. type
	User     string // 用户名 user
	Password string // 密码 password
	Host     string // IP host
	Port     string // 端口 port
	DBName   string // 数据库名 dbName
	Charset  string // 编码方式 charset
	Showsql  bool   // 是否显示数据库查询语句 showsql
}

// Email 邮箱配置 email
type Email struct {
	User     string // 邮箱账号 user
	Password string // 密码 password
	Host     string // 邮箱服务器地址 host
}

// Redis缓存数据库 redis
type Redis struct {
	Addr     string // 数据库地址 addr
	Password string // 数据库连接密码 password
	Db       int    // 使用的数据库 db
}

// ResourceStore 资源文件存储位置  resourcestore
type ResourceStore struct {
	UserAvatar string // 用户头像存储位置 useravatar
}
