//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: user.go
// description: 用户实体表
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-14
//

package entity

// UserInfo 用户实体
type UserInfo struct {
	UserAccount  int64  // 用户账号
	UserEmail    string // 邮箱号
	UserName     string // 用户名
	Signature    string // 用户个性签名
	Avatar       string // 用户头像
	UserPassword string // 用户密码
	UserAge      int64  // 用户年龄
	UserSex      int64  // 用户性别
	Delete       int8   // 是否被删除
	Online       int8   // 用户在线状态
	OwnGroups    int64  // 创建的群聊个数
}

// UserFriend 朋友间的相互联系
type UserFriend struct {
	Launcher int64 // 好友发起者
	Accepter int64 // 好友接受者
}

// FriendInfo 用户朋友信息
type FriendInfo struct {
	UserAccount int64      // 用户名
	Friends     []UserInfo // 好友信息
}

// GroupInfo 群聊信息
type GroupInfo struct {
	Groupid     int64  // 群聊号
	GroupName   string // 群聊名称
	GroupIntro  string // 群聊简介
	GroupOwner  int64  // 群主
	GroupAvatar string // 群头像
}
