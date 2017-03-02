// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unifiedorder

import (
	"strconv"
	"time"

	"github.com/issue9/wechat/pay"
)

// Return 表示统一下单功能的返回值类型。
type Return struct {
	pay       *pay.Pay
	TradeType string
	PrepayID  string
	CodeURL   string // 二维码链接
}

// BrandWCPayRequest 返回给微信浏览器的数据
type BrandWCPayRequest struct {
	AppID       string `json:"appId"`
	TimeStamp   string `json:"timeStamp"`
	NonceString string `json:"nonceStr"`
	Package     string `json:"package"`
	SignType    string `json:"signType"`
	PaySign     string `json:"paySign"`
}

// GetBrandWCPayRequest 获取 BrandWCPayRequest 数据
func (r *Return) GetBrandWCPayRequest(signType string) *BrandWCPayRequest {
	now := time.Now().Unix()
	ret := &BrandWCPayRequest{
		AppID:       r.pay.AppID(),
		TimeStamp:   strconv.FormatInt(now, 10),
		NonceString: pay.NonceString(),
		Package:     "prepay_id=" + r.PrepayID,
		SignType:    signType,
	}

	params := map[string]string{
		"appId":     ret.AppID,
		"timeStamp": ret.TimeStamp,
		"nonceStr":  ret.NonceString,
		"package":   ret.Package,
		"signType":  ret.SignType,
	}
	ret.PaySign = r.pay.Sign(signType, params)

	return ret
}
