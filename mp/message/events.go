// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package message

import (
	"encoding/xml"
	"strings"
)

// 模板的发送状态值
const (
	TemplateSendStatusSuccess int8 = 1
	TemplateSendStatusUserBlock
	TemplateSendStatusSystemFailed
)

// 模板的事件类型
const (
	EventTypeSubscribe             = "subscribe"
	EventTypeUnsubscribe           = "unsibuscribe"
	EventTypeScan                  = "SCAN"
	EventTypeLocation              = "LOCATION"
	EventTypeClick                 = "CLICK"
	EventTypeView                  = "VIEW"
	EventTypeTemplateSendJobFinish = "TEMPLATESENDJOBFINISH"
)

// Eventer 事件接口
type Eventer interface {
	Messager
	EventType() string
}

type event struct {
	base
	Event string `xml:"Event"` // 事件类型
}

// EventScan 表示通过扫描带参数的二维码事件
//
// subscribe 表示已关注下的扫描事件，SCAN 未关注下的扫描事件
// 若 IsScan() 为 false，则 subscribe 表示关注，unsbuscribe 表示取消关注
type EventScan struct {
	event
	EventKey string `xml:"EventKey"`
	Ticket   string `xml:"Ticket"`
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
	EventKey string `xml:"EventKey"`
}

// EventTemplateSendJobFinish 模板消息发送事件
type EventTemplateSendJobFinish struct {
	event
	MsgID  int64  `xml:"MsgID"`
	Status string `xml:"Status"`
}

func (e *event) EventType() string {
	return e.Event
}

// IsScan 是扫描产生的事件还是普通的关注事件
func (e *EventScan) IsScan() bool {
	return len(e.EventKey) > 0 || len(e.Ticket) > 0
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

// eventType 这不是一个真实存在的消息类型，
// 只是用于解析 xml 中的 Event 字段的具体值用的。
type eventType struct {
	Event string `xml:"Event"`
}

// 从指定的数据中分析其消息的类型
func getEventType(data []byte) (string, error) {
	obj := &eventType{}
	if err := xml.Unmarshal(data, obj); err != nil {
		return "", err
	}

	return obj.Event, nil
}

func getEventObj(data []byte) (Eventer, error) {
	eventType, err := getEventType(data)
	if err != nil {
		return nil, err
	}

	var obj Eventer
	switch eventType {
	case EventTypeSubscribe, EventTypeUnsubscribe, EventTypeScan:
		obj = &EventScan{}
		err = xml.Unmarshal(data, obj)
	case EventTypeLocation:
		obj = &EventLocation{}
		err = xml.Unmarshal(data, obj)
	case EventTypeView, EventTypeClick:
		obj = &EventClickView{}
		err = xml.Unmarshal(data, obj)
	case EventTypeTemplateSendJobFinish:
		obj = &EventTemplateSendJobFinish{}
		err = xml.Unmarshal(data, obj)
	}

	if err != nil {
		return nil, err
	}
	return obj, nil
}
