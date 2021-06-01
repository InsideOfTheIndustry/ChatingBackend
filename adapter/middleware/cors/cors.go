//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: cors.go
// description: 解决跨域问题
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-19
//

package cors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")                                                                   // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization,token,account") // header内允许被跨域的字段
		c.Header("Access-Control-Allow-Methods", "POST, GET,PUT,OPTIONS,DELETE")                                       // 可以方形的请求
		//c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
