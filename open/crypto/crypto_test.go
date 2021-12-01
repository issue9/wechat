// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package crypto

import (
	"encoding/xml"
	"strconv"
	"testing"
	"time"

	"github.com/issue9/assert/v2"
	"github.com/issue9/rands"
)

type message struct {
	Root         xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   string   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgID        string   `xml:"MsgId"`
}

func TestCrypto_encrypt_decrypt(t *testing.T) {
	a := assert.New(t, false)
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

func TestCrypto_Encrypt_DecryptObject(t *testing.T) {
	a := assert.New(t, false)
	c, err := New("wx123458de9ae3rdew", "token", rands.String(43, 44, randstr))
	a.NotError(err).NotNil(c)

	timesamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := nonceString()

	// encryptObject
	obj := &message{
		MsgID:      "1234567890123456",
		CreateTime: "1348831860",
		Content:    "示例内容",
	}
	text, sign, err := c.EncryptObject(obj, timesamp, nonce)
	a.NotError(err).NotNil(text)

	// decryptObject
	msgobj := &message{}
	a.NotError(c.DecryptObject(text, sign, timesamp, nonce, msgobj))

	a.Equal(msgobj.MsgID, obj.MsgID)
	a.Equal(msgobj.CreateTime, obj.CreateTime)
	a.Equal(msgobj.Content, obj.Content)
}
