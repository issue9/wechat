// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package notify

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"

	"github.com/issue9/wechat/internal"
	"github.com/issue9/wechat/pay"
)

// Response 向微信返馈的信息
type Response struct {
	XMLName xml.Name       `xml:"xml"`
	Code    internal.CData `xml:"return_code"`
	Message internal.CData `xml:"return_msg"`
}

// Success 构建一个表示正常的 Response 实例
func Success() *Response {
	return &Response{
		Code:    internal.CData{Text: pay.Success},
		Message: internal.CData{Text: "OK"},
	}
}

// Fail 构建一个表示出错的 Response 实例，message 为出错信息
func Fail(message string) *Response {
	return &Response{
		Code:    internal.CData{Text: pay.Fail},
		Message: internal.CData{Text: message},
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
	w.Header().Set("Content-Type", "application/xml")

	n, err := r.writeTo(w)
	if err != nil {
		return err
	}
	if n <= 0 {
		return errors.New("写入的内容大小不正确")
	}

	return nil
}
