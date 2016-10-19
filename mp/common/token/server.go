// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package token

import (
	"log"
	"os"
	"time"

	"github.com/issue9/wechat/mp/common/config"
	"github.com/issue9/wechat/mp/common/result"
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

	// 获取相关的配置项
	Config() *config.Config
}

// DefaultServer 默认的 access_token 中控服务器
type DefaultServer struct {
	conf   *config.Config
	errlog *log.Logger
	token  *AccessToken
}

// NewDefaultServer 声明一个默认的 access_token 中控服务器
//
// 若将 errlog 指定为 nil，则会将错误信息输出到 stderr 中。
func NewDefaultServer(conf *config.Config, errlog *log.Logger) (*DefaultServer, error) {
	if len(conf.AppID) == 0 {
		return nil, result.New(40002)
	}

	if len(conf.AppSecret) == 0 {
		return nil, result.New(41004)
	}

	if errlog == nil {
		errlog = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)
	}

	at := &DefaultServer{
		conf:   conf,
		errlog: errlog,
	}
	at.refresh()

	return at, nil
}

// Token 获取当前的 *AccessToken
func (s *DefaultServer) Token() *AccessToken {
	return s.token
}

// Refresh 刷新 AccessToken，并获取新的 token
func (s *DefaultServer) Refresh() (*AccessToken, error) {
	token, err := Refresh(s.conf)
	if err != nil {
		return nil, err
	}
	s.token = token

	return token, nil
}

// Config 获取相关的配置对象
func (s *DefaultServer) Config() *config.Config {
	return s.conf
}

// 定时刷新
func (s *DefaultServer) refresh() {
	if _, err := s.Refresh(); err != nil {
		s.errlog.Println(err)
	}

	// 提前10分钟刷新
	time.AfterFunc(time.Duration(s.token.ExpiresIn-600)*time.Second, func() {
		s.refresh()
	})
}

// URL 生成指定地址的 URL，会在查询参数中添中 access_token 的相关设置
func URL(s Server, path string, queries map[string]string) string {
	if queries == nil {
		queries = make(map[string]string, 1)
	}
	queries["access_token"] = s.Token().AccessToken
	return s.Config().URL(path, queries)
}
