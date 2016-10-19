// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"sort"
)

// Handler 消息处理函数。
//
// 而返回值则即为需要向请求源输出的字符串，一般为 reply.go 中的相关类型。
type Handler func(Messager) ([]byte, error)

// Server 消息管理服务器。
type Server struct {
	token   string
	handler Handler
	errlog  *log.Logger
}

// NewServer 声明一个新的消息管理服务器。
//
// 若将 errlog 指定为 nil，则会将错误信息输出到 stderr 中。
func NewServer(token string, h Handler, errlog *log.Logger) *Server {
	if errlog == nil {
		errlog = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)
	}

	return &Server{
		token:   token,
		handler: h,
		errlog:  errlog,
	}
}

// Signature 验证签名，GET 方法
func (s *Server) Signature(w http.ResponseWriter, r *http.Request) {
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")

	if signature == sign(s.token, timestamp, nonce) {
		w.Write([]byte(r.FormValue("echostr")))
		return
	}
}

// Message 消息处理，POST 方法
func (s *Server) Message(w http.ResponseWriter, r *http.Request) {
	obj, err := getMsgObj(r.Body)
	if err != nil {
		s.errlog.Println(err)
		return
	}

	bs, err := s.handler(obj)
	if err != nil {
		s.errlog.Println(err)
		return
	}

	w.Write(bs)
}

// sign 微信接口地址验证方法
func sign(token, timestamp, nonce string) string {
	strs := sort.StringSlice{token, timestamp, nonce}
	strs.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce))
	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)

	hash := sha1.Sum(buf)
	return hex.EncodeToString(hash[:])
}

// TransferCustomerService 是一个默认的 Handler 实现，仅仅是实现了对消息的转发。
func TransferCustomerService(m Messager) ([]byte, error) {
	if m.Type() == TypeEvent {
		return nil, errors.New("事件不允许转发给客服")
	}

	return NewReplyTranferCustomerService(m).Bytes()
}
