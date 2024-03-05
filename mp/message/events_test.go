// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package message

import (
	"testing"

	"github.com/issue9/assert/v4"
)

var _ Eventer = &event{}
var _ Eventer = &EventScan{}
var _ Eventer = &EventLocation{}
var _ Eventer = &EventClickView{}

func TestGetEventType(t *testing.T) {
	a := assert.New(t, false)

	data := []byte(`<xml>
	<ToUserName><![CDATA[12d]]></ToUserName>
	<FromUserName><![CDATA[dddadfaee]]></FromUserName>
	<CreateTime>12345555</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[subscribe]]></Event>
	</xml>`)

	event, err := getEventType(data)
	a.NotError(err).Equal(event, "subscribe")
}

func TestGetEventObj(t *testing.T) {
	a := assert.New(t, false)

	// subscribe
	data := []byte(`<xml>
	<ToUserName><![CDATA[12d]]></ToUserName>
	<FromUserName><![CDATA[dddadfaee]]></FromUserName>
	<CreateTime>12345555</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[subscribe]]></Event>
	</xml>`)

	event, err := getEventObj(data)
	a.NotError(err)
	_, ok := event.(*EventScan)
	a.True(ok)

	// click
	data = []byte(`<xml>
	<ToUserName><![CDATA[12d]]></ToUserName>
	<FromUserName><![CDATA[dddadfaee]]></FromUserName>
	<CreateTime>12345555</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[CLICK]]></Event>
	</xml>`)

	event, err = getEventObj(data)
	a.NotError(err)
	_, ok = event.(*EventClickView)
	a.True(ok)
}
