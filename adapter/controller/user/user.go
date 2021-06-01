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
	"chatting/domain/user/entity"
	"chatting/infrastructure/logServer"
	CommonError "chatting/infrastructure/utils/error"
	"errors"
	"strconv"
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

	var registeruserinfo = entity.UserInfo{
		UserEmail:    userinfo.UserEmail,
		UserName:     userinfo.UserName,
		UserPassword: userinfo.UserPassword,
		Signature:    userinfo.Signature,
		UserAge:      userinfo.UserAge,
		UserSex:      userinfo.UserSex,
		Avatar:       userinfo.Avatar,
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
	useraccount, err := Service.UserService.RegisterUser(registeruserinfo)
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
	err := c.BindJSON(&emailinfo)
	if err != nil {
		logServer.Error("数据绑定失败:(%s)", err.Error())
		c.JSON(500, CommonError.NewFieldError(err.Error()).MarshalMap())
		return
	}

	err = Service.UserService.SendVerificationCode(emailinfo.UserEmail)
	if err != nil {
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
	err := c.BindJSON(&userinfoget)
	if err != nil {
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
		"message":"查询用户好友信息成功！",
		"friends": friendInfo.Friends,
	})


}

// TokenVerify token鉴权
func TokenVerify(c *gin.Context) {
	token := c.GetHeader("token")
	useraccount := c.GetHeader("account")
	logServer.Info("token:%s,useraccount %s", token, useraccount)
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
