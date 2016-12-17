// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package utils

import (
	"bytes"
	"encoding/hex"
	"hash"
	"sort"
	"strings"
	"time"

	"github.com/issue9/rands"
)

var (
	bs      = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randSrv *rands.Rands
)

func init() {
	randSrv, _ = rands.New(time.Now().Unix(), 100, 20, 21, bs)
}

// NonceStr 获取随机字符串
func NonceStr() string {
	if randSrv != nil {
		return randSrv.String()
	}
	return rands.String(20, 21, bs)
}

// Sign 微信支付签名
//
// apikey 支付用的 apikey
// params 签名用的参数
// fn 签名的类型，为空则为 md5
func Sign(apikey string, params map[string]string, h hash.Hash) string {
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

	h.Write(buf.Bytes())
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
