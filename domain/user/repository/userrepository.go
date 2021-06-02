//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: userrepository.go
// description: 用户领域相关存储库信息
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-14
//

package repository

import (
	"chatting/domain/user/entity"
)

//UserRepository 用户存储库
type UserRepository interface {
	Create(user *entity.UserInfo) (int64, error)               // 创建新用户 返回用户账号信息
	Update(user *entity.UserInfo) error                        // 更新用户信息（此处不包括头像）
	Query(useraccount int64) (*entity.UserInfo, error)         // 查询用户信息
	QueryEmailIfAlreadyUse(email string) (bool, error)         // 查询邮箱是否已经注册
	QueryFriends(useraccount int64) (entity.FriendInfo, error) // 查询用户好友信息
	UpdateAvatar(user *entity.UserInfo) error                  // 修改用户头像信息
}

//UserCacheRepository 用户缓存存储库
type UserCacheRepository interface {
	SetVerificationCode(emailaddr, verificationcode string) error // 设置邮箱验证码
	GetVerificationCode(emailaddr string) (string, error)         // 获取邮箱验证码
	SetToken(useraccount int64, token string) error               // 设置token信息
	GetToken(useraccount int64) (string, error)                   // 获取token信息
}

// UserEmailServer 用户邮箱服务接口
type UserEmailServer interface {
	SendEmail(message, subject, receiver string) error // 发送邮件
}
