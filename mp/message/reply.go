// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"encoding/xml"

	"github.com/issue9/wechat/common"
)

// ReplySuccess 成功返回的内容
var ReplySuccess = []byte("success")

// ReplyTransferCustomerService 转发消息
type ReplyTransferCustomerService struct {
	XMLName      xml.Name     `xml:"xml"`
	ToUserName   common.CData `xml:"ToUserName"`   // 开发者微信号
	FromUserName common.CData `xml:"FromUserName"` // 发送方帐号（一个 OpenID）
	MsgType      common.CData `xml:"MsgType"`      // 消息类型
	CreateTime   int64        `xml:"CreateTime"`   // 消息创建时间 （整型）
}

// NewReplyTranferCustomerService 将所有的消息进行转发
func NewReplyTranferCustomerService(m Messager) *ReplyTransferCustomerService {
	return &ReplyTransferCustomerService{
		ToUserName:   common.CData{Text: m.From()},
		FromUserName: common.CData{Text: m.To()},
		MsgType:      common.CData{Text: TypeTransferCustomerService},
		CreateTime:   m.Created(),
	}
}

// Bytes 返回 []byte 内容
func (t *ReplyTransferCustomerService) Bytes() ([]byte, error) {
	return xml.Marshal(t)
}

// TODO 其它的回复类型
