// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package token

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/issue9/wechat/mp/common/config"
	"github.com/issue9/wechat/mp/common/result"
)

// AccessToken 用于描述从 https://api.weixin.qq.com/cgi-bin/token 正常返回的数据结构。
type AccessToken struct {
	AccessToken string        `json:"access_token"`
	ExpiresIn   time.Duration `json:"expires_in"`
	Created     time.Time     // 该 access_token 的获取时间
}

// IsExpired 该 access_token 是否还在有效期之内。
func (at *AccessToken) IsExpired() bool {
	return time.Now().After(at.Created.Add(at.ExpiresIn))
}

// Refresh 刷新 access_token
//
// 用户最好自己实现一个处理 access_token 的中控服务器来集中处理 access_token 的获取与更新。
func Refresh(conf *config.Config) (*AccessToken, error) {
	queries := map[string]string{
		"grant_type": "client_credential",
		"appid":      conf.AppID,
		"secret":     conf.AppSecret,
	}
	resp, err := http.Get(conf.URL("cgi-bin/token", queries))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 { // 400 以上的状态码，直接输出错误信息
		return nil, &result.Result{
			Code:    resp.StatusCode,
			Message: resp.Status,
		}
	}

	return parseAccessToken(resp.Body)
}

// 将 r 中的数据分析到 AccessToken 中
func parseAccessToken(r io.Reader) (*AccessToken, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// access_token
	at := &AccessToken{}
	if err := json.Unmarshal(data, at); err != nil {
		return nil, err
	}
	if len(at.AccessToken) > 0 { // 正常读取，必须有 access_token 字段
		at.Created = time.Now()
		return at, nil
	}

	rslt := &result.Result{}
	if err := json.Unmarshal(data, rslt); err != nil {
		return nil, err
	}
	return nil, rslt
}
