// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package open

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	preAuthCodeURL          = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token="
	componentBindURL        = "https://mp.weixin.qq.com/safe/bindcomponent?action=bindcomponent&no_scan=1&component_appid=%s&pre_auth_code=%s&redirect_uri=%s&%s#wechat_redirect"
	componentLoginPageURL   = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s"
	componentAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	apiQueryAuthURL         = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
	apiAuthorizerTokenURL   = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token="
	apiGetAuthorizerInfoURL = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token="
)

// PreAuthCode 从预授权地址返回的内容
type PreAuthCode struct {
	Code    string `json:"pre_auth_code"`
	Expires int    `json:"expires_in"`
}

// GetPreAuthCode 获取预授权码
func GetPreAuthCode(appid, accessToken string) (*PreAuthCode, error) {
	url := preAuthCodeURL + accessToken
	body := bytes.NewBufferString(`{"component_appid":"` + appid + `"}`)

	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	p := &PreAuthCode{}
	if err := json.Unmarshal(data, p); err != nil {
		return nil, err
	}
	return p, nil
}

// GetComponentLoginPageURL 生成授权引导页网址
func GetComponentLoginPageURL(appid, preAuthCode, redirectURL string, authType AuthType) string {
	url := fmt.Sprintf(componentLoginPageURL, appid, preAuthCode, redirectURL)
	if authType == 0 {
		authType = AuthTypeAll
	}

	if authType != AuthTypeAll {
		url += "&auth_type=" + strconv.Itoa(int(authType))
	}

	return url
}

// GetComponentBindURL 生成授权链接
func GetComponentBindURL(appid, preAuthCode, redirectURL string, authType AuthType, bizAppID string) string {
	var typ string
	if bizAppID != "" {
		typ = "biz_appid=" + bizAppID
	} else {
		if authType == 0 {
			authType = AuthTypeAll
		}
		typ = "auth_type=" + strconv.Itoa(int(authType))
	}

	return fmt.Sprintf(componentBindURL, appid, preAuthCode, redirectURL, typ)
}

// GetComponentAccessToken 获取第三方平台 component_access_token
func GetComponentAccessToken(appid, appsecret, verityTicket string) (token string, expiresIn int, err error) {
	type request struct {
		AppID  string `json:"component_appid"`
		Secret string `json:"component_appsecret"`
		Ticket string `json:"component_verify_ticket"`
	}

	type response struct {
		Token   string `json:"component_access_token"`
		Expired int    `json:"expires_in"`
	}

	data, err := xml.Marshal(&request{
		AppID:  appid,
		Secret: appsecret,
		Ticket: verityTicket,
	})
	if err != nil {
		return "", 0, err
	}

	resp, err := http.Post(componentAccessTokenURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", 0, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	obj := &response{}
	if err = xml.Unmarshal(data, obj); err != nil {
		return "", 0, err
	}
	return obj.Token, obj.Expired, nil
}

type QueryAuth struct {
	AuthorizationInfo *AuthorizationInfo `json:"authorization_info"`
}

type AuthorizationInfo struct {
	AppID        string      `json:"authorizer_appid"`
	AccessToken  string      `json:"authorizer_access_token"`
	ExpiresIn    int         `json:"expires_in"`
	RefreshToken string      `json:"authorizer_refresh_token"`
	FuncInfo     []*FuncInfo `json:"func_info"`
}

type FuncInfo struct {
	Scope *Scope `json:"funcscope_category"`
}

type Scope struct {
	ID string `json:"id"`
}

// GetQueryAuth 使用授权码换取公众号的接口调用凭据和授权信息
func GetQueryAuth(appid, authorizationCode, componentAccessToken string) (*QueryAuth, error) {
	type request struct {
		AppID string `json:"component_appid"`
		Code  string `json:"authorization_code"`
	}

	data, err := xml.Marshal(&request{
		AppID: appid,
		Code:  authorizationCode,
	})
	if err != nil {
		return nil, err
	}

	url := apiQueryAuthURL + componentAccessToken
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	obj := &QueryAuth{}
	if err = xml.Unmarshal(data, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

type AuthorizerToken struct {
	AccessToken  string `json:"authorizer_access_token"`  // 授权方令牌
	ExpiresIn    int    `json:"expires_in"`               // 有效期，为2小时
	RefreshToken string `json:"authorizer_refresh_token"` // 刷新令牌
}

// GetAuthorizerToken 获取（刷新）授权公众号的接口调用凭据（令牌）
func GetAuthorizerToken(appid, authorizerAppid, authorizerRefreshToken, componentAccessToken string) (*AuthorizerToken, error) {
	type request struct {
		AppID            string `json:"component_appid"`
		AuthAppid        string `json:"authorizer_appid"`
		AuthRefreshToken string `json:"authorizer_refresh_token"`
	}

	data, err := xml.Marshal(&request{
		AppID:            appid,
		AuthAppid:        authorizerAppid,
		AuthRefreshToken: authorizerRefreshToken,
	})
	if err != nil {
		return nil, err
	}

	url := apiAuthorizerTokenURL + componentAccessToken
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	obj := &AuthorizerToken{}
	if err = xml.Unmarshal(data, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

type AuthorizerInfo struct {
	Nickname      string `json:"nick_name"`
	Headimg       string `json:"head_img"`
	Username      string `json:"user_name"`
	PrincipalName string `json:"principal_name"`
	QrcodeURL     string `json:"qrcode_url"`
	Signature     string `json:"signature"`
}

type AuthorizerObj struct {
	Info *AuthorizerInfo `json:"authorizer_info"`
}

// GetAuthorizerInfo 获取授权方的公众号帐号基本信息
func GetAuthorizerInfo(appid, authorizerAppid, componentAccessToken string) (*AuthorizerObj, error) {
	type request struct {
		AppID     string `json:"component_appid"`
		AuthAppid string `json:"authorizer_appid"`
	}

	data, err := xml.Marshal(&request{
		AppID:     appid,
		AuthAppid: authorizerAppid,
	})
	if err != nil {
		return nil, err
	}

	url := apiGetAuthorizerInfoURL + componentAccessToken
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	obj := &AuthorizerObj{}
	if err = xml.Unmarshal(data, obj); err != nil {
		return nil, err
	}
	return obj, nil
}
