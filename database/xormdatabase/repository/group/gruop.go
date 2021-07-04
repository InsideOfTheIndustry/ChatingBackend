//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: gruop.go
// description: 群聊存储实现
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-07-01
//

package group

import (
	"chatting/database/xormdatabase"
	UserRepo "chatting/database/xormdatabase/repository/user"
	"chatting/domain/group/entity"
	"chatting/infrastructure/logServer"
	"errors"
)

// GroupRepository 群聊存储库具体实现
type GroupRepository struct {
	*xormdatabase.XormEngine
}

// QueryGroupInfo 查询群聊信息
func (gr GroupRepository) QueryGroupInfo(groupid int64) (entity.GroupInfo, error) {
	var groupinfo = GroupInfo{Groupid: groupid}
	ok, err := gr.Get(&groupinfo)
	if err != nil {
		logServer.Error("查询群聊信息失败:%s", err.Error())
		return entity.GroupInfo{}, err
	}
	logServer.Info("ok:%v,value:%v", ok, groupinfo)
	if groupinfo.Deleted == 1 {
		return entity.GroupInfo{}, nil
	}

	var groupinforeturn = entity.GroupInfo{
		Groupid:     groupinfo.Groupid,
		GroupName:   groupinfo.GroupName,
		GroupIntro:  groupinfo.GroupIntro,
		GroupAvatar: groupinfo.GroupAvatar,
		GroupOwner:  groupinfo.GroupOwner,
		CreateAt:    groupinfo.CreateAt,
	}

	return groupinforeturn, nil
}

// CreateGroup 创建群聊
func (gr GroupRepository) CreateGroup(groupinfo entity.GroupInfo) error {
	var infowritein = GroupInfo{
		GroupIntro: groupinfo.GroupIntro,
		GroupName:  groupinfo.GroupName,
		GroupOwner: groupinfo.GroupOwner,
		CreateAt:   groupinfo.CreateAt,
	}
	sess := gr.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		logServer.Error("事务启动失败:%s", err.Error())
		return err
	}

	var user = UserRepo.UserInfo{
		UserAccount: groupinfo.GroupOwner,
	}

	_, err := gr.Get(&user)
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
		logServer.Error("插入新群聊失败:%s", err.Error())
		if err := sess.Rollback(); err != nil {
			logServer.Error("事务回滚失败:%s", err.Error())
			return err
		}
		return err
	}

	if _, err := sess.Get(&infowritein); err != nil {
		logServer.Error("查询新群聊失败:%s", err.Error())
		if err := sess.Rollback(); err != nil {
			logServer.Error("事务回滚失败:%s", err.Error())
			return err
		}
		return err
	}

	// 插入用户 群聊关系 准备加一个时间字段
	var usergroup = UserRepo.UserGroup{
		Useraccount:     groupinfo.GroupOwner,
		Groupid:         infowritein.Groupid,
		UserNameInGroup: user.UserName,
	}
	if _, err := sess.InsertOne(usergroup); err != nil {
		logServer.Error("将用户加入默认群聊失败:%s", err.Error())
		if err := sess.Rollback(); err != nil {
			logServer.Error("事务回滚失败:%s", err.Error())
			return err
		}
		return err
	}

	var userinfo = UserRepo.UserInfo{
		OwnGroups: user.OwnGroups + 1,
	}
	if _, err := sess.Where("useraccount = ?", user.UserAccount).Cols("owngroups").Update(&userinfo); err != nil {

		logServer.Error("更新用户群聊个数失败:%s", err.Error())
		if err := sess.Rollback(); err != nil {
			logServer.Error("事务回滚失败:%s", err.Error())
			return err
		}
		return err
	}

	return sess.Commit()
}

// QueryGroupMember 查询群内用户
func (gr GroupRepository) QueryGroupMember(groupid int64) (entity.GroupMember, error) {
	var groupmembers = make([]GroupMemberInfo, 0)
	var groupmember = entity.GroupMember{
		Groupid: groupid,
	}

	if err := gr.Sql("SELECT username,avatar,UserGroup.useraccount,usersex,usernameingroup FROM UserInfo , UserGroup   WHERE  UserGroup.groupid = ? AND UserInfo.useraccount = UserGroup.useraccount", groupid).Find(&groupmembers); err != nil {
		return groupmember, err
	}

	gm := make([]entity.GroupMemberInfo, len(groupmembers))

	for i := range groupmembers {
		gm[i] = entity.GroupMemberInfo(groupmembers[i])
	}

	groupmember.GroupMembers = gm

	return groupmember, nil
}

// DeleteGroup(groupid int64) error 删除群聊
func (gr GroupRepository) DeleteGroup(groupid int64) error {
	var deletegroupinfo = GroupInfo{
		Groupid: groupid,
		Deleted: 1,
	}
	if _, err := gr.Where("groupid = ?", groupid).Cols("deleted").Update(deletegroupinfo); err != nil {
		logServer.Error("删除群聊失败:%s", err.Error())
		return err
	}
	return nil
}

// UpdateGroupAvatar  修改群头像
func (gr GroupRepository) UpdateGroupAvatar(groupid int64, avatar string) error {
	var groupinfo = GroupInfo{
		Groupid:     groupid,
		GroupAvatar: avatar,
	}
	if _, err := gr.Where("groupid = ?", groupid).Cols("groupavatar").Update(groupinfo); err != nil {
		logServer.Error("更新头像失败:%s", err.Error())
		return err
	}
	return nil
}

// UpdateGroupInfo 更新群聊信息
func (gr GroupRepository) UpdateGroupInfo(groupinfo entity.GroupInfo) error {
	var updategroupinfo = GroupInfo{
		GroupName:  groupinfo.GroupName,
		GroupIntro: groupinfo.GroupIntro,
	}
	if _, err := gr.Where("groupid = ?", groupinfo.Groupid).Update(updategroupinfo); err != nil {
		logServer.Error("更新群聊信息失败:%s", err.Error())
		return err

	}
	return nil
}

// UpdateUserNameInGroup 更新群内用户名
func (gr GroupRepository) UpdateUserNameInGroup(groupid, useraccount int64, usernameingroup string) error {
	var updateinfo = UserGroup{
		UserNameInGroup: usernameingroup,
	}
	if _, err := gr.Table("UserGroup").Where("useraccount = ? ", useraccount).And("groupid = ?", groupid).Cols("usernameingroup").Update(updateinfo); err != nil {
		logServer.Error("更新群内用户名失败:%s", err.Error())
		return err
	}

	return nil
}
