// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package refund 执行退款操作
//
//  p := pay.New(...)
//  r := refund.Refund{
//      Pay: p,
//      OpUserID: "10001",
//      SignType: pay.SignTypeMD5,
//  }
//
//  // 执行退款操作
//  r.OutTradeNO(...)
package refund

import (
	"strconv"

	"github.com/issue9/wechat/pay"
)

// Refund 退款数据
type Refund struct {
	Pay           *pay.Pay
	DeviceInfo    string // 设备信息
	SignType      string // 签名类型
	RefundFeeType string // 货币类型
	OpUserID      string // 操作员帐号, 默认为商户号
	RefundAccount string // 退款资金来源
}

func (r *Refund) params(outRefundNO, outTradeNO, transactionID string, totalFee, refundFee int) map[string]string {
	return map[string]string{
		"device_info":     r.DeviceInfo,
		"sign_type":       r.SignType,
		"out_trade_no":    outTradeNO,
		"transaction_id":  transactionID,
		"out_refund_no":   outRefundNO,
		"total_fee":       strconv.Itoa(totalFee),
		"refund_fee":      strconv.Itoa(refundFee),
		"refund_fee_type": r.RefundFeeType,
		"op_user_id":      r.OpUserID,
		"refund_account":  r.RefundAccount,
	}
}

// OutTradeNO 通过 outTradeNO 执行退款操作
func (r *Refund) OutTradeNO(outRefundNO, outTradeNO string, totalFee, refundFee int) (*Return, error) {
	params := r.params(outRefundNO, outTradeNO, "", totalFee, refundFee)
	return r.refund(params)
}

// TransactionID 通过 transcationID 执行退款操作
func (r *Refund) TransactionID(outRefundNO, transactionID string, totalFee, refundFee int) (*Return, error) {
	params := r.params(outRefundNO, "", transactionID, totalFee, refundFee)
	return r.refund(params)
}

func (r *Refund) refund(params map[string]string) (*Return, error) {
	maps, err := r.Pay.Refund(params)
	if err != nil {
		return nil, err
	}

	if err = r.Pay.ValidateAll(r.SignType, maps); err != nil {
		return nil, err
	}

	return newReturn(maps)
}
