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
	TypeVideo                   = "video"
	TypeShortVideo              = "shortvideo"
	TypeLocation                = "location"
	TypeLink                    = "link"
	TypeEvent                   = "event"
	TypeTransferCustomerService = "transfer_customer_service" // 只能用于回复消息中
)

// Messager 表示消息的基本结构。
type Messager interface {
	// 获取消息的类型
	Type() string
}

// MsgText 文本消息
type Text struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType,cdata"`      // 消息类型
	MsgID        int64  `xml:"MsgId"`              // 消息 id，64 位整型
	Content      string `xml:"Content,cdata"`      // 文本消息内容
}

// Image 图片消息
type Image struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType,cdata"`      // 消息类型
	MsgID        int64  `xml:"MsgId"`              // 消息 id，64 位整型
	PicURL       string `xml:"PicUrl,cdata"`
	MediaID      string `xml:"MediaId,cdata"`
}

// Voice 语音消息
type Voice struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType,cdata"`      // 消息类型
	MsgID        int64  `xml:"MsgId"`              // 消息 id，64 位整型
	MediaID      string `xml:"MediaId,cdata"`
	Format       string `xml:"Format,cdata"`
	Recognition  string `xml:"Recognition,cdata,omitempty"` // 语音识别结果
}

// Video 视频消息
type Video struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType,cdata"`      // 消息类型
	MsgID        int64  `xml:"MsgId"`              // 消息 id，64 位整型
	MediaID      string `xml:"MediaId,cdata"`
	ThumbMediaID string `xml:"ThumbMediaId,cdata"`
}

// shortVideo 短视频消息
type ShortVideo struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType,cdata"`      // 消息类型
	MsgID        int64  `xml:"MsgId"`              // 消息 id，64 位整型
	MediaID      string `xml:"MediaId,cdata"`
	ThumbMediaID string `xml:"ThumbMediaId,cdata"`
}

// Location 位置消息
type Location struct {
	ToUserName   string  `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string  `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64   `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string  `xml:"MsgType,cdata"`      // 消息类型
	MsgID        int64   `xml:"MsgId"`              // 消息 id，64 位整型
	X            float64 `xml:"Location_X"`         // 维度
	Y            float64 `xml:"Location_Y"`         // 经度
	Scale        int     `xml:"Scale"`
	Label        string  `xml:"Label,cdata"` // 地理位置信息
}

// Link 链接消息
type Link struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType,cdata"`      // 消息类型
	MsgID        int64  `xml:"MsgId"`              // 消息 id，64 位整型
	Title        string `xml:"Title,cdata"`
	Description  string `xml:"Description,cdata"`
	URL          string `xml:"Url,cdata"`
}

func (t *Text) Type() string {
	return TypeText
}

func (i *Image) Type() string {
	return TypeImage
}

func (v *Voice) Type() string {
	return TypeVoice
}

func (v *Video) Type() string {
	return TypeVideo
}

func (v *ShortVideo) Type() string {
	return TypeShortVideo
}

func (l *Location) Type() string {
	return TypeLocation
}

func (l *Link) Type() string {
	return TypeLink
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

func getMsgObj(r io.Reader) (Messager, error) {
	data := make([]byte, 0, 1000)
	_, err := io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	typ, err := getMsgType(data)
	if err != nil {
		return nil, err
	}

	var obj Messager
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
