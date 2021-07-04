//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: userservice.go
// description: 用户相关领域服务
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-14
//

package service

import (
	commonservice "chatting/domain/common/service"
	"chatting/domain/user/entity"
	"chatting/domain/user/repository"

	"errors"
	"math/rand"
	"strconv"
)

//UserService 用户领域服务
type UserService struct {
	UserRepository      repository.UserRepository      // 用户存储库
	UserCacheRepository repository.UserCacheRepository // 用户缓存库
}

//RegisterUser 用户注册
func (us *UserService) RegisterUser(useremail, username, userpassword, signature, avatar string, userage, usersex int64) (int64, error) {
	var user = entity.UserInfo{
		UserEmail:    useremail,
		UserName:     username,
		UserPassword: userpassword,
		Signature:    signature,
		UserAge:      userage,
		UserSex:      usersex,
		Avatar:       avatar,
	}

	// 判断邮箱是否已被注册
	emailifused, err := us.UserRepository.QueryEmailIfAlreadyUse(user.UserEmail)
	if err != nil {
		return 0, err
	}
	if emailifused {
		return 0, nil
	}

	// 创建新用户
	useraccount, err := us.UserRepository.Create(&user)
	if err != nil {
		return 0, err
	}
	return useraccount, nil
}

// IfCouldSendVerifyCodeForRegister 发送邮箱验证码
func (us *UserService) IfCouldSendVerifyCodeForRegister(emailaddr string) (string, error) {

	ifemailhasused, err := us.UserRepository.QueryEmailIfAlreadyUse(emailaddr)

	if err != nil {
		return "", err
	}
	if ifemailhasused {
		return "", errors.New("邮箱已被注册。")
	}

	// 生成验证码
	var messagecode = ""
	for i := 0; i < 4; i++ {
		number := rand.Intn(10)
		word := strconv.Itoa(number)
		messagecode += word
	}

	return messagecode, nil

}

// SendUseraccount 注册成功后向用户发送账号信息
func (us *UserService) SendUseraccount(cs commonservice.CommonService, useraccount int64, emailaddr string) error {
	useraccountmessage := strconv.FormatInt(useraccount, 10)
	err := cs.SendEmail("您的账号为："+useraccountmessage, "webchatting用户账号", emailaddr)
	if err != nil {
		return err
	}
	return err
}

// VerifyLoginInfo 验证用户是否存在
func (us *UserService) VerifyLoginInfo(useraccount int64, userpassword string) (bool, error) {
	userinfo, err := us.UserRepository.Query(useraccount)
	if err != nil {
		return false, err
	}
	if userinfo.UserPassword != userpassword || userinfo.Delete == 1 {
		return false, nil
	}
	return true, nil
}

// GetUserInfo 获取用户信息
func (us *UserService) GetUserInfo(useraccount int64) (entity.UserInfo, error) {
	userinfo, err := us.UserRepository.Query(useraccount)

	return *userinfo, err
}

// GetUserFriendInfo 获取用户好友信息
func (us *UserService) GetUserFriendInfo(useraccount int64) (entity.FriendInfo, error) {
	userFriendInfo, err := us.UserRepository.QueryFriends(useraccount)
	if err != nil {
		return userFriendInfo, err
	}
	for i := range userFriendInfo.Friends {
		userFriendInfo.Friends[i].UserPassword = "************"
	}
	return userFriendInfo, nil
}

// UpdateUserAvatar 更新用户头像
func (us *UserService) UpdateUserAvatar(useraccount int64, avatar string) (bool, error) {
	var userinfo = &entity.UserInfo{
		UserAccount: useraccount,
		Avatar:      avatar,
	}
	if err := us.UserRepository.UpdateAvatar(userinfo); err != nil {
		return false, err
	}

	return true, nil
}

// UpdateUserInfo 更新用户信息 不包括头像
func (us *UserService) UpdateUserInfo(useraccount, userage, usersex int64, username, signature string) error {
	var userinfo = entity.UserInfo{
		UserAccount: useraccount,
		UserAge:     userage,
		UserSex:     usersex,
		UserName:    username,
		Signature:   signature,
	}

	if err := us.UserRepository.Update(&userinfo); err != nil {
		return err
	}
	return nil
}
