// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"net/url"
	"path"

	"github.com/issue9/wechat/mp/common/result"
)

// Config 微信的基本配置内容
type Config struct {
	AppID     string
	AppSecret string
	Host      string // 主机，不包含协议和端口
}

// New 声明一个 Config 实例。
//
// Host 表示的微信的接口域名，留空，会选择通用的域名：
// api.weixin.qq.com
//
// 详细说明在：
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1465199793_BqlKA&token=&lang=zh_CN
func New(appid, appsecret, host string) (*Config, error) {
	if len(appid) == 0 {
		return nil, result.New(40002)
	}

	if len(appsecret) == 0 {
		return nil, result.New(41004)
	}

	if len(host) == 0 {
		host = "api.weixin.qq.com"
	}

	return &Config{
		AppID:     appid,
		AppSecret: appsecret,
		Host:      host,
	}, nil
}

// URL 生成调用 api 的地址。根据 c.Host 不同，生成不同的地址。
func (c *Config) URL(urlpath string, queries map[string]string) string {
	us := make(url.Values, len(queries))
	for k, v := range queries {
		us.Add(k, v)
	}

	return "https://" + path.Join(c.Host, urlpath) + "?" + us.Encode()
}
