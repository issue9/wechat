// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package message

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

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
