//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: common.go
// description: 通用缓存实现
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-06-30
//

package common

import (
	"chatting/database/redisdatabase"
	"chatting/infrastructure/logServer"
	"time"

	"github.com/go-redis/redis/v8"
)

// CommonCacheRepository 通用存储库实现
type CommonCacheRepository struct {
	*redisdatabase.RedisEngine
}

// SetVerificationCode(emailaddr, verificationcode string) error 实现存储库接口 设置验证码
func (ccr CommonCacheRepository) SetVerificationCode(emailaddr, verificationcode string) error {
	err := ccr.Set(redisdatabase.CtxRedis, emailaddr, verificationcode, time.Duration(120*time.Second)).Err()
	if err != nil {
		logServer.Error("redis缓存设置验证码失败:%s", err.Error())
		return err
	}
	logServer.Info("设置验证码缓存成功。")
	return nil
}

// 	GetVerificationCode(emailaddr string)(string, error) 获取验证码
func (ccr CommonCacheRepository) GetVerificationCode(emailaddr string) (string, error) {
	VerificationCode, err := ccr.Get(redisdatabase.CtxRedis, emailaddr).Result()
	if err == redis.Nil {
		logServer.Error("邮箱:(%s)的验证码不存在:(%s)", emailaddr, err.Error())
		return "", nil
	} else if err != nil {
		logServer.Error("读取发现错误:(%s)", err.Error())
		return "", err
	} else {
		logServer.Info("读取邮箱：(%s)的验证码成功", emailaddr)
		return VerificationCode, nil
	}
}

// SetRequestInfo 设置请求信息
func (ccr CommonCacheRepository) SetRequestInfo(path string, verifycode string, fretime int64) error {

	err := ccr.Set(redisdatabase.CtxRedis, path, verifycode, time.Duration(time.Duration(fretime)*time.Second)).Err()
	if err != nil {
		logServer.Error("写入失败：%s", err.Error())
		return err
	}
	return nil
}

// GetRequestInfo 获取请求信息
func (ccr CommonCacheRepository) GetRequestInfo(path string) (bool, error) {
	_, err := ccr.Get(redisdatabase.CtxRedis, path).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		logServer.Error("读取发现错误:(%s)", err.Error())
		return false, err
	} else {
		return true, nil
	}
}
