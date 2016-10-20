// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"errors"
	"fmt"
)

// Handler 消息处理函数。
//通过向 NewServer 注册 Handler 函数，获取对消息处理的权限。
//
// 参数 Messager 为从微信端传递过来的 xml 数据对象实例，都已定义在 messages.go
// 文件中。
//
// 函数的返回值，被当作消息被动回复内容传递给微信调用方。在 reply.go 中
// 定义了大部分可能用到返回类型，可以拿来直接使用。
//
// NOTE 所有的 Handler 必须在 5 秒内有返回数据，否则微信端会再次发起同样的请求
type Handler func(Messager) ([]byte, error)

// HandlerBus 为 Handler 接口的管理器，方便用户按类别来注册消息处理。
//
//  h := &HandlerBus{}
//  h.RegisterMessage(TypeText, h1)
//  h.RegisterMessage(TypeImage, h2)
//  srv := NewServer("token", h.Handler, nil)
type HandlerBus struct {
	messageHandlers map[string]Handler
	eventHandlers   map[string]Handler
}

// NewHandlerBus 声明一个新的 HandlerBus。
func NewHandlerBus() *HandlerBus {
	return &HandlerBus{
		messageHandlers: make(map[string]Handler, 9),
		eventHandlers:   make(map[string]Handler, 7),
	}
}

// RegisterMessage 注册消息处理函数。
// typ 的值若为 event，可以注册，但不会实际有作用。
func (b *HandlerBus) RegisterMessage(typ string, h Handler) {
	b.messageHandlers[typ] = h
}

// RegisterEvent 注册事件处理函数
func (b *HandlerBus) RegisterEvent(event string, h Handler) {
	b.eventHandlers[event] = h
}

func (b *HandlerBus) Handler(m Messager) ([]byte, error) {
	var h Handler
	var found bool
	typ := m.Type()

	if typ == TypeEvent {
		event := m.(Eventer).EventType()
		h, found = b.eventHandlers[event]
		if !found {
			return nil, fmt.Errorf("事件[%v]的处理函数不存在")
		}
		return h(m)
	} else {
		h, found = b.eventHandlers[typ]
		if !found { // 消息处理函数不存在的情况下，实行转发
			h = TransferCustomerService
		}
	}

	return h(m)
}

// TransferCustomerService 是 Handler 的一种实现，实现了对消息的转发。
func TransferCustomerService(m Messager) ([]byte, error) {
	if m.Type() == TypeEvent {
		return nil, errors.New("事件不允许转发给客服")
	}

	return NewReplyTranferCustomerService(m).Bytes()
}
