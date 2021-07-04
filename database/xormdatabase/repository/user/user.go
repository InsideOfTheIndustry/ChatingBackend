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

const (
	DefaultGroup  int64 = 10000
	DefaultFriend int64 = 100009
)

// UserRepository 用户的dao操作
type UserRepository struct {
	*xormdatabase.XormEngine
}

// Create(user *entity.UserInfo) (int64,error) // 创建新用户 返回用户账号信息
func (ud UserRepository) Create(user *entity.UserInfo) (int64, error) {
	sess := ud.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		logServer.Error("事务启动失败:%s", err.Error())
		return 0, err
	}

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
	_, err := sess.InsertOne(userindatabase)
	if err != nil {
		logServer.Error("创建用户失败：（%s）", err.Error())
		if err := sess.Rollback(); err != nil {
			logServer.Error("事务回滚失败:%s", err.Error())
			return 0, err
		}
		return 0, err
	}

	var usernew = UserInfo{}
	if _, err := sess.Where("useremail = ?", user.UserEmail).Get(&usernew); err != nil {
		logServer.Error("查询用户信息失败:%s", err.Error())
		if err := sess.Rollback(); err != nil {
			logServer.Error("事务回滚失败:%s", err.Error())
			return 0, err
		}
		return 0, err
	}

	// 插入默认群聊
	var usergroup = UserGroup{
		Useraccount:     usernew.UserAccount,
		Groupid:         DefaultGroup,
		UserNameInGroup: usernew.UserName,
	}
	if _, err := sess.InsertOne(usergroup); err != nil {
		logServer.Error("将用户加入默认群聊失败:%s", err.Error())
		if err := sess.Rollback(); err != nil {
			logServer.Error("事务回滚失败:%s", err.Error())
			return 0, err
		}
		return 0, err
	}

	// 加入默认好友
	var defaultfriend = UserFriend{
		Launcher: DefaultFriend,
		Accepter: usernew.UserAccount,
	}

	if _, err := sess.InsertOne(defaultfriend); err != nil {
		logServer.Error("加入默认好友失败:%s", err.Error())
		if err := sess.Rollback(); err != nil {
			logServer.Error("事务回滚失败:%s", err.Error())
			return 0, err
		}
		return 0, err
	}

	if err := sess.Commit(); err != nil {
		return 0, err
	}
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

	if userinfo.Delete == 1 {
		logServer.Info("用户已删除:%v", useraccount)
		return &userinfoentity, nil
	}

	userinfoentity.UserAccount = userinfo.UserAccount
	userinfoentity.UserEmail = userinfo.UserEmail
	userinfoentity.UserName = userinfo.UserName
	userinfoentity.Signature = userinfo.Signature
	userinfoentity.Avatar = userinfo.Avatar
	userinfoentity.UserPassword = userinfo.UserPassword
	userinfoentity.UserAge = userinfo.UserAge
	userinfoentity.UserSex = userinfo.UserSex
	userinfoentity.OwnGroups = userinfo.OwnGroups
	userinfoentity.Online = userinfo.Online

	return &userinfoentity, nil
}

// Update(*entity.UserInfo) error // 更新用户信息 不包括头像信息
func (ud UserRepository) Update(user *entity.UserInfo) error {
	var userindatabase = UserInfo{
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		Signature:   user.Signature,
		UserAge:     user.UserAge,
		UserSex:     user.UserSex,
	}

	if _, err := ud.Where("useraccount = ?", user.UserAccount).Update(userindatabase); err != nil {
		logServer.Error("更新用户失败：（%s）", err.Error())
		return err
	}
	logServer.Error("更新用户成功。")
	return nil
}

// UpdateAvatar 更新用户头像信息
func (ud UserRepository) UpdateAvatar(user *entity.UserInfo) error {
	var userindatabase = UserInfo{
		Avatar:      user.Avatar,
		UserAccount: user.UserAccount,
	}
	_, err := ud.Where("useraccount = ?", user.UserAccount).Update(userindatabase)
	if err != nil {
		logServer.Error("更新用户头像失败：（%s）", err.Error())
		return err
	}
	logServer.Error("更新用户头像成功。")
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
	if err := ud.Where("launcher = ?", useraccount).Find(&friendlauchers); err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}
	if err := ud.Where("accepter = ?", useraccount).Find(&friendaccepters); err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}

	var friendsinfolist = make([]entity.UserInfo, 0)
	// 查询朋友的具体信息
	for i := range friendlauchers {
		var friendinfo = UserInfo{}
		if _, err := ud.Where("useraccount = ?", friendlauchers[i].Accepter).Get(&friendinfo); err != nil {
			logServer.Error("查询出现错误:(%s)", err.Error())
			return friendsinfo, err
		}
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
		if _, err := ud.Where("useraccount = ?", friendaccepters[i].Launcher).Get(&friendinfo); err != nil {
			logServer.Error("查询出现错误:(%s)", err.Error())
			return friendsinfo, err
		}
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

	var friendsinfo = entity.FriendInfo{
		UserAccount: useraccount,
	}

	var friendsinfolistaccepter = make([]entity.UserInfo, 0)
	var friendsinfolistlauncher = make([]entity.UserInfo, 0)

	if err := ud.Table("UserInfo").Join("INNER", "UserFriend", "UserInfo.useraccount = UserFriend.launcher").Where("accepter = ?", useraccount).Find(&friendsinfolistlauncher); err != nil {
		logServer.Error("join出现错误:%s", err.Error())
		return friendsinfo, err
	}

	if err := ud.Table("UserInfo").Join("INNER", "UserFriend", "UserInfo.useraccount = UserFriend.accepter").Where("launcher = ?", useraccount).Find(&friendsinfolistaccepter); err != nil {
		logServer.Error("join出现错误:%s", err.Error())
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

// DeleteUser 删除用户
func (ud UserRepository) DeleteUser(useraccount int64) error {
	var ui = UserInfo{
		UserAccount: useraccount,
		Delete:      1,
	}
	if _, err := ud.Cols("delete").Update(&ui); err != nil {
		logServer.Error("删除用户失败:%s", err.Error())
		return err
	}
	return nil
}

// QueryGroupOfUser 查询用户所在的群
func (ud UserRepository) QueryGroupOfUser(useraccount int64) ([]entity.GroupInfo, error) {
	var usergroupinfo = make([]UserGroup, 0)
	if err := ud.Where("useraccount = ?", useraccount).Find(&usergroupinfo); err != nil {
		logServer.Error("查询用户所在的群失败:%s", err.Error())
		return make([]entity.GroupInfo, 0), err
	}

	var groupidlist = make([]int64, len(usergroupinfo))
	for i := range usergroupinfo {
		groupid := usergroupinfo[i].Groupid
		groupidlist = append(groupidlist, groupid)
	}

	var grouplist = make([]GroupInfo, 0)
	if err := ud.In("groupid", groupidlist).Find(&grouplist); err != nil {
		logServer.Error("查询用户所在群信息失败:%s", err.Error())
	}

	var groupentitylist = make([]entity.GroupInfo, 0)
	for i := range grouplist {
		entitygroup := (entity.GroupInfo)(grouplist[i])
		groupentitylist = append(groupentitylist, entitygroup)
	}

	return groupentitylist, nil
}
