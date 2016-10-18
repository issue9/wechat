// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package token

import (
	"log"
	"os"
	"time"

	"github.com/issue9/wechat/result"
)

// Server 表示中控服务器接口
type Server interface {
	// 获取中控服务器缓存的 access_token。
	Token() *AccessToken

	// 刷新中控服务器的 access_token。
	//
	// 中控服务器应该提供自动刷新机制。
	// 此函数的存在，仅仅是为了在某些特定的情况下，手动刷 access_token 使用。
	Refresh() (*AccessToken, error)
}

// 默认的 access_token 中控服务器
type AccessTokenServer struct {
	appid  string
	secret string
	errlog *log.Logger

	token *AccessToken
}

// NewAccessTokenSever 声明一个默认的 access_token 中控服务器
func NewAccessTokenServer(appid, secret string, errlog *log.Logger) (*AccessTokenServer, error) {
	if len(appid) == 0 {
		return nil, result.New(40002)
	}

	if len(secret) == 0 {
		return nil, result.New(41004)
	}

	if errlog == nil {
		errlog = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)
	}

	at := &AccessTokenServer{
		appid:  appid,
		secret: secret,
		errlog: errlog,
	}
	at.refresh()

	return at, nil
}

// Token 获取当前的 *AccessToken
func (s *AccessTokenServer) Token() *AccessToken {
	return s.token
}

// Refresh 刷新 AccessToken，并获取新的 token
func (s *AccessTokenServer) Refresh() (*AccessToken, error) {
	token, err := Refresh(s.appid, s.secret)
	if err != nil {
		return nil, err
	}
	s.token = token

	return token, nil
}

// 定时刷新
func (s *AccessTokenServer) refresh() {
	if _, err := s.Refresh(); err != nil {
		s.errlog.Println(err)
	}

	// 提交10分钟刷
	time.AfterFunc(time.Duration(s.token.ExpiresIn-600)*time.Second, func() {
		s.refresh()
	})
}
