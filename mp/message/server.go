// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

// Handler 消息处理函数。
//通过向 NewServer 注册 Handler 函数，获取对消息处理的权限。
//
// 参数 Messager 为从微信端传递过来的 xml 数据对象实例，都已定义在 messages.go
// 文件中。
//
// 函数的返回值，被当作消息被动回复内容传递给微信调用方。在 reply.go 中
// 定义了大部分可能用到返回类型，可以拿来直接使用。
type Handler func(Messager) ([]byte, error)

// Server 消息管理服务器。
type Server struct {
	token   string
	handler Handler
	errlog  *log.Logger
}

// NewServer 声明一个新的消息管理服务器。
//
// 若将 h 参数指定为 nil，则会被自动赋予 TransferCustomerService 函数。
// 若将 errlog 指定为 nil，则会将错误信息输出到 stderr 中。
func NewServer(token string, h Handler, errlog *log.Logger) *Server {
	if errlog == nil {
		errlog = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)
	}

	if h == nil {
		h = TransferCustomerService
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
	obj, err := getObj(r.Body)
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

func getObj(r io.Reader) (Messager, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return getMessageObj(data)
}

// TransferCustomerService 是一个默认的 Handler 实现，仅仅是实现了对消息的转发。
func TransferCustomerService(m Messager) ([]byte, error) {
	if m.Type() == TypeEvent {
		return nil, errors.New("事件不允许转发给客服")
	}

	return NewReplyTranferCustomerService(m).Bytes()
}
