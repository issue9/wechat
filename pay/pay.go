// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package pay 微信支付的相关接口
package pay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"hash"
	"net/http"
	"sort"
	"strings"

	"github.com/issue9/rands"
)

// Post 发送请求，会优先使用 params 中的相关参数。
// 比如：若已经指定了 appid，会不会使用 conf.AppID；
// 若使用了 sign，则不会再计算 sign 值。
func Post(conf *Config, url string, params Paramser, ret Returner) error {
	ps, err := params.Params()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = map2XML(conf, ps, buf); err != nil {
		return err
	}
	resp, err := http.Post(url, "application/xml", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return ret.From(resp.Body)
}

var nonceStringChars = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// NonceString 产生一段随机字符串
func NonceString() string {
	return rands.String(24, 31, nonceStringChars)
}

func map2XML(conf *Config, params map[string]string, buf *bytes.Buffer) error {
	if params["appid"] == "" {
		params["appid"] = conf.AppID
	}

	if params["mch_id"] == "" {
		params["mch_id"] = conf.MchID
	}

	if params["nonce_str"] == "" {
		params["nonce_str"] = conf.NonceStr
	}

	if params["sign_type"] == "" {
		params["sign_type"] = conf.SignType
	}

	if params["sign"] == "" {
		params["sign"] = sign(conf.APIKey, params)
	}

	buf.WriteString("<xml>")
	for k, v := range params {
		if v == "" {
			continue
		}

		buf.WriteByte('<')
		buf.WriteString(k)
		buf.WriteByte('>')

		if err := xml.EscapeText(buf, []byte(v)); err != nil {
			return err
		}

		buf.WriteString("</")
		buf.WriteString(k)
		buf.WriteByte('>')
	}
	buf.WriteString("</xml>")

	return nil
}

// 微信支付签名
//
// apikey 支付用的 apikey
// params 签名用的参数
// fn 签名的类型，为空则为 md5
func sign(apikey string, params map[string]string) string {
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
