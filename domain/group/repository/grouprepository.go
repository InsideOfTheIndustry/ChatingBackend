//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: grouprepository.go
// description: 群聊信息管理存储库
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-06-30
//

package repository

import "chatting/domain/group/entity"

// GroupRepository 群聊仓储库
type GroupRepository interface {
	CreateGroup(groupinfo entity.GroupInfo) error                                   // 创建群聊
	UpdateGroupInfo(groupinfo entity.GroupInfo) error                               // 修改群聊信息
	UpdateGroupAvatar(groupid int64, avatar string) error                           // 更新群聊头像
	DeleteGroup(groupid int64) error                                                // 删除群聊
	QueryGroupInfo(groupid int64) (entity.GroupInfo, error)                         // 查询群聊信息
	QueryGroupMember(groupid int64) (entity.GroupMember, error)                     // 查找群成员信息
	UpdateUserNameInGroup(groupid, useraccount int64, usernameingroup string) error // 更新本人群内昵称
}

// GroupCacheRepository 群聊缓存库
type GroupCacheRepository interface {
}
