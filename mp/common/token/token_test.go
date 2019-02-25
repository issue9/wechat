// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package token

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/wechat/mp/common/result"
)

func TestToken_IsExpired(t *testing.T) {
	a := assert.New(t)

	token := &AccessToken{
		ExpiresIn: 7200 * time.Second,
		Created:   time.Now(),
	}
	a.False(token.IsExpired())

	time.Sleep(1 * time.Second)
	token.ExpiresIn = 1
	a.True(token.IsExpired())
}

func TestParseAccessToken(t *testing.T) {
	a := assert.New(t)

	r := bytes.NewReader([]byte(`{"errcode":40001,"errmsg":"error msg"}`))
	at, err := parseAccessToken(r)
	rslt, ok := err.(*result.Result)
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
