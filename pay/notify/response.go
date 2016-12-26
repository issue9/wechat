// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package notify

import (
	"encoding/xml"

	"github.com/issue9/wechat/pay"
)

// Response 向微信返馈的信息
type Response struct {
	XMLName xml.Name `xml:"xml"`
	Code    string   `xml:"return_code"`
	Message string   `xml:"return_msg"`
}

// Success 构建一个表示正常的 Response 实例
func Success() *Response {
	return &Response{
		Code:    pay.Success,
		Message: "OK",
	}
}

// Fail 构建一个表示出错的 Response 实例，message 为出错信息
func Fail(message string) *Response {
	return &Response{
		Code:    pay.Fail,
		Message: message,
	}
}
