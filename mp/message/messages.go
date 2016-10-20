// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import "encoding/xml"

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

// Messager 表示消息和事件的基本结构。
type Messager interface {
	// 消息类型，对应 MsgType 字段
	Type() string

	// 开发者微信号，对应 ToUserName 字段
	To() string

	// 发送方账号，对应 FromUserName 字段
	From() string

	// 创建时间，对应 CreateTime 字段
	Created() int64
}

// Message 表示消息的基本结构，不包含事件
type Message interface {
	Messager

	// 表示消息的 ID
	ID() int64
}

// 所有消息的基本内容，包含事件
type base struct {
	ToUserName   string `xml:"ToUserName"`   // 开发者微信号
	FromUserName string `xml:"FromUserName"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`   // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType"`      // 消息类型
}

type message struct {
	base
	MsgID int64 `xml:"MsgId"` // 消息 id，64 位整型
}

// Text 文本消息
type Text struct {
	message
	Content string `xml:"Content"` // 文本消息内容
}

// Image 图片消息
type Image struct {
	message
	PicURL  string `xml:"PicUrl"`
	MediaID string `xml:"MediaId"`
}

// Voice 语音消息
type Voice struct {
	message
	MediaID     string `xml:"MediaId"`
	Format      string `xml:"Format"`
	Recognition string `xml:"Recognition"` // 语音识别结果
}

// Video 视频消息
type Video struct {
	message
	MediaID      string `xml:"MediaId"`
	ThumbMediaID string `xml:"ThumbMediaId"`
}

// shortVideo 短视频消息
type ShortVideo struct {
	message
	MediaID      string `xml:"MediaId"`
	ThumbMediaID string `xml:"ThumbMediaId"`
}

// Location 位置消息
type Location struct {
	message
	X     float64 `xml:"Location_X"` // 维度
	Y     float64 `xml:"Location_Y"` // 经度
	Scale int     `xml:"Scale"`
	Label string  `xml:"Label"` // 地理位置信息
}

// Link 链接消息
type Link struct {
	message
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	URL         string `xml:"Url"`
}

func (b *base) To() string {
	return b.ToUserName
}

func (b *base) From() string {
	return b.FromUserName
}

func (b *base) Created() int64 {
	return b.CreateTime
}

func (b *base) Type() string {
	return b.MsgType
}

func (m *message) ID() int64 {
	return m.MsgID
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
		return "", err
	}

	return obj.MsgType, nil
}

// 根据类型或事件，获取相应的初始化对象
func getMessageObj(data []byte) (Messager, error) {
	typ, err := getMsgType(data)
	if err != nil {
		return nil, err
	}

	var obj Messager
	switch typ {
	case TypeText:
		obj = &Text{}
		err = xml.Unmarshal(data, obj)
	case TypeImage:
		obj = &Image{}
		err = xml.Unmarshal(data, obj)
	case TypeVoice:
		obj = &Voice{}
		err = xml.Unmarshal(data, obj)
	case TypeVideo:
		obj = &Video{}
		err = xml.Unmarshal(data, obj)
	case TypeShortVideo:
		obj = &ShortVideo{}
		err = xml.Unmarshal(data, obj)
	case TypeLocation:
		obj = &Location{}
		err = xml.Unmarshal(data, obj)
	case TypeLink:
		obj = &Link{}
		err = xml.Unmarshal(data, obj)
	case TypeEvent:
		return getEventObj(data)
	}

	if err != nil {
		return nil, err
	}
	return obj, nil
}
