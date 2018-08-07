// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package auth 小程序登录验证的相关操作
package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/issue9/wechat/mp/common/config"
	"github.com/issue9/wechat/mp/common/result"
)

// TODO 采用 mp/common/config 包的配置

const (
	grantType = "authorization_code"
)

// Response 返回的数据
type Response struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"` // 某些情况下存在

	created time.Time
}

// Authorization 执行登录验证，并获取相应的数据
func Authorization(conf *config.Config, jscode string) (*Response, error) {
	queries := map[string]string{
		"grant_type": grantType,
		"appid":      conf.AppID,
		"secret":     conf.AppSecret,
		"js_code":    jscode,
	}

	resp, err := http.Get(conf.URL("sns/jscode2session", queries))
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
