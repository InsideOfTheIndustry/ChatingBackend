//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: group.go
// description: 群聊相关api
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-07-01
//

package group

import (
	Service "chatting/adapter/controller"
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/logServer"
	CommonError "chatting/infrastructure/utils/error"
	"errors"
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
)

//SendNewGroupVerificationCode 发送验证码
func SendNewGroupVerificationCode(c *gin.Context) {
	var userinfo = UserInfoGet{}

	// 解析数据

	if err := c.BindJSON(&userinfo); err != nil {
		logServer.Error("数据绑定失败:(%s)", err.Error())
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	// 设置邮件发送频率
	useraccount := strconv.FormatInt(userinfo.UserAccount, 10)
	path := c.FullPath()
	frequent, err := Service.CommonService.JudgeRequestFrequence(path+useraccount+"email", 60)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	if frequent {
		c.JSON(500, CommonError.NewRequestsTooFrequentError().MarshalMap())
		return
	}

	user, err := Service.UserService.GetUserInfo(userinfo.UserAccount)
	if err != nil {
		c.JSON(500, CommonError.NewSendVerificationCodeError(err.Error()).MarshalMap())
		return
	}

	if user.OwnGroups >= 2 {
		c.JSON(500, CommonError.NewSendVerificationCodeError(errors.New("number of the group you created allready more than 2, you can not create new group").Error()).MarshalMap())
		return
	}
	// 生成验证码
	var messagecode = ""
	for i := 0; i < 4; i++ {
		number := rand.Intn(10)
		word := strconv.Itoa(number)
		messagecode += word
	}

	// 验证码缓存
	if _, err := Service.CommonService.SetVerificationCode(user.UserEmail+"group", messagecode); err != nil {
		c.JSON(500, CommonError.NewSendVerificationCodeError(err.Error()).MarshalMap())
		return
	}

	// 发送验证码
	if err := Service.CommonService.SendEmail(messagecode, "webchatting群聊创建验证码", user.UserEmail); err != nil {
		c.JSON(500, CommonError.NewSendVerificationCodeError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "验证码发送成功",
		"status":  200,
	})
}

// CreateNewGroup 新建群聊
func CreateNewGroup(c *gin.Context) {
	var groupinfo = GroupInfoCreate{}
	if err := c.BindJSON(&groupinfo); err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}
	// 获取信息判断还能否创建
	user, err := Service.UserService.GetUserInfo(groupinfo.UserAccount)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	if user.OwnGroups >= 2 {
		c.JSON(500, CommonError.NewSendVerificationCodeError(errors.New("number of the group you created allready more than 2, you can not create new group").Error()).MarshalMap())
		return
	}
	// 验证码是否相等
	vc, err := Service.CommonService.GetVerificationCode(user.UserEmail + "group")
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	if vc != groupinfo.VerificationCode {
		c.JSON(400, CommonError.NewServerInternalError(errors.New("验证码错误或失效!").Error()).MarshalMap())
		return
	}

	if err := Service.GroupService.CreateNewGroup(groupinfo.GroupName, groupinfo.GroupIntro, groupinfo.VerificationCode, groupinfo.CreateAt, groupinfo.UserAccount); err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "新建群聊成功！",
		"success": true,
	})
}

// QueryGroupInfo 查询群信息
func QueryGroupInfo(c *gin.Context) {
	var groupid = c.Param("groupid")

	groupidint, err := strconv.ParseInt(groupid, 10, 64)
	if err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	groupinfo, err := Service.GroupService.QueryGroupInfo(groupidint)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message":   "查询群聊信息成功！",
		"groupinfo": groupinfo,
	})
}

// QueryGroupMember 查询群成员信息
func QueryGroupMember(c *gin.Context) {
	var groupid = c.Param("groupid")

	groupidint, err := strconv.ParseInt(groupid, 10, 64)
	if err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	memberinfo, err := Service.GroupService.QueryGroupMember(groupidint)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message":      "查询群成员成功！",
		"groupmembers": memberinfo,
	})
}

// UpdateGroupInfo 更新群信息
func UpdateGroupInfo(c *gin.Context) {
	var groupinfoupdate = GroupInfoUpdate{}
	if err := c.BindJSON(&groupinfoupdate); err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	groupinfos, err := Service.GroupService.QueryGroupInfo(groupinfoupdate.Groupid)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	if groupinfos.GroupOwner != groupinfoupdate.UserAccount {
		c.JSON(401, CommonError.NewAuthorizationError().MarshalMap())
		return
	}

	success, err := Service.GroupService.UpdateGroupInfo(groupinfoupdate.UserAccount, groupinfoupdate.Groupid, groupinfoupdate.GroupIntro, groupinfoupdate.GroupName)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()))
		return
	}

	if !success {
		c.JSON(401, CommonError.NewAuthorizationError())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "更新群信息成功！",
	})

}

// UpdateGroupAvatar 更新群聊头像
func UpdateGroupAvatar(c *gin.Context) {
	var avatarupdateinfo = GroupInfoUpdateAvatar{}

	if err := c.BindJSON(&avatarupdateinfo); err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	groupinfo, err := Service.GroupService.QueryGroupInfo(avatarupdateinfo.Groupid)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	if groupinfo.GroupOwner != avatarupdateinfo.UserAccount {
		c.JSON(401, CommonError.NewAuthorizationError().MarshalMap())
		return
	}

	groupid := strconv.FormatInt(avatarupdateinfo.Groupid, 10)
	avartarremotepath, err := Service.CommonService.EncodeImageAndStore("./resourcelocation/groupavatar/", avatarupdateinfo.Avatar, groupid, "", configServer.Applicationcfg.ServerIp, configServer.Applicationcfg.Port, "/groupavatar/")
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	if _, err := Service.GroupService.UpdateGroupAvatar(avatarupdateinfo.UserAccount, avatarupdateinfo.Groupid, avartarremotepath); err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "更新头像成功！",
	})
}

// UpdateUserNameInGroup 更新群内名称
func UpdateUserNameInGroup(c *gin.Context) {
	var usernameupdate = UserNameInGroupUpdateInfo{}

	if err := c.BindJSON(&usernameupdate); err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	if _, err := Service.GroupService.UpdateUserNameInGroup(usernameupdate.Groupid, usernameupdate.UserAccount, usernameupdate.UserNameInGroup); err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "更新昵称成功！",
	})
}
