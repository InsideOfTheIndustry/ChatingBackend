//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: userErrorInfo.go
// description: 与用户领域相关的错误定义
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-18
//

package error

const (
	SendEmailError            ErrorCode = "SendEmailError"            // 发送邮件错误
	VerificationCodeError     ErrorCode = "VerificationCodeError"     // 验证码错误
	FieldError                ErrorCode = "FieldError"                // 字段解析失败
	ServerInternalError       ErrorCode = "ServerInternalError"       // 服务器内部错误
	RegisterError             ErrorCode = "RegisterError"             // 注册失败
	SendVerificationCodeError ErrorCode = "SendVerificationCodeError" // 验证码发送失败
	AccountOrPasswordError    ErrorCode = "AccountOrPasswordError"    // 用户名或密码验证失败
	AuthorizationError        ErrorCode = "AuthorizationError"        // 鉴权失败
)

// NewVerificationCodeError 验证码错误
func NewVerificationCodeError() Error {
	return NewCommonError(VerificationCodeError, "验证码输入错误！")
}

// NewFieldError 输入的字段有错误
func NewFieldError(errinfo string) Error {
	return NewCommonError(FieldError, "请检查输入字段！具体错误原因为：%s", errinfo)
}

// NewServerInternalError 服务器内部错误
func NewServerInternalError(errinfo string) Error {
	return NewCommonError(ServerInternalError, "服务器内部错误！具体错误原因为：%s", errinfo)
}

// NewRegisterError 注册失败错误
func NewRegisterError(errinfo string) Error {
	return NewCommonError(RegisterError, "注册失败：%s!", errinfo)
}

// NewSendVerificationCodeError 验证码发送失败
func NewSendVerificationCodeError(errinfo string) Error {
	return NewCommonError(SendVerificationCodeError, "验证码发送失败！具体原因为：%s", errinfo)
}

// NewAccountOrPasswordError 用户名或密码错误
func NewAccountOrPasswordError() Error {
	return NewCommonError(AccountOrPasswordError, "用户名或密码错误！")
}

// NewAuthorizationError 鉴权失败
func NewAuthorizationError() Error {
	return NewCommonError(AuthorizationError,"token验证失败!")
}
