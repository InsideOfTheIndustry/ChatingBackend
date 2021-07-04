//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: group.go
// description: 群聊提供的服务
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-06-30
//

package service

import (
	"chatting/domain/group/entity"
	"chatting/domain/group/repository"
	"time"
)

type GroupService struct {
	GroupRepository repository.GroupRepository
}

// CreateNewGroup 新建一个群聊
func (gs GroupService) CreateNewGroup(groupname, groupintro, verificationcode string, createat time.Time, groupowner int64) error {

	var groupinfo = entity.GroupInfo{
		GroupName:  groupname,
		GroupIntro: groupintro,
		GroupOwner: groupowner,
		CreateAt:   createat,
	}
	if err := gs.GroupRepository.CreateGroup(groupinfo); err != nil {
		return err
	}
	return nil
}

// QueryGroupInfo 查询群聊信息
func (gs GroupService) QueryGroupInfo(groupid int64) (entity.GroupInfo, error) {
	groupinfo, err := gs.GroupRepository.QueryGroupInfo(groupid)

	return groupinfo, err
}

// UpdateGroupInfo 更新群聊信息
func (gs GroupService) UpdateGroupInfo(useraccount, groupid int64, groupintro string, groupname string) (bool, error) {

	var updateinfo = entity.GroupInfo{
		Groupid:    groupid,
		GroupIntro: groupintro,
		GroupName:  groupname,
	}

	if err := gs.GroupRepository.UpdateGroupInfo(updateinfo); err != nil {
		return false, err
	}
	return true, nil
}

// UpdateGroupAvatar 更新群头像
func (gs GroupService) UpdateGroupAvatar(useraccount, groupid int64, avatar string) (bool, error) {

	if err := gs.GroupRepository.UpdateGroupAvatar(groupid, avatar); err != nil {
		return false, err
	}

	return true, nil
}

// QueryGroupMember 查询群成员信息
func (gs GroupService) QueryGroupMember(groupid int64) (entity.GroupMember, error) {
	members, err := gs.GroupRepository.QueryGroupMember(groupid)

	return members, err
}

// UpdateUserNameInGroup 更新群内昵称
func (gs GroupService) UpdateUserNameInGroup(groupid, useraccount int64, usernameingroup string) (bool, error) {
	if err := gs.GroupRepository.UpdateUserNameInGroup(groupid, useraccount, usernameingroup); err != nil {
		return false, err
	}

	return true, nil
}
