// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import "strings"

const (
	TemplateSendStatusSuccess int8 = 1
	TemplateSendStatusUserBlock
	TemplateSendStatusSystemFailed
)

// Eventer 事件接口
type Eventer interface {
	Messager
	EventType() string
}

type event struct {
	base
	Event string `xml:"Event,cdata"` // 事件类型
}

// EventSubscribe 表示普通的关注和取消关注事件
//
// subscribe 表示关注，unsbuscribe 表示取消关注
type EventSubscribe struct {
	event
}

// EventScan 表示通过扫描带参数的二维码事件
//
// subscribe 表示已关注下的扫描事件，SCAN 未关注下的扫描事件
type EventScan struct {
	event
	EventKey string `xml:"EventKey,cdata"`
	Ticket   string `xml:"Ticket,cdata"`
}

// EventLocation 表示通过扫描带参数的二维码事件
type EventLocation struct {
	event
	Latitude  float64 `xml:"Latitude"`  // 纬度
	Longitude float64 `xml:"Longitude"` // 经度
	Precision float64 `xml:"Precision"` // 精度
}

// EventClickView 表示点击事件，可以菜单或是链接。
//
// 若是点击菜单，则 EventKey 表示菜单的 key，若
// 点击的是链接，则 EventKey 表示的是要点击的链接。
type EventClickView struct {
	event
	EventKey string `xml:"EventKey,cdata"`
}

// EventTemplateSendJobFinish 模板消息发送事件
type EventTemplateSendJobFinish struct {
	event
	MsgID  int64  `xml:"MsgID"`
	Status string `xml:"Status,cdata"`
}

func (e *event) EventType() string {
	return e.Event
}

// StatusType 当前事例的状态
func (e *EventTemplateSendJobFinish) StatusType() int8 {
	switch {
	case e.Status == "success":
		return TemplateSendStatusSuccess
	case strings.Index(e.Status, "user") >= 0:
		return TemplateSendStatusUserBlock
	case strings.Index(e.Status, "system") >= 0:
		return TemplateSendStatusSystemFailed
	}

	return TemplateSendStatusSuccess
}
