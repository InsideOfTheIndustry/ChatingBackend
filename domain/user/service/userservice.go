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
	"time"
)

//UserService 用户领域服务
type UserService struct {
	UserRepository      repository.UserRepository      // 用户存储库
	UserCacheRepository repository.UserCacheRepository // 用户缓存库
	UserEmailServer     repository.UserEmailServer     // 邮件服务
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
	if userinfo.UserPassword != userpassword || userinfo.Delete == 1 {
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

// CreateNewGroup 新建一个群聊
func (us *UserService) CreateNewGroup(groupname, groupintro, verificationcode string, createat time.Time, groupowner int64) error {

	userinfo, err := us.UserRepository.Query(groupowner)

	if err != nil {
		return err
	}

	vc, err := us.UserCacheRepository.GetVerificationCode(userinfo.UserEmail + "group")
	if err != nil {
		return err
	}

	if vc != verificationcode {
		return errors.New("验证码错误或失效!")
	}

	var groupinfo = entity.GroupInfo{
		GroupName:  groupname,
		GroupIntro: groupintro,
		GroupOwner: groupowner,
		CreateAt:   createat,
	}
	if err := us.UserRepository.CreateGroup(groupinfo); err != nil {
		return err
	}
	return nil
}

// SendCreateNewGroupVerifyCode 在创建群聊时向对方邮箱发送验证码
func (us *UserService) SendCreateNewGroupVerifyCode(useraccount int64) error {

	userinfo, err := us.UserRepository.Query(useraccount)

	if err != nil {
		return err
	}

	if userinfo.OwnGroups >= 2 {
		return errors.New("number of the group you created allready more than 2, you can not create new group")
	}

	// 生成验证码
	var messagecode = ""
	for i := 0; i < 4; i++ {
		number := rand.Intn(10)
		word := strconv.Itoa(number)
		messagecode += word
	}

	// 设置验证码缓存
	err = us.UserCacheRepository.SetVerificationCode(userinfo.UserEmail+"group", messagecode)
	if err != nil {
		return err
	}

	// 发送验证码
	if err := us.UserEmailServer.SendEmail(messagecode, "webchatting群聊创建验证码", userinfo.UserEmail); err != nil {
		return err
	}
	return err
}

// JudgeRequestFrequence 判断请求访问频率
func (us *UserService) JudgeRequestFrequence(path string, fretime int64) (bool, error) {
	ifexist, err := us.UserCacheRepository.GetRequestInfo(path)

	// 出现错误返回错误
	if err != nil {
		return false, err
	}
	// 表示请求访问过于频繁
	if ifexist {
		return true, nil
	}
	// 合理时间内的访问
	if err := us.UserCacheRepository.SetRequestInfo(path, "5", fretime); err != nil {
		return false, err
	}

	return false, err

}
