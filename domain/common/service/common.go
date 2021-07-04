//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: common.go
// description: 通用服务
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-06-30
//

package service

import (
	"chatting/domain/common/repository"
	"chatting/infrastructure/logServer"
	"encoding/base64"
	"os"
	"strings"
)

// CommonService 通用服务
type CommonService struct {
	CommonRepository      repository.CommonRepository
	CommonCacheRepository repository.CommonCacheRepository
	CommonEmailRepository repository.CommonEmailRepository
}

// SetVerificationCode 设置通用验证码缓存
func (cs CommonService) SetVerificationCode(setkey, setvalue string) (bool, error) {
	if err := cs.CommonCacheRepository.SetVerificationCode(setkey, setvalue); err != nil {
		return false, err
	}

	return true, nil
}

// GetVerificationCode 获取通用验证码
func (cs CommonService) GetVerificationCode(setkey string) (string, error) {
	code, err := cs.CommonCacheRepository.GetVerificationCode(setkey)
	if err != nil {
		return "", err
	}

	return code, err
}

// JudgeRequestFrequence 判断请求访问频率
func (us CommonService) JudgeRequestFrequence(path string, fretime int64) (bool, error) {
	ifexist, err := us.CommonCacheRepository.GetRequestInfo(path)

	// 出现错误返回错误
	if err != nil {
		return false, err
	}
	// 表示请求访问过于频繁
	if ifexist {
		return true, nil
	}
	// 合理时间内的访问
	if err := us.CommonCacheRepository.SetRequestInfo(path, "5", fretime); err != nil {
		return false, err
	}

	return false, err

}

// SendEmail  发送邮件
func (cs CommonService) SendEmail(message, subject, receiver string) error {

	if err := cs.CommonEmailRepository.SendEmail(message, subject, receiver); err != nil {
		return err
	}

	return nil
}

// EncodeImageAndStore 编码图片并存储
func (cs CommonService) EncodeImageAndStore(storagepath, avartarencode, markname, randomcode, serverip, port, remoteloc string) (string, error) {
	// 获取图片类型
	avartarslice := strings.Split(avartarencode, ";")
	// 11-last: data:image/type;xxxx...
	avatartype := avartarslice[0][11:]

	// 对图片进行转码
	avatarbyte, err := base64.StdEncoding.DecodeString(avartarslice[1][7:])
	if err != nil {
		logServer.Error("图片解码失败:%s", err.Error())
		return "", err
	}

	avartarfile, err := os.OpenFile(storagepath+markname+randomcode+"."+avatartype, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		logServer.Error("打开文件错误：%s", err.Error())
		return "", err
	}
	defer avartarfile.Close()

	if _, err := avartarfile.Write(avatarbyte); err != nil {
		logServer.Error("写入数据错误：%s", err.Error())
		return "", err
	}
	var avatarremotepath = "http://" + serverip + ":" + port + remoteloc + markname + randomcode + "." + avatartype
	return avatarremotepath, nil
}
