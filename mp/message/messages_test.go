// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package message

import (
	"testing"

	"github.com/issue9/assert/v4"
)

var _ Message = &Text{}
var _ Message = &Image{}
var _ Message = &Voice{}
var _ Message = &Video{}
var _ Message = &ShortVideo{}
var _ Message = &Location{}
var _ Message = &Link{}
var _ Message = &message{}

func TestGetMsgType(t *testing.T) {
	a := assert.New(t, false)

	data := []byte(`<xml>
	<ToUserName><![CDATA[12d]]></ToUserName>
	<FromUserName><![CDATA[dddadfaee]]></FromUserName>
	<CreateTime>12345555</CreateTime>
	<MsgType><![CDATA[image]]></MsgType>
	</xml>`)

	typ, err := getMsgType(data)
	a.NotError(err).Equal(typ, "image")
}

func TestGetMessageObj(t *testing.T) {
	a := assert.New(t, false)

	data := []byte(`<xml>
	<ToUserName><![CDATA[12d]]></ToUserName>
	<FromUserName><![CDATA[dddadfaee]]></FromUserName>
	<CreateTime>12345555</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[subscribe]]></Event>
	</xml>`)

	msg, err := getMessageObj(data)
	a.NotError(err)
	_, ok := msg.(*EventScan)
	a.True(ok)

	// 消息
	data = []byte(`<xml>
	<ToUserName><![CDATA[12d]]></ToUserName>
	<FromUserName><![CDATA[dddadfaee]]></FromUserName>
	<CreateTime>12345555</CreateTime>
	<MsgType><![CDATA[text]]></MsgType>
	</xml>`)

	msg, err = getMessageObj(data)
	a.NotError(err)
	_, ok = msg.(*Text)
	a.True(ok)
}
