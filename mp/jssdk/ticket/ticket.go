// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package ticket

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/issue9/wechat/mp/common/token"
)

// Ticket 表示 jsjapi 的 ticket 类型
type Ticket struct {
	Code      int    `json:"errcode"`
	Msg       string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// Config 表示 jssdk 中 wx.config 中的参数。
type Config struct {
	Debug       bool     `json:"debug"`
	AppID       string   `json:"appId"`
	Timestamp   int64    `json:"timestamp"`
	NonceString string   `json:"nonceStr"`
	Signature   string   `json:"signature"`
	APIList     []string `json:"jsApiList"`
}

// Refresh 获取相关的 Ticket 值。
func Refresh(srv token.Server) (*Ticket, error) {
	url := token.URL(srv, "/cgi-bin/ticket/getticket", map[string]string{"type": "jsapi"})
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	t := &Ticket{}
	if err := json.Unmarshal(data, t); err != nil {
		return nil, err
	}

	return t, nil
}
