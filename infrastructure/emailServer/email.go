//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: email.go
// description: 邮件服务
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package emailServer

import (
	"chatting/infrastructure/configServer"
	"chatting/infrastructure/logServer"
	"crypto/tls"
	"net/smtp"
	"strings"
)

// 邮件服务引擎
type EmailEngine struct {
	User     string
	Host     string
	Password string
}

var EmailEngineModel *EmailEngine

// InitEmailEngine 初始化邮件服务引擎
func InitEmailEngine() {
	emailConfig := configServer.Emailcfg
	var emailEngine = EmailEngine{
		User:     emailConfig.User,
		Host:     emailConfig.Host,
		Password: emailConfig.Password,
	}
	EmailEngineModel = &emailEngine
}

// SendEmail 发送邮件信息
func (en *EmailEngine) SendEmail(message, subject, receiver string) error {
	host := strings.Split(en.Host, ":")
	auth := smtp.PlainAuth("", en.User, en.Password, host[0])

	msg := []byte("To: " + receiver + "\r\nFrom: Webchatting<" + en.User + ">\r\nSubject: " + subject + "\r\n" + "Content-Type: text/html; charset=UTF-8" + "\r\n\r\n" + "<html><body><h3>" + message + "</h3></body>")
	// err := smtp.SendMail(en.Host, auth, en.User, []string{receiver}, []byte(msg))

	// 阿里云禁用了 25 端口 因此需要使用ssl 465端口
	err := en.SendEmailUseTsl(auth, msg, receiver)
	if err != nil {
		logServer.Error("邮件发送至(%s)失败(%s)", receiver, err.Error())
		return err
	}
	logServer.Info("邮件发送至（%s）成功。", receiver)
	return nil
}

func (en *EmailEngine) SendEmailUseTsl(auth smtp.Auth, message []byte, receiver string) error {
	conn, err := tls.Dial("tcp", en.Host, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		logServer.Error("连接到邮箱服务器失败:%s...", err.Error())
		return err
	}
	host := strings.Split(en.Host, ":")
	smtpclient, err := smtp.NewClient(conn, host[0])
	if err != nil {
		logServer.Error("创建新的邮箱客户端失败:%s...", err.Error())
		return err
	}
	defer smtpclient.Close()

	if ok, _ := smtpclient.Extension("AUTH"); ok {
		if err = smtpclient.Auth(auth); err != nil {
			logServer.Error("Error during AUTH", err)
			return err
		}
	}
	// 设置发送者
	if err := smtpclient.Mail(en.User); err != nil {
		logServer.Error("设置邮件发送者出现错误:%s...", err.Error())
		return err
	}
	// 设置接受者
	if err := smtpclient.Rcpt(receiver); err != nil {
		logServer.Error("设置邮件接受者出现错误:%s...", err.Error())
		return err
	}

	// 获取写入句柄
	writer, err := smtpclient.Data()
	if err != nil {
		logServer.Error("初始化写入失败:%s...", err.Error())
		return err
	}

	// 写入信息
	if _, err := writer.Write(message); err != nil {
		logServer.Error("写入邮件信息失败:%s...", err.Error())
		return err
	}

	if err := writer.Close(); err != nil {
		logServer.Error("邮箱写入句柄关闭失败:%s...", err.Error())
		return err
	}

	return smtpclient.Quit()

}
