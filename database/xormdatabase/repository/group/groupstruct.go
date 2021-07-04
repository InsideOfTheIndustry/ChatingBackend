//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: groupstruct.go
// description: 群聊相关结构
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-07-01
//
package group

import "time"

// GroupInfo 群聊信息
type GroupInfo struct {
	Groupid     int64     `xorm:"groupid"`     // 群聊号
	GroupName   string    `xorm:"groupname"`   // 群聊名称
	GroupIntro  string    `xorm:"groupintro"`  // 群聊简介
	GroupOwner  int64     `xorm:"groupowner"`  // 群主
	GroupAvatar string    `xorm:"groupavatar"` // 群头像
	Deleted     int8      `xorm:"delete"`      // 是否删除
	CreateAt    time.Time `xorm:"createat"`    // 创建时间
}

// GroupMemberInfo 群内用户的信息
type GroupMemberInfo struct {
	UserName        string `xorm:"username"`        // 用户名
	UserAccount     int64  `xorm:"useraccount"`     // 用户账号
	UserSex         int8   `xorm:"usersex"`         // 用户性别
	Avatar          string `xorm:"avatar"`          // 用户头像
	UserNameInGroup string `xorm:"usernameingroup"` // 用户在群中的昵称
}

// UserGroup 用户-群聊信息表
type UserGroup struct {
	UserNameInGroup string `xorm:"usernameingroup"` // 用户i

}
