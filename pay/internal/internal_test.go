// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func TestMapFromReader(t *testing.T) {
	a := assert.New(t)

	buf := bytes.NewBufferString(`<xml>
	<appid>1234567</appid>
	<mch_id><![CDATA[mch_id123]]></mch_id>
	</xml>`)

	m, err := MapFromXMLReader(buf)
	a.NotError(err).NotNil(m)
	a.Equal(m["appid"], "1234567")
	a.Equal(m["mch_id"], "mch_id123")
}

func TestMap2XMLObj(t *testing.T) {
	a := assert.New(t)

	maps := map[string]string{
		"appid":      "12345",
		"mch_id":     "mch_id123",
		"count":      "55",
		"not_exists": "not_exists",
	}
	obj := &struct {
		AppID string `xml:"appid"`
		MchID string `xml:"mch_id"`
		Count int    `xml:"count"`
	}{}

	err := Map2XMLObj(maps, obj)
	a.NotError(err)

	a.Equal(obj.AppID, "12345").
		Equal(obj.MchID, "mch_id123").
		Equal(obj.Count, 55)
}
