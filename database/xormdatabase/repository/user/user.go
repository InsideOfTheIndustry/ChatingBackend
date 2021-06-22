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
	"errors"
)

const (
	DefaultGroup int64 = 10000
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
		sess.Rollback()
		return 0, err
	}

	var usernew = UserInfo{}
	if _, err := sess.Where("useremail = ?", user.UserEmail).Get(&usernew); err != nil {
		logServer.Error("查询用户信息失败:%s", err.Error())
		sess.Rollback()
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
		sess.Rollback()
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

// CreateGroup 创建群聊
func (ud UserRepository) CreateGroup(groupinfo entity.GroupInfo) error {
	var infowritein = GroupInfo{
		GroupIntro: groupinfo.GroupIntro,
		GroupName:  groupinfo.GroupName,
		GroupOwner: groupinfo.GroupOwner,
		CreateAt:   groupinfo.CreateAt,
	}
	sess := ud.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		logServer.Error("事务启动失败:%s", err.Error())
		return err
	}

	user, err := ud.Query(groupinfo.GroupOwner)
	if err != nil {
		logServer.Error("查询用户失败:%s", err.Error())
		return errors.New("your account has problem")
	}
	if user.OwnGroups >= 2 {
		logServer.Error("用户群聊已达上限")
		return errors.New("you can not create more groups")
	}

	// 插入新群聊
	if _, err := sess.InsertOne(infowritein); err != nil {
		sess.Rollback()
		logServer.Error("插入新群聊失败:%s", err.Error())
		return err
	}

	if _, err := sess.Get(&infowritein); err != nil {
		sess.Rollback()
		logServer.Error("查询新群聊失败:%s", err.Error())
		return err
	}

	// 插入用户 群聊关系 准备加一个时间字段
	var usergroup = UserGroup{
		Useraccount:     groupinfo.GroupOwner,
		Groupid:         infowritein.Groupid,
		UserNameInGroup: user.UserName,
	}
	if _, err := sess.InsertOne(usergroup); err != nil {
		logServer.Error("将用户加入默认群聊失败:%s", err.Error())
		sess.Rollback()
		return err
	}

	var userinfo = UserInfo{
		OwnGroups: user.OwnGroups + 1,
	}
	if _, err := sess.Where("useraccount = ?", user.UserAccount).Cols("owngroups").Update(&userinfo); err != nil {
		sess.Rollback()
		logServer.Error("更新用户群聊个数失败:%s", err.Error())
		return err
	}

	return sess.Commit()

}

// UpdateGroup 更新群聊信息
func (ud UserRepository) UpdateGroup(groupinfo entity.GroupInfo) error {
	return nil
}

// QueryGroupInfo 查询群聊信息
func (ud UserRepository) QueryGroupInfo(groupid int64) (entity.GroupInfo, error) {
	var groupinfo = GroupInfo{}
	ok, err := ud.Where("groupid = ?", groupid).Get(&groupinfo)
	if err != nil || !ok {
		logServer.Error("查询群聊信息失败:%s", err.Error())
		return entity.GroupInfo{}, err
	}

	var groupinforeturn = entity.GroupInfo{
		Groupid:     groupinfo.Groupid,
		GroupName:   groupinfo.GroupName,
		GroupIntro:  groupinfo.GroupIntro,
		GroupAvatar: groupinfo.GroupAvatar,
		GroupOwner:  groupinfo.GroupOwner,
	}

	return groupinforeturn, nil
}
