// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package pay 微信支付的相关接口
package pay

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"sort"
	"strings"
	"time"

	"github.com/issue9/errwrap"
	"github.com/issue9/rands"
)

var randSrv *rands.Rands

func init() {
	srv, err := rands.New(time.Now().Unix(), 100, 24, 32, rands.AlphaNumber)
	if err != nil {
		panic(err)
	}

	randSrv = srv
}

// NonceString 产生一段随机字符串
func NonceString() string {
	return randSrv.String()
}

// Sign 微信支付签名
//
// apikey 支付用的 apikey
// params 签名用的参数
// fn 签名的类型，为空则为 md5
func Sign(apikey, signType string, params map[string]string) (string, error) {
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
	buf.WString("key=").WString(apikey)
	if buf.Err != nil {
		return "", buf.Err
	}

	var h hash.Hash
	switch signType {
	case "", SignTypeMD5:
		h = md5.New()
	case SignTypeHmacSha256:
		//h = hmac.New(sha256.New, []byte("123"))
	default:
		return "", fmt.Errorf("无效的签名类型：%v", signType)
	}

	h.Write(buf.Bytes())
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}
