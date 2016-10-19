// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"net/url"
	"path"
)

// Config 微信的基本配置内容
type Config struct {
	AppID     string
	AppSecret string
	Host      string // 主机，不包含协议和端品
}

// URL 生成调用 api 的地址
func (c *Config) URL(urlpath string, queries map[string]string) string {
	us := make(url.Values, len(queries))
	for k, v := range queries {
		us.Add(k, v)
	}

	return "https://" + path.Join(c.Host, urlpath) + "?" + us.Encode()
}
