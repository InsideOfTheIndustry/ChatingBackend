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
	UserEmailServer     repository.UserEmailServer     // 邮件服务
}

//RegisterUser 用户注册
func (us *UserService) RegisterUser(user entity.UserInfo) (int64, error) {

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

// SendVerificationCode 发送邮箱验证码
func (us *UserService) SendVerificationCode(emailaddr string) error {

	ifemailhasused, err := us.UserRepository.QueryEmailIfAlreadyUse(emailaddr)

	if err != nil {
		return err
	}
	if ifemailhasused {
		return errors.New("邮箱已被注册。")
	}

	// 生成验证码
	var messagecode = ""
	for i := 0; i < 4; i++ {
		number := rand.Intn(10)
		word := strconv.Itoa(number)
		messagecode += word
	}
	// 设置验证码缓存
	err = us.UserCacheRepository.SetVerificationCode(emailaddr, messagecode)
	if err != nil {
		return err
	}
	// 发送验证码
	if err := us.UserEmailServer.SendEmail(messagecode, "webchatting用户注册验证码", emailaddr); err != nil {
		return err
	}
	return err
}

// SendUseraccount 注册成功后向用户发送账号信息
func (us *UserService) SendUseraccount(useraccount int64, emailaddr string) error {
	useraccountmessage := strconv.FormatInt(useraccount, 10)
	err := us.UserEmailServer.SendEmail("您的账号为："+useraccountmessage, "webchatting用户账号", emailaddr)
	if err != nil {
		return err
	}
	return err
}

// VerifyCode验证邮箱验证码是否正确
func (us *UserService) VerifyCode(code, emailaddr string) (bool, error) {
	codeincache, err := us.UserCacheRepository.GetVerificationCode(emailaddr)
	if err != nil {
		return false, err
	}
	if code == codeincache {
		return true, nil
	}
	return false, nil
}

// VerifyLoginInfo 验证用户是否存在
func (us *UserService) VerifyLoginInfo(useraccount int64, userpassword string) (bool, error) {
	userinfo, err := us.UserRepository.Query(useraccount)
	if err != nil {
		return false, err
	}
	if userinfo.UserPassword != userpassword {
		return false, nil
	}
	return true, nil
}

// GetUserFriendInfo 获取用户好友信息
func (us *UserService) GetUserFriendInfo(useraccount int64) (entity.FriendInfo, error) {
	userFriendInfo, err := us.UserRepository.QueryFriends(useraccount)
	if err != nil {
		return userFriendInfo, err
	}

	return userFriendInfo, nil
}
