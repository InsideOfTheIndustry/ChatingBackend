//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: user.go
// description: 具体的数据库操作实现
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package user

import (
	"chatting/database/xormdatabase"
	"chatting/domain/user/entity"
	"chatting/infrastructure/logServer"
)

// UserRepository 用户的dao操作
type UserRepository struct {
	*xormdatabase.XormEngine
}

// Create(user *entity.UserInfo) (int64,error) // 创建新用户 返回用户账号信息
func (ud UserRepository) Create(user *entity.UserInfo) (int64, error) {
	var userindatabase = UserInfo{
		UserAccount:  user.UserAccount,
		UserEmail:    user.UserEmail,
		UserName:     user.UserName,
		Signature:    user.Signature,
		Avatar:       user.Avatar,
		UserPassword: user.UserPassword,
		UserAge:      user.UserAge,
		UserSex:      user.UserSex,
	}
	_, err := ud.InsertOne(userindatabase)
	if err != nil {
		logServer.Error("创建用户失败：（%s）", err.Error())
		return 0, err
	}

	var usernew = UserInfo{}
	ud.Where("useremail = ?", user.UserEmail).Get(&usernew)

	return usernew.UserAccount, nil

}

// Query(useraccount int64) (*entity.UserInfo, error) // 查询用户信息
func (ud UserRepository) Query(useraccount int64) (*entity.UserInfo, error) {
	var userinfo = UserInfo{}
	_, err := ud.Where("useraccount = ?", useraccount).Get(&userinfo)

	var userinfoentity = entity.UserInfo{}

	if err != nil {
		logServer.Error("查询数据出错:(%s)", err.Error())
		return &userinfoentity, err
	}
	userinfoentity.UserAccount = userinfo.UserAccount
	userinfoentity.UserEmail = userinfo.UserEmail
	userinfoentity.UserName = userinfo.UserName
	userinfoentity.Signature = userinfo.Signature
	userinfoentity.Avatar = userinfo.Avatar
	userinfoentity.UserPassword = userinfo.UserPassword
	userinfoentity.UserAge = userinfo.UserAge
	userinfoentity.UserSex = userinfo.UserSex

	return &userinfoentity, nil
}

// Update(*entity.UserInfo) error // 更新用户信息 不包括头像信息
func (ud UserRepository) Update(user *entity.UserInfo) error {
	var userindatabase = UserInfo{
		UserAccount:  user.UserAccount,
		UserName:     user.UserName,
		Signature:    user.Signature,
		UserPassword: user.UserPassword,
		UserAge:      user.UserAge,
		UserSex:      user.UserSex,
	}
	_, err := ud.Where("useraccount = ?", user.UserAccount).Update(userindatabase)
	if err != nil {
		logServer.Error("更新用户失败：（%s）", err.Error())
		return err
	}
	logServer.Error("更新用户成功。")
	return nil
}

func (ud UserRepository) UpadateAvatar(user *entity.UserInfo) error {
	var userindatabase = UserInfo{
		Avatar:      user.Avatar,
		UserAccount: user.UserAccount,
	}
	_, err := ud.Where("useraccount = ?", user.UserAccount).Update(userindatabase)
	if err != nil {
		logServer.Error("更新用户失败：（%s）", err.Error())
		return err
	}
	logServer.Error("更新用户成功。")
	return nil
}

// QueryFriends(useraccount int64)([]entity.FriendInfo, error) // 查询用户好友信息 复杂了 下面改进了
func (ud UserRepository) QueryFriendsComplex(useraccount int64) (entity.FriendInfo, error) {
	var friendlauchers = make([]UserFriend, 0)
	var friendaccepters = make([]UserFriend, 0)
	var friendsinfo = entity.FriendInfo{
		UserAccount: useraccount,
	}

	// 查询朋友信息时 需要从发起者和接收者两处查询
	err := ud.Where("launcher = ?", useraccount).Find(&friendlauchers)
	if err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}
	err = ud.Where("accepter = ?", useraccount).Find(&friendaccepters)
	if err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}

	var friendsinfolist = make([]entity.UserInfo, 0)
	// 查询朋友的具体信息
	for i := range friendlauchers {
		var friendinfo = UserInfo{}
		ud.Where("useraccount = ?", friendlauchers[i].Accepter).Get(&friendinfo)
		var entityuserinfo = entity.UserInfo{
			UserAccount: friendinfo.UserAccount,
			UserName:    friendinfo.UserName,
			Signature:   friendinfo.Signature,
			Avatar:      friendinfo.Avatar,
			UserAge:     friendinfo.UserAge,
			UserSex:     friendinfo.UserSex,
		}
		friendsinfolist = append(friendsinfolist, entityuserinfo)
	}

	for i := range friendaccepters {
		var friendinfo = UserInfo{}
		ud.Where("useraccount = ?", friendaccepters[i].Launcher).Get(&friendinfo)
		var entityuserinfo = entity.UserInfo{
			UserAccount: friendinfo.UserAccount,
			UserName:    friendinfo.UserName,
			Signature:   friendinfo.Signature,
			Avatar:      friendinfo.Avatar,
			UserAge:     friendinfo.UserAge,
			UserSex:     friendinfo.UserSex,
		}
		friendsinfolist = append(friendsinfolist, entityuserinfo)
	}
	friendsinfo.Friends = friendsinfolist
	return friendsinfo, nil

}

// QueryFriends(useraccount int64)([]entity.FriendInfo, error) // 查询用户好友信息
func (ud UserRepository) QueryFriends(useraccount int64) (entity.FriendInfo, error) {
	var friendlauchers = make([]UserFriend, 0)
	var friendaccepters = make([]UserFriend, 0)
	var friendsinfo = entity.FriendInfo{
		UserAccount: useraccount,
	}

	// 查询朋友信息时 需要从发起者和接收者两处查询
	if err := ud.Where("launcher = ?", useraccount).Find(&friendlauchers); err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}

	if err := ud.Where("accepter = ?", useraccount).Find(&friendaccepters); err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}

	// 获取被定义为接受者的朋友
	var accepter = make([]int64, 0, len(friendlauchers))
	for i := range friendlauchers {
		accepter = append(accepter, friendlauchers[i].Accepter)
	}
	// 获取定义为发起者的朋友
	var launcher = make([]int64, 0, len(friendaccepters))
	for i := range friendaccepters {
		launcher = append(launcher, friendaccepters[i].Launcher)
	}

	var friendsinfolistaccepter = make([]entity.UserInfo, 0, len(accepter))
	var friendsinfolistlauncher = make([]entity.UserInfo, 0, len(launcher))

	if err := ud.In("useraccount", accepter).Find(&friendsinfolistaccepter); err != nil {
		logServer.Error("查询朋友信息出错:%s", err.Error())
		return friendsinfo, err
	}
	if err := ud.In("useraccount", launcher).Find(&friendsinfolistlauncher); err != nil {
		logServer.Error("查询朋友信息出错:%s", err.Error())
		return friendsinfo, err
	}
	friendsinfo.Friends = append(friendsinfo.Friends, friendsinfolistaccepter...)
	friendsinfo.Friends = append(friendsinfo.Friends, friendsinfolistlauncher...)
	return friendsinfo, nil

}

// QueryEmailIfAlreadyUse(email string) (bool, error)           // 查询邮箱是否已经注册
func (ud UserRepository) QueryEmailIfAlreadyUse(email string) (bool, error) {
	var userinfo = UserInfo{}
	count, err := ud.Where("useremail = ?", email).Count(userinfo)
	if err != nil {
		logServer.Error("查询邮箱是否被注册出现错误:(%s)", err.Error())
		return true, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
