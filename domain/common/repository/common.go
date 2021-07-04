//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: common.go
// description: 通用服务仓储库
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-06-30
//

package repository

// CommonRepository 通用存储接口
type CommonRepository interface {
}

// CommonCacheRepository 通用存储缓存接口
type CommonCacheRepository interface {
	SetVerificationCode(setkey, setvalue string) error                  // 设置通用验证缓存信息
	GetVerificationCode(setkey string) (string, error)                  // 获取通用验证码
	GetRequestInfo(path string) (bool, error)                           // 获取访问信息（防止连续多次操作）
	SetRequestInfo(path string, verifycode string, fretime int64) error // 设置访问信息（防止连续多次操作）
}

// CommonEmailRepository 用户邮箱服务接口
type CommonEmailRepository interface {
	SendEmail(message, subject, receiver string) error // 发送邮件
}
