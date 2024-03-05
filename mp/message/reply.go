// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package message

import (
	"encoding/xml"

	"github.com/issue9/wechat/internal/xxml"
)

// ReplySuccess 成功返回的内容
var ReplySuccess = []byte("success")

// ReplyTransferCustomerService 转发消息
type ReplyTransferCustomerService struct {
	XMLName      xml.Name   `xml:"xml"`
	ToUserName   xxml.CData `xml:"ToUserName"`   // 开发者微信号
	FromUserName xxml.CData `xml:"FromUserName"` // 发送方帐号（一个 OpenID）
	MsgType      xxml.CData `xml:"MsgType"`      // 消息类型
	CreateTime   int64      `xml:"CreateTime"`   // 消息创建时间 （整型）
}

// NewReplyTranferCustomerService 将所有的消息进行转发
func NewReplyTranferCustomerService(m Messager) *ReplyTransferCustomerService {
	return &ReplyTransferCustomerService{
		ToUserName:   xxml.CData{Text: m.From()},
		FromUserName: xxml.CData{Text: m.To()},
		MsgType:      xxml.CData{Text: TypeTransferCustomerService},
		CreateTime:   m.Created(),
	}
}

// Bytes 返回 []byte 内容
func (t *ReplyTransferCustomerService) Bytes() ([]byte, error) {
	return xml.Marshal(t)
}

// TODO 其它的回复类型
