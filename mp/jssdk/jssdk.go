// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package jssdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/issue9/wechat/mp/common/config"
	"github.com/issue9/wechat/mp/common/result"
)

// 授权作用域，供 GetCodeURL 使用。
const (
	SnsapiUserinfo = "snsapi_uesrnfo"
	SnsapiBase     = "snsapi_base"

	// 获取 code 的地址
	codeURL = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%v&redirect_uri=%v&response_type=code&scope=%v&state=%v#wechat_redirect"
)

// GetCodeURL 获取 code
func GetCodeURL(conf *config.Config, redirectURI, scope, state string) string {
	redirectURI = url.QueryEscape(redirectURI)
	return fmt.Sprintf(codeURL, conf.AppID, redirectURI, scope, state)
}

// GetAccessToken 根据 code 获取 access_token
func GetAccessToken(conf *config.Config, code string) (*AccessToken, error) {
	queries := map[string]string{
		"appid":      conf.AppID,
		"secret":     conf.AppSecret,
		"code":       code,
		"grant_type": "authorization_code",
	}
	url := conf.URL("sns/oauth2/access_token", queries)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseAccessToken(resp.Body)
}

// RefreshAccessToken 刷新 access_token
func RefreshAccessToken(conf *config.Config, token *AccessToken) (*AccessToken, error) {
	queries := map[string]string{
		"appid":         conf.AppID,
		"refresh_token": token.RefreshToken,
		"grant_type":    "refresh_token",
	}
	url := conf.URL("sns/oauth2/refresh_token", queries)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseAccessToken(resp.Body)
}

// GetUserInfo 获取用户基本信息
//
// 若不指定 lang 则使用 zh_CN 作为其默认值。
func GetUserInfo(conf *config.Config, token *AccessToken, lang string) (*UserInfo, error) {
	if len(lang) == 0 {
		lang = "zh_CN"
	}
	queries := map[string]string{
		"openid":       token.OpenID,
		"access_token": token.AccessToken,
		"lang":         lang,
	}
	url := conf.URL("sns/userinfo", queries)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &UserInfo{}
	if err = json.Unmarshal(data, info); err != nil {
		return nil, err
	}
	if len(info.OpenID) > 0 {
		return info, nil
	}

	rslt := &result.Result{}
	if err = json.Unmarshal(data, rslt); err != nil {
		return nil, err
	}
	return nil, rslt
}

// AuthAccessToken 验证 access_token 是否有效
func AuthAccessToken(conf *config.Config, token *AccessToken) (bool, error) {
	queries := map[string]string{
		"openid":       token.OpenID,
		"access_token": token.AccessToken,
	}
	url := conf.URL("sns/auth", queries)

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	rslt := result.From(data)
	return rslt.Code == 0, rslt
}

// 分析 r 中的数据到 AccessToken 或是 result.Result 对象中。
func parseAccessToken(r io.Reader) (*AccessToken, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	token := &AccessToken{}
	if err := json.Unmarshal(data, token); err != nil {
		return nil, err
	}
	if len(token.AccessToken) > 0 || token.ExpiresIn > 0 {
		token.Created = time.Now()
		return token, nil
	}

	rslt := &result.Result{}
	if err := json.Unmarshal(data, rslt); err != nil {
		return nil, err
	}
	return nil, rslt
}
