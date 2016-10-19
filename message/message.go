// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"encoding/xml"
	"io"

	"github.com/issue9/wechat/common/result"
)

// 消息类型
const (
	TypeText                    = "text"
	TypeImage                   = "image"
	TypeVoice                   = "voice"
	TypeShortVideo              = "shortvideo"
	TypeLocation                = "location"
	TypeLink                    = "link"
	TypeEvent                   = "event"
	TypeTransferCustomerService = "transfer_customer_service" // 只能用于回复消息中
)

// MsgText 文本消息
type Text struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType,cdata"`      // 消息类型
	Content      string `xml:"Content,cdata"`      // 文本消息内容
	MsgID        int64  `xml:"MsgId"`              // 消息 id，64 位整型
}

type Image struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType,cdata"`      // 消息类型
	MsgID        int64  `xml:"MsgId"`              // 消息 id，64 位整型
	PicUrl       string `xml:"PicUrl,cdata"`
	MediaID      string `xml:"MediaId,cdata"`
}

// msgType 这不是一个真实存在的消息类型，
// 只是用于解析 xml 中的 MsgType 字段的具体值用的。
type msgType struct {
	MsgType string `xml:"MsgType"`
}

// 从指定的数据中分析其消息的类型
func getMsgType(data []byte) (string, error) {
	obj := &msgType{}
	if err := xml.Unmarshal(data, obj); err != nil {
		return "", result.New(600)
	}

	return obj.MsgType, nil
}

func getMsgObj(r io.Reader) (interface{}, error) {
	data := make([]byte, 0, 1000)
	_, err := io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	typ, err := getMsgType(data)
	if err != nil {
		return nil, err
	}

	var obj interface{}
	switch typ {
	case TypeText:
		obj = &Text{}
	case TypeImage:
		obj = &Image{}
	}

	if err = xml.Unmarshal(data, obj); err != nil {
		return nil, err
	}
	return obj, nil
}
