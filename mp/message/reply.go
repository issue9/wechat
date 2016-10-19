// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"encoding/xml"
)

var ReplySuccess = []byte("success")

// ReplyTransferCustomerService 转发消息
type ReplyTransferCustomerService struct {
	base
}

// ReplyText 回复文本消息
type ReplyText struct {
	Text
}

// NewReplyTranferCustomerService 将所有的消息进行转发
func NewReplyTranferCustomerService(m Messager) *ReplyTransferCustomerService {
	return &ReplyTransferCustomerService{base: base{
		ToUserName:   m.From(),
		FromUserName: m.To(),
		MsgType:      m.Type(),
		CreateTime:   m.Created(),
	}}
}

func (t *ReplyTransferCustomerService) Bytes() ([]byte, error) {
	return xml.Marshal(t)
}

// NewReplyText 从文本消息中构建一个回复文本消息
func NewReplyText(t *Text) *ReplyText {
	ret := &ReplyText{}
	ret.ToUserName = t.From()
	ret.FromUserName = t.To()
	ret.MsgType = t.Type()
	ret.CreateTime = t.Created()
	ret.Content = t.Content

	return ret
}

func (t *ReplyText) Bytes() ([]byte, error) {
	return xml.Marshal(t)
}
