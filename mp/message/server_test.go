// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package message

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestGetObj(t *testing.T) {
	a := assert.New(t, false)

	data := []byte(`<xml>
	<MsgType>event</MsgType>
	<Event>subscribe</Event>
	</xml>`)
	buf := bytes.NewReader(data)

	msg, err := getObj(buf)
	a.NotError(err)
	obj1, ok := msg.(*EventScan)
	a.True(ok).False(obj1.IsScan())

	// 消息
	data = []byte(`<xml>
	<MsgType>text</MsgType>
	<Content>cc</Content>
	</xml>`)
	buf = bytes.NewReader(data)

	_, err = getObj(buf)
	a.NotError(err)
	//obj2, ok := msg.(*Text)
	//a.True(ok).Equal(obj2.Content, "cc")
}
