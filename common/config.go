// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package common

import (
	"net/url"
	"path"
)

// Config 微信的基本配置内容
type Config struct {
	AppID     string
	AppSecret string
	Host      string // 主机，不包含协议和端口
}

// NewConfig 声明一个 [Config] 实例
//
// Host 表示的微信的接口域名，留空，会选择通用的域名：api.weixin.qq.com 可参考[说明]。
//
// [说明]: https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1465199793_BqlKA&token=&lang=zh_CN
func NewConfig(appid, appsecret, host string) *Config {
	if len(appid) == 0 {
		panic("参数 appid 不能为空")
	}

	if len(appsecret) == 0 {
		panic("参数 appsecret 不能为空")
	}

	if len(host) == 0 {
		host = "api.weixin.qq.com"
	}

	return &Config{
		AppID:     appid,
		AppSecret: appsecret,
		Host:      host,
	}
}

// URL 生成调用 api 的地址
//
// 根据 c.Host 不同，生成不同的地址。
func (c *Config) URL(urlpath string, queries map[string]string) string {
	us := make(url.Values, len(queries))
	for k, v := range queries {
		us.Add(k, v)
	}

	return "https://" + path.Join(c.Host, urlpath) + "?" + us.Encode()
}
