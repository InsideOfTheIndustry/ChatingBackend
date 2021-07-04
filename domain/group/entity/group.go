//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: group.go
// description: 群聊信息实体
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-06-30
//

package entity

import (
	"time"
)

// GroupInfo 群聊信息
type GroupInfo struct {
	Groupid     int64     // 群聊id
	GroupName   string    // 群聊名称
	GroupIntro  string    // 群简介
	GroupOwner  int64     // 群拥有者
	GroupAvatar string    // 群头像
	CreateAt    time.Time // 创建时间
	Deleted     int8      //是否删除
}

// GroupMember 群聊成员
type GroupMember struct {
	Groupid      int64             // 群聊id
	GroupMembers []GroupMemberInfo // 群聊成员
}

// GroupMemberInfo 群友成员信息
type GroupMemberInfo struct {
	UserName        string // 用户名
	UserAccount     int64  // 用户账号
	UserSex         int8   // 用户性别
	Avatar          string // 用户头像
	UserNameInGroup string // 群内用户名
}
