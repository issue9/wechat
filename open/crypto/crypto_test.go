// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package crypto

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/rands"
)

func TestCrypto_encrypt_decrypt(t *testing.T) {
	a := assert.New(t)
	c, err := New("wx123458de9ae3rdew", "token", rands.String(43, 44, randstr))
	a.NotError(err).NotNil(c)

	msg := `<xml>
	<ToUserName><![CDATA[示例内容]]></ToUserName>
	<FromUserName><![CDATA[示例内容]]></FromUserName>
	<CreateTime>1348831860</CreateTime>
	<MsgType><![CDATA[示例内容]]></MsgType>
	<Content><![CDATA[示例内容]]></Content>
	<MsgId>1234567890123456</MsgId>
	</xml>`

	text, err := c.encrypt([]byte(msg))
	a.NotError(err).NotNil(text)

	detext, err := c.decrypt(text)
	a.NotError(err).NotNil(detext)

	a.Equal(string(detext), msg)
}
