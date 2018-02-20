// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package notify

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"

	"github.com/issue9/wechat/common"
	"github.com/issue9/wechat/pay"
)

// Response 向微信返馈的信息
type Response struct {
	XMLName xml.Name     `xml:"xml"`
	Code    common.CData `xml:"return_code"`
	Message common.CData `xml:"return_msg"`
}

// Success 构建一个表示正常的 Response 实例
func Success() *Response {
	return &Response{
		Code:    common.CData{Text: pay.Success},
		Message: common.CData{Text: "OK"},
	}
}

// Fail 构建一个表示出错的 Response 实例，message 为出错信息
func Fail(message string) *Response {
	return &Response{
		Code:    common.CData{Text: pay.Fail},
		Message: common.CData{Text: message},
	}
}

func (r *Response) writeTo(w io.Writer) (int, error) {
	bs, err := xml.Marshal(r)
	if err != nil {
		return 0, err
	}

	return w.Write(bs)
}

// Render 输出到客户端
func (r *Response) Render(state int, w http.ResponseWriter) error {
	w.WriteHeader(state)
	w.Header().Set("ContentType", "application/xml")

	n, err := r.writeTo(w)
	if err != nil {
		return err
	}
	if n <= 0 {
		return errors.New("写入的内容大小不正确")
	}

	return nil
}
