// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unifiedorder

// Return 表示统一下单功能的返回值类型。
type Return struct {
	TradeType string `xml:"trade_type"`
	PrepayID  string `xml:"prepay_id"`
	CodeURL   string `xml:"code_url"`
}

func newReturn(params map[string]string) *Return {
	return &Return{
		TradeType: params["trade_type"],
		PrepayID:  params["prepay_id"],
		CodeURL:   params["code_url"],
	}
}
