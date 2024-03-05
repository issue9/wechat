// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package ticket

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/issue9/errwrap"

	"github.com/issue9/wechat/mp/common/token"
	"github.com/issue9/wechat/pay"
)

// Server 表示中控服务器接口
type Server interface {
	// 获取中控服务器缓存的 access_token。
	Ticket() *Ticket

	// 刷新中控服务器的 access_token。
	//
	// 中控服务器应该提供自动刷新机制。
	// 此函数的存在，仅仅是为了在某些特定的情况下，手动刷 access_token 使用。
	Refresh() (*Ticket, error)

	// 根据当前的 Ticket 生成相应的 Config 实例。
	Config(string) (*Config, error)
}

// DefaultServer 默认的 access_token 中控服务器
type DefaultServer struct {
	tokenSrv token.Server
	errlog   *log.Logger
	ticket   *Ticket
}

// NewDefaultServer 声明一个默认的 access_token 中控服务器
//
// 若将 errlog 指定为 nil，则会将错误信息输出到 stderr 中。
func NewDefaultServer(tksrv token.Server, errlog *log.Logger) *DefaultServer {
	if errlog == nil {
		errlog = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)
	}

	srv := &DefaultServer{
		tokenSrv: tksrv,
		errlog:   errlog,
	}
	srv.refresh()

	return srv
}

// Ticket 获取当前的 *Ticket
func (s *DefaultServer) Ticket() *Ticket {
	return s.ticket
}

// Refresh 刷新 Ticket，并获取新的 token
func (s *DefaultServer) Refresh() (*Ticket, error) {
	ticket, err := Refresh(s.tokenSrv)
	if err != nil {
		return nil, err
	}
	s.ticket = ticket

	return ticket, nil
}

// Config 表示 Config 实例
func (s *DefaultServer) Config(url string) (*Config, error) {
	now := time.Now()
	nonceStr := pay.NonceString()

	sign, err := sign(map[string]string{
		"noncestr":     nonceStr,
		"jsapi_ticket": s.Ticket().Ticket,
		"timestamp":    strconv.FormatInt(now.Unix(), 10),
		"url":          url,
	})
	if err != nil {
		return nil, err
	}

	return &Config{
		Debug:       true,
		AppID:       s.tokenSrv.Config().AppID,
		Timestamp:   now.Unix(),
		NonceString: nonceStr,
		Signature:   sign,
	}, nil

}

// 定时刷新
func (s *DefaultServer) refresh() {
	if _, err := s.Refresh(); err != nil {
		s.errlog.Println(err)
	}

	// 提前10分钟刷新
	time.AfterFunc(time.Duration(s.ticket.ExpiresIn-600)*time.Second, func() {
		s.refresh()
	})
}

// Sign 微信支付签名
//
// params 签名用的参数
func sign(params map[string]string) (string, error) {
	/* 排序 */
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}

		keys = append(keys, k)
	}
	sort.Strings(keys)

	/* 拼接字符串 */
	var buf errwrap.Buffer
	for _, k := range keys {
		v := params[k]
		if len(v) == 0 {
			continue
		}

		buf.WString(k).WByte('=').WString(v).WByte('&')
	}
	buf.Truncate(buf.Len() - 1) // 去掉最后的 &
	if buf.Err != nil {
		return "", buf.Err
	}

	h := sha1.New()
	h.Write(buf.Bytes())
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}
