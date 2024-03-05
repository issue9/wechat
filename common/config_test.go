// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package common

import (
	"testing"

	"github.com/issue9/assert/v4"
)

func TestNewConfig(t *testing.T) {
	a := assert.New(t, false)

	conf := NewConfig("appid", "appsecret", "")
	a.Equal(conf.Host, "api.weixin.qq.com")
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
