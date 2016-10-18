// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package token

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/issue9/wechat/result"
)

// AccessToken 用于描述从 https://api.weixin.qq.com/cgi-bin/token 正常返回的数据结构。
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Created     int64  // 该 access_token 的获取时间
}

// IsExpired 该 access_token 是否还在有效期之内。
func (at *AccessToken) IsExpired() bool {
	return time.Now().Unix() > (at.Created + at.ExpiresIn)
}

// Refresh 获取 access_token 值。
//
// 用户最好自己实现一个处理 access_token 的中控服务器来集中处理 access_token 的获取与更新。
func Refresh(appid, secret string) (*AccessToken, error) {
	if len(appid) == 0 {
		return nil, result.New(41002)
	}

	if len(secret) == 0 {
		return nil, result.New(41004)
	}

	resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret)
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 错误码
	if bytes.Index(data, []byte("errcode")) > 0 {
		return nil, result.From(data)
	}

	// access_token
	at := &AccessToken{}
	if err := json.Unmarshal(data, at); err != nil {
		return nil, err
	}
	at.Created = time.Now().Unix()
	return at, nil
}
