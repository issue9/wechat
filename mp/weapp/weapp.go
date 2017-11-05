// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package weapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	grantType = "authorization_code"

	loginURL = "https://api.weixin.qq.com/sns/jscode2session"
)

// LoginResponse 返回的数据
type LoginResponse struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ExpiresIn  string `json:"expires_in"`
}

// Authorization 执行登录验证，并获取相应的数据
func Authorization(appid, secret, jscode string) (*LoginResponse, error) {
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

	data := &LoginResponse{}
	if err := json.Unmarshal(bs, data); err != nil {
		return nil, err
	}

	return data, nil
}
