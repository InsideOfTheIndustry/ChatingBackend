//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: group_struct.go
// description: 存放访问请求时需要的数据结构
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-07-01
//
package group

import "time"

// UserInfoGet 获取用户信息时的输入结构
type UserInfoGet struct {
	UserAccount int64 `json:"useraccount"` // 用户账号
}

// GroupInfoGet 获取用户信息时的输入结构
type GroupInfoGet struct {
	Groupid int64 `json:"groupid"` // 群聊号
}

// GroupInfoCreate 用于群聊创建和修改
type GroupInfoCreate struct {
	UserAccount      int64     `json:"useraccount"`      // 创建者账号
	GroupName        string    `json:"groupname"`        // 群聊名称
	GroupIntro       string    `json:"groupintro"`       // 群聊简介
	VerificationCode string    `json:"verificationcode"` // 验证码
	CreateAt         time.Time `json:"createat"`         // 创建时间
}

// GroupInfoCreateAndUpdate 用于群聊创建和修改
type GroupInfoUpdate struct {
	UserAccount int64  `json:"useraccount"` // 创建者账号
	GroupName   string `json:"groupname"`   // 群聊名称
	GroupIntro  string `json:"groupintro"`  // 群聊简介
	Groupid     int64  `json:"groupid"`     // 群聊号
}

// GroupInfoUpdateAvatar 用户修改群聊头像时需要的信息
type GroupInfoUpdateAvatar struct {
	UserAccount int64  `json:"useraccount"` // 用户账号
	Groupid     int64  `json:"groupid"`     // 群聊id
	Avatar      string `json:"avatar"`      // 用户上传的头像信息：图片经过base64编码后的内容
}

// UserNameInGroupUpdateInfo 群内用户名更新信息
type UserNameInGroupUpdateInfo struct {
	UserAccount     int64  `json:"useraccount"`     // 用户账号信息
	Groupid         int64  `json:"groupid"`         // 群聊id
	UserNameInGroup string `json:"usernameingroup"` // 群内显示的用户名
}
