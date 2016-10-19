// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

type Eventer interface {
	Event() string
}

// EventSubscribe 表示普通的关注和取消关注事件
type EventSubscribe struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	EventType    string `xml:"Event,cdata"`        // 事件类型
}

// EventScan 表示通过扫描带参数的二维码事件
type EventScan struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	EventType    string `xml:"Event,cdata"`        // 事件类型
	EventKey     string `xml:"EventKey,cdata"`
	Ticket       string `xml:"Ticket,cdata"`
}

// EventLocation 表示通过扫描带参数的二维码事件
type EventLocation struct {
	ToUserName   string  `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string  `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64   `xml:"CreateTime"`         // 消息创建时间 （整型）
	EventType    string  `xml:"Event,cdata"`        // 事件类型
	Latitude     float64 `xml:"Latitude"`           // 纬度
	Longitude    float64 `xml:"Longitude"`          // 经度
	Precision    float64 `xml:"Precision"`          // 精度
}

// EventClickView 表示点击事件，可以菜单或是链接。
//
// 若是点击菜单，则 EventKey 表示菜单的 key，若
// 点击的是链接，则 EventKey 表示的是要点击的链接。
type EventClickView struct {
	ToUserName   string `xml:"ToUserName,cdata"`   // 开发者微信号
	FromUserName string `xml:"FromUserName,cdata"` // 发送方帐号（一个 OpenID）
	CreateTime   int64  `xml:"CreateTime"`         // 消息创建时间 （整型）
	EventType    string `xml:"Event,cdata"`        // 事件类型
	EventKey     string `xml:"EventKey,cdata"`
}

// Event 表示事件类型，subscribe 表示关注，unsbuscribe 表示取消关注
func (e *EventSubscribe) Event() string {
	return e.EventType
}

// Event 表示事件类型，subscribe 表示已关注下的扫描事件，SCAN 未关注下的扫描事件
func (e *EventScan) Event() string {
	return e.EventType
}

func (e *EventLocation) Event() string {
	return "LOCATION"
}

// Event 表示事件类型，可以是 CLICK 或是 VIEW
func (e *EventClickView) Event() string {
	return e.EventKey
}
