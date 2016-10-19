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
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
}

func NewReplyTransferCustomerService(m Messager) *ReplyTransferCustomerService {
	return &ReplyTransferCustomerService{
		ToUserName:   m.From(),
		FromUserName: m.To(),
		CreateTime:   m.Created(),
		MsgType:      MsgTypeTransferCustomerService,
	}
}

func (t *ReplyTransferCustomerService) Bytes() ([]byte, error) {
	return xml.Marshal(t)
}
