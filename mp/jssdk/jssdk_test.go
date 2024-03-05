// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package jssdk

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v4"

	"github.com/issue9/wechat/common"
)

func TestParseAccessToken(t *testing.T) {
	a := assert.New(t, false)

	r := bytes.NewReader([]byte(`{"errcode":40001,"errmsg":"error msg"}`))
	at, err := parseAccessToken(r)
	rslt, ok := err.(*common.Result)
	a.Error(err).Nil(at)
	a.True(ok).False(rslt.IsOK())

	// 带 errmsg 不会被误判
	r = bytes.NewReader([]byte(`{"access_token":"errmsg","expires_in":13334232}`))
	at, err = parseAccessToken(r)
	a.NotError(err).
		NotNil(at).
		Equal(at.AccessToken, "errmsg").
		Equal(at.ExpiresIn, 13334232).
		True(at.Created.Unix() > 0)

	// 解析错误
	r = bytes.NewReader([]byte(`{"access_token":123,"expires_in":13334232}`))
	at, err = parseAccessToken(r)
	a.Error(err).Nil(at)

	// 参数错误
	r = bytes.NewReader(nil)
	at, err = parseAccessToken(r)
	a.Error(err).Nil(at)
}
