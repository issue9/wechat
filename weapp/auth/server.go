// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/issue9/wechat/common"
)

// ErrWeappUnauthorization 表示微信未登录，或是登录已经过期
var ErrWeappUnauthorization = errors.New("微信未登录")

// Server 小程序状态管理服务
type Server struct {
	conf *common.Config

	// 保存从服务端返回的内容，键名为 openid，键值为整个返回值
	tokens       map[string]*Response
	tokensLocker sync.RWMutex
	expired      time.Duration
}

// NewServer 声明一个新的 Server
//
// cap 表示初始容量。
// gctick 表示 GC 的启动频率，根据业务量自定义一个合理的值。
func NewServer(conf *common.Config, cap int, gctick, expired time.Duration) *Server {
	srv := &Server{
		conf:    conf,
		tokens:  make(map[string]*Response, cap),
		expired: expired,
	}

	go srv.gc(gctick)
	return srv
}

// New 申请一个新的登录 token
func (srv *Server) New(jscode string) (*Response, error) {
	resp, err := Authorization(srv.conf, jscode)
	if err != nil {
		return nil, err
	}

	resp.created = time.Now()
	srv.tokensLocker.Lock()
	srv.tokens[resp.Openid] = resp // 已经存在则替换，否则添加
	srv.tokensLocker.Unlock()

	return resp, nil
}

// Decode 解码
func (srv *Server) Decode(openid, data, iv string) ([]byte, *Watermark, error) {
	srv.tokensLocker.RLock()
	defer srv.tokensLocker.RUnlock()

	resp, found := srv.tokens[openid]
	if !found {
		return nil, nil, ErrWeappUnauthorization
	}

	return Decode(srv.conf.AppID, resp.SessionKey, data, iv)
}

func (srv *Server) gc(tick time.Duration) {
	c := time.Tick(tick)

	for t := range c {
		func() { // 包装在函数中，才能保证 tokensLocker 不死锁
			srv.tokensLocker.Lock()
			defer srv.tokensLocker.Unlock()

			for _, token := range srv.tokens {
				if token.created.Add(srv.expired).Before(t) {
					delete(srv.tokens, token.Openid)
				}
			}
		}()
	}
}
