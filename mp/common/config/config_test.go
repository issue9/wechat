// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/issue9/assert/v2"
)

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	conf, err := New("appid", "appsecret", "")
	a.NotError(err).Equal(conf.Host, "api.weixin.qq.com")
}

func TestConfig_URL(t *testing.T) {
	a := assert.New(t, false)

	conf := &Config{
		AppID:     "appid",
		AppSecret: "secret",
		Host:      "api.domain",
	}

	url := conf.URL("test", map[string]string{"a": "b", "c": "d"})
	a.Equal(url, "https://api.domain/test?a=b&c=d")

	url = conf.URL("test", nil)
	a.Equal(url, "https://api.domain/test?")
}
