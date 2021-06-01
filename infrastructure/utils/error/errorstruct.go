//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: errorstruct.go
// description: 自定义错误
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-18
//

package error

import "fmt"

// ErrorCode 错误码定义
type ErrorCode string

type Error interface {
	MarshalMap() map[string]interface{} // 将错误数据序列化为Map对象
}

// CommonError 通用错误
type commonError struct {
	code    ErrorCode // 错误码定义
	message string    // 错误信息
}

// MarshalMap 将错误数据序列化为Map对象
func (e commonError) MarshalMap() map[string]interface{} {
	mapValue := map[string]interface{}{
		"code":    e.code,
		"message": e.message,
	}

	return mapValue
}

//NewCommonError 新建一个错误对象
// New 创建新的自定义错误
func NewCommonError(code ErrorCode, message string, args ...interface{}) Error {
	if len(args) > 0 {
		return &commonError{code, fmt.Sprintf(message, args...)}
	}
	return &commonError{code, message}
}
