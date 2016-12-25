// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package pay 微信支付的相关接口
package pay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"hash"
	"sort"
	"strings"
	"time"

	"github.com/issue9/rands"
)

var (
	// 函数可用的字符
	nonceStringChars = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSrv *rands.Rands
)

func init() {
	randSrv, _ = rands.New(time.Now().Unix(), 100, 24, 32, nonceStringChars)
}

// NonceString 产生一段随机字符串
func NonceString() string {
	if randSrv == nil {
		return rands.String(24, 32, nonceStringChars)
	}

	return randSrv.String()
}

// Sign 微信支付签名
//
// apikey 支付用的 apikey
// params 签名用的参数
// fn 签名的类型，为空则为 md5
func Sign(apikey string, params map[string]string) string {
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
	buf := new(bytes.Buffer)
	for _, k := range keys {
		v := params[k]
		if len(v) == 0 {
			continue
		}

		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
		buf.WriteByte('&')
	}
	buf.WriteString("key=")
	buf.WriteString(apikey)

	var h hash.Hash
	switch params["sign_type"] {
	case SignTypeMD5:
		h = md5.New()
	case SignTypeHmacSha256:
		//h = hmac.New(sha256.New, []byte("123"))
	}

	h.Write(buf.Bytes())
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
