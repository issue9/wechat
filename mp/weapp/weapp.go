// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package weapp 小程序的相关操作
package weapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/issue9/wechat/mp/common/result"
)

const (
	grantType = "authorization_code"

	loginURL = "https://api.weixin.qq.com/sns/jscode2session"
)

// Response 返回的数据
type Response struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"` // 某些情况下存在

	created time.Time
}

// Authorization 执行登录验证，并获取相应的数据
func Authorization(appid, secret, jscode string) (*Response, error) {
	vals := url.Values{}
	vals.Set("grant_type", grantType)
	vals.Set("appid", appid)
	vals.Set("secret", secret)
	vals.Set("js_code", jscode)

	url := loginURL + "?" + vals.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	data := &Response{}
	if err := json.Unmarshal(bs, data); err != nil {
		return nil, err
	}

	if data.Openid != "" { // 正常数据，肯定有 openid
		return data, nil
	}

	rslt := &result.Result{}
	if err := json.Unmarshal(bs, rslt); err != nil {
		return nil, err
	}
	return nil, rslt
}
