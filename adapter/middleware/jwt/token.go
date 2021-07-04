//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: token.go
// description: token验证中间件
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-06-30
//

package jwt

import (
	"chatting/infrastructure/logServer"
	"strconv"
	"time"

	Service "chatting/adapter/controller"
	CommonError "chatting/infrastructure/utils/error"

	"github.com/gin-gonic/gin"
)

// TokenVerify token鉴权 其中加入了请求频率限制
func TokenVerify(c *gin.Context) {

	token := c.GetHeader("token")
	useraccount := c.GetHeader("account")
	path := c.FullPath()
	method := c.Request.Method

	frequent, err := Service.CommonService.JudgeRequestFrequence(path+method+useraccount, 5)
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

	jwtclaim, err := ParseToken(token)
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

// TokenVerifyWithoutTImeLimit 其中每页请求频率限制
func TokenVerifyWithoutTImeLimit(c *gin.Context) {

	token := c.GetHeader("token")
	useraccount := c.GetHeader("account")

	useraccountint, err := strconv.ParseInt(useraccount, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(401, CommonError.NewAuthorizationError().MarshalMap())
		return
	}

	if token == "" || useraccount == "" {
		c.AbortWithStatusJSON(401, CommonError.NewAuthorizationError().MarshalMap())
		return
	}

	jwtclaim, err := ParseToken(token)
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
