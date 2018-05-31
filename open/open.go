// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package open

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/issue9/wechat/open/crypto"
)

const (
	getAPIGetAuthorizerOptionURL = "https://api.weixin.qq.com/cgi-bin/component/ api_get_authorizer_option?component_access_token=%s"
	getAPISetAuthorizerOptionURL = "https://api.weixin.qq.com/cgi-bin/component/ api_set_authorizer_option?component_access_token=%s"
)

// 表示几种验证类型
const (
	AuthTypeMP    = 1
	AuthTypeWeapp = 2
	AuthTypeAll   = 3
)

// AuthType 验证类型
type AuthType int8

// ParseAuthorizationCode 从回调地址中获取授权码等信息
//
// 用户的确认授权之后的回调地址，数据从 query 中获取。
func ParseAuthorizationCode(r *http.Request) (code string, expiresIn int, err error) {
	expiresIn, err = strconv.Atoi(r.FormValue("expires_in"))
	if err != nil {
		return "", 0, err
	}

	return r.FormValue("auth_code"), expiresIn, nil
}

// VerifyTicket 返回的验证信息
type VerifyTicket struct {
	Root                  xml.Name `xml:"xml"`
	AppID                 string   `xml:"AppId"`
	CreateTime            string   `xml:"CreateTime"`
	InfoType              string   `xml:"InfoType"`
	ComponentVerifyTicket string   `xml:"ComponentVerifyTicket"`
}

// ParseVerifyTicket 处理 component_verify_ticket 事件中返回的数据
func ParseVerifyTicket(c *crypto.Crypto, w http.ResponseWriter, r *http.Request) (*VerifyTicket, error) {
	w.Write([]byte("success")) // 先返回内容

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	ticket := &VerifyTicket{}
	sign := r.FormValue("msg_signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")
	if err = c.DecryptObject(data, sign, timestamp, nonce, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
