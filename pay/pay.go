// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

import (
	"bytes"
	"encoding/xml"
	"net/http"

	"github.com/issue9/wechat/pay/internal/utils"
)

// 返回状态的值
const (
	Success = "SUCCESS"
	Fail    = "FAIL"
)

// 接口地址
const (
	UnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	OrderQueryURL   = "https://api.mch.weixin.qq.com/pay/orderquery"
	CloseOrderURL   = "https://api.mch.weixin.qq.com/pay/closeorder"
	RefundURL       = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	RefundQueryURL  = "https://api.mch.weixin.qq.com/pay/refundquery"
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
		params["sign"] = utils.Sign(conf.APIKey, params)
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
