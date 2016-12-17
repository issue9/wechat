// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unifiedorder

import "strings"

const (
	success = "SUCCESS"
	fail    = "FAIL"
)

// Return 表示统一下单功能的返回值类型。
type Return struct {
	Code    string `xml:"return_code"`
	Message string `xml:"return_msg"`

	// 当 code == success 时，返回值拥有以下值
	AppID       string `xml:"appid"`
	MchID       string `xml:"mch_id"`
	DeviceInfo  string `xml:"device_info"`
	NonceStr    string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	ResultCode  string `xml:"result_code"`
	ErrCode     string `xml:"err_code"`
	ErrCodeDesc string `xml:"err_code_des"`

	// 当 ResultCode == success 时，拥有的返回值
	TradeType string `xml:"trade_type"`
	PrepayID  string `xml:"prepay_id"`
	CodeURL   string `xml:"code_url"`
}

// IsOK 通信是否正常，即 Result.Code 是否为 SUCCESS
func (r *Return) IsOK() bool {
	return strings.ToUpper(r.Code) == success
}

// IsSuccess 交易是否正常，即 (Result.Code == SUCCESS) && (Result.ResultCode == SUCCESS)
func (r *Return) IsSuccess() bool {
	return r.IsOK() && strings.ToUpper(r.ResultCode) == success
}
