// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// result 包含了对微信接口所有错误返回信息的定义
package result

// Result 描述了大部分微信接口的返回类型。
//
// NOTE: 当将 Result 当作一个 error 时，请确保 Result.IsOK() 为 false。
type Result struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

// New 声明一个 Result 实例
func New(code int) *Result {
	message, found := messages[code]
	if !found {
		return &Result{
			Code:    -2,
			Message: "不存在该错误代码",
		}
	}

	return &Result{
		Code:    code,
		Message: message,
	}
}

// IsOK 该结果是否正常返回。
func (r *Result) IsOK() bool {
	return r.Code == 0
}

// Error 实现 error 接口内容，在 IsOK 为 false 时，可以将 Result 当作一个 error 实例来使用。
func (r *Result) Error() string {
	return r.Message
}
