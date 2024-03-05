// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

// result 包含了对微信接口所有错误返回信息的定义
package result

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Result 描述了微信和 http 状态码，其中
// [0,400) 区间的为非错误代码；
// [400,600) 之间的为 HTTP 错误代码；
// [600,40000) 之间的自定义的错误代码；
// [40000,) 为微信的错误代码
// NOTE: 当将 Result 当作一个 error 时，请确保 Result.IsOK() 为 false。
type Result struct {
	Code    int    `json:"errcode,omitempty"`
	Message string `json:"errmsg,omitempty"`
}

// New 声明一个 Result 实例
func New(code int) *Result {
	var message string
	if code >= 200 && code < 600 { // HTTP 信息
		message = http.StatusText(code)
	} else {
		message = messages[code]

	}

	if len(message) == 0 {
		message = "不存在该错误代码" + strconv.Itoa(code)
		code = 601
	}

	return &Result{
		Code:    code,
		Message: message,
	}
}

// From 将一段字符串转换成 Result 实例。
func From(data []byte) *Result {
	r := &Result{}
	if err := json.Unmarshal(data, r); err != nil {
		r.Code = 600
		r.Message = err.Error()
	}

	return r
}

// IsOK 该结果是否正常返回。
func (r *Result) IsOK() bool {
	return r.Code >= 0 && r.Code < 400
}

// Error 实现 error 接口内容，在 IsOK 为 false 时，可以将 Result 当作一个 error 实例来使用。
func (r *Result) Error() string {
	return r.Message
}
