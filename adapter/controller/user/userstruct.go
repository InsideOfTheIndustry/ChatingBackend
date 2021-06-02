//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: userstruct.go
// description: 用户相关结构体定义
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-16
//

package user

// User新建用户时传入的用户数据
type User struct {
	UserAccount      int64  `json:"useraccount"`      // 用户账号
	UserEmail        string `json:"useremail"`        // 邮箱号
	UserName         string `json:"username"`         // 用户名
	Signature        string `json:"signature"`        // 用户个性签名
	Avatar           string `json:"avatar"`           // 用户头像
	UserPassword     string `json:"userpassword"`     // 用户密码
	UserAge          int64  `json:"userage"`          // 用户年龄
	UserSex          int64  `json:"usersex"`          // 用户性别
	VerificationCode string `json:"verificationcode"` // 验证码
}

// VerificationCode 发送验证码时接收的参数信息
type VerificationCode struct {
	UserEmail string `json:"useremail"` // 邮箱号
}

// LoginInfo 用户登录信息
type LoginInfo struct {
	UserAccount  int64  `json:"useraccount"`  // 用户账号
	UserPassword string `json:"userpassword"` // 用户密码
}

// UserInfoGet 获取用户信息时的输入结构
type UserInfoGet struct {
	UserAccount int64 `json:"useraccount"` // 用户账号
}

// UserInfoUpadateAvatar 用户上传头像时需要的信息
type UserInfoUpdateAvatar struct {
	UserAccount int64  `json:"useraccount"` // 用户账号
	Avatar      string `json:"avatar"`      // 用户上传的头像信息：图片经过base64编码后的内容
}
