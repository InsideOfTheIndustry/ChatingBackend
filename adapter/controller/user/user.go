//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: user.go
// description: 用户相关的控制器
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-16
//

package user

import (
	Service "chatting/adapter/controller"
	Jwt "chatting/adapter/middleware/jwt"
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/logServer"
	CommonError "chatting/infrastructure/utils/error"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterAccount 用户注册
func RegisterAccount(c *gin.Context) {
	var userinfo = User{}

	// 解析数据
	err := c.BindJSON(&userinfo)
	if err != nil {
		logServer.Error("用户信息格式有误:(%s)", err.Error())
		c.JSON(400, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	// 验证验证码与邮箱
	ifhasemail, err := Service.UserService.VerifyCode(userinfo.VerificationCode, userinfo.UserEmail)

	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	if !ifhasemail {
		c.JSON(400, CommonError.NewVerificationCodeError().MarshalMap())
		return
	}

	// 正式注册
	useraccount, err := Service.UserService.RegisterUser(userinfo.UserEmail, userinfo.UserName, userinfo.UserPassword, userinfo.Signature, userinfo.Avatar, userinfo.UserAge, userinfo.UserSex)
	if err != nil {

		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	if useraccount == 0 {

		c.JSON(500, CommonError.NewRegisterError(errors.New("邮箱已被注册。").Error()).MarshalMap())
		return
	}

	// 注册成功发送邮件通知
	err = Service.UserService.SendUseraccount(useraccount, userinfo.UserEmail)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message":     "账号注册成功，请查看您的邮箱获取账号信息。",
		"useraccount": useraccount,
	})

}

//SendVerificationCode 发送验证码
func SendVerificationCode(c *gin.Context) {
	var emailinfo = VerificationCode{}

	// 解析数据
	if err := c.BindJSON(&emailinfo); err != nil {
		logServer.Error("数据绑定失败:(%s)", err.Error())
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}
	path := c.FullPath()
	frequent, err := Service.UserService.JudgeRequestFrequence(path+emailinfo.UserEmail, 60)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	if frequent {
		c.JSON(500, CommonError.NewRequestsTooFrequentError().MarshalMap())
		return
	}

	if err := Service.UserService.SendVerificationCode(emailinfo.UserEmail); err != nil {
		c.JSON(500, CommonError.NewSendVerificationCodeError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "验证码发送成功",
		"status":  200,
	})
}

//UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var userinfo = LoginInfo{}
	err := c.BindJSON(&userinfo)
	if err != nil {
		logServer.Error("数据绑定失败！")
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	// 验证用户信息
	loginsuccess, err := Service.UserService.VerifyLoginInfo(userinfo.UserAccount, userinfo.UserPassword)

	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()))
		return
	}
	if !loginsuccess {
		c.JSON(401, CommonError.NewAccountOrPasswordError().MarshalMap())
		return
	}

	// 生成token
	token, err := Jwt.GenarateToken(userinfo.UserAccount)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	// 保存token到redis
	err = Service.UserService.UserCacheRepository.SetToken(userinfo.UserAccount, token)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "登录成功！",
		"token":   token,
	})
}

//GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	var userinfoget = UserInfoGet{}
	err := c.BindJSON(&userinfoget)
	if err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}
	useraccountint := userinfoget.UserAccount

	userinfo, err := Service.UserService.UserRepository.Query(useraccountint)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	userinfo.UserPassword = "******"
	c.JSON(200, map[string]interface{}{
		"message":  "查询成功！",
		"userinfo": userinfo,
	})
}

// GetFriendInfo 获取用户好友信息
func GetFriendInfo(c *gin.Context) {
	var userinfoget = UserInfoGet{}

	if err := c.BindJSON(&userinfoget); err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}
	useraccountint := userinfoget.UserAccount

	friendInfo, err := Service.UserService.GetUserFriendInfo(useraccountint)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "查询用户好友信息成功！",
		"friends": friendInfo.Friends,
	})

}

// UpdateUserAvatar 更新用户头像信息
func UpdateUserAvatar(c *gin.Context) {

	var updateinfo = UserInfoUpdateAvatar{}

	if err := c.BindJSON(&updateinfo); err != nil {
		logServer.Error("错误为:%s", err.Error())
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	// 获取图片类型
	avartarslice := strings.Split(updateinfo.Avatar, ";")
	// 11-last: data:image/type;xxxx...
	avatartype := avartarslice[0][11:]

	// 对图片进行转码
	avatarbyte, err := base64.StdEncoding.DecodeString(avartarslice[1][7:])
	if err != nil {
		logServer.Error("错误为:%s", err.Error())
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	// useraccount 转为 string
	useraccount := strconv.FormatInt(updateinfo.UserAccount, 10)

	avartarfile, err := os.OpenFile("./resourcelocation/useravatar/"+useraccount+"."+avatartype, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	defer avartarfile.Close()

	if _, err := avartarfile.Write(avatarbyte); err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	var avatarremotepath = "http://" + configServer.Applicationcfg.ServerIp + ":" + configServer.Applicationcfg.Port + "/useravatar/" + useraccount + "." + avatartype

	if _, err := Service.UserService.UpdateUserAvatar(updateinfo.UserAccount, avatarremotepath); err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"avatar":  avatarremotepath,
		"message": "success",
	})

}

// TokenVerify token鉴权 其中加入了请求频率限制
func TokenVerify(c *gin.Context) {

	token := c.GetHeader("token")
	useraccount := c.GetHeader("account")
	path := c.FullPath()

	frequent, err := Service.UserService.JudgeRequestFrequence(path+useraccount, 5)
	if err != nil {
		c.AbortWithStatusJSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	if frequent {
		c.AbortWithStatusJSON(500, CommonError.NewRequestsTooFrequentError().MarshalMap())
		return
	}

	useraccountint, err := strconv.ParseInt(useraccount, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(401, CommonError.NewAuthorizationError().MarshalMap())
		return
	}

	if token == "" || useraccount == "" {
		c.AbortWithStatusJSON(401, CommonError.NewAuthorizationError().MarshalMap())
		return
	}

	jwtclaim, err := Jwt.ParseToken(token)
	if err != nil {
		logServer.Error("token解析失败:(%s)", err.Error())
		c.AbortWithStatusJSON(401, CommonError.NewAuthorizationError().MarshalMap())
		return
	}

	if jwtclaim.UserAccount != useraccountint || jwtclaim.ExpiresAt < time.Now().Unix() {
		c.AbortWithStatusJSON(401, CommonError.NewAuthorizationError().MarshalMap())
		return
	}
	c.Next()
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	var userupdateinfo = UserInfoUpdate{}
	if err := c.BindJSON(&userupdateinfo); err != nil {
		logServer.Error("绑定数据出现错误:%s", err.Error())
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	if err := Service.UserService.UpdateUserInfo(userupdateinfo.UserAccount, userupdateinfo.UserAge, userupdateinfo.UserSex, userupdateinfo.UserName, userupdateinfo.Signature); err != nil {
		logServer.Error("更新用户信息失败出现错误:%s", err.Error())
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	userinfo, err := Service.UserService.UserRepository.Query(userupdateinfo.UserAccount)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	userinfo.UserPassword = "******"
	c.JSON(200, map[string]interface{}{
		"message":  "用户信息更新成功！",
		"userinfo": userinfo,
	})
}

// GetUserGroupInfo 查询用户所在的群信息
func GetUserGroupInfo(c *gin.Context) {
	var userinfoget = UserInfoGet{}
	if err := c.BindJSON(&userinfoget); err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}
	groupsinfo, err := Service.UserService.UserRepository.QueryGroupOfUser(userinfoget.UserAccount)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message":   "查询用户群聊信息成功！",
		"groupinfo": groupsinfo,
	})
}

// CreateNewGroup 新建群聊
func CreateNewGroup(c *gin.Context) {
	var groupinfo = GroupInfoCreateAndUpdate{}
	if err := c.BindJSON(&groupinfo); err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}
	if err := Service.UserService.CreateNewGroup(groupinfo.GroupName, groupinfo.GroupIntro, groupinfo.VerificationCode, groupinfo.CreateAt, groupinfo.UserAccount); err != nil {
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
	var groupinfoget = GroupInfoGet{}

	if err := c.BindJSON(&groupinfoget); err != nil {
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	groupinfo, err := Service.UserService.UserRepository.QueryGroupInfo(groupinfoget.Groupid)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message":   "查询群聊信息成功！",
		"groupinfo": groupinfo,
	})
}

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
	frequent, err := Service.UserService.JudgeRequestFrequence(path+useraccount+"email", 60)
	if err != nil {
		c.JSON(500, CommonError.NewServerInternalError(err.Error()).MarshalMap())
		return
	}
	if frequent {
		c.JSON(500, CommonError.NewRequestsTooFrequentError().MarshalMap())
		return
	}

	if err := Service.UserService.SendCreateNewGroupVerifyCode(userinfo.UserAccount); err != nil {
		c.JSON(500, CommonError.NewSendVerificationCodeError(err.Error()).MarshalMap())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "验证码发送成功",
		"status":  200,
	})
}
