// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package refund

import (
	"strconv"

	"github.com/issue9/wechat/pay"
)

// Refund 退款数据
type Refund struct {
	DeviceInfo    string // 设备信息
	SignType      string // 签名类型
	TransactionID string // 微信订单号
	OutTradeNO    string // 商户订单号，与微信订单号，必须二选一
	OutRefundNO   string // 商户退款单号，商户系统内部唯一
	TotalFee      int    // 订单金额，单位为分，只能为整数
	RefundFee     int    // 退款总金额，订单总金额，单位为分，只能为整数
	RefundFeeType string // 货币类型
	OpUserID      string // 操作员帐号, 默认为商户号
	RefundAccount string // 退款资金来源 REFUND_SOURCE_RECHARGE_FUNDS

	p *pay.Pay
}

// New 声明新的 Refund 实例
func New(p *pay.Pay) *Refund {
	return &Refund{
		p: p,
	}
}

func (r *Refund) params() (map[string]string, error) {
	return map[string]string{
		"device_info":     r.DeviceInfo,
		"sign_type":       r.SignType,
		"out_trade_no":    r.OutTradeNO,
		"transaction_id":  r.TransactionID,
		"out_refund_no":   r.OutRefundNO,
		"total_fee":       strconv.Itoa(r.TotalFee),
		"refund_fee":      strconv.Itoa(r.RefundFee),
		"refund_fee_type": r.RefundFeeType,
		"op_user_id":      r.OpUserID,
		"refund_account":  r.RefundAccount,
	}, nil
}

// Do 执行退款操作
func (r *Refund) Do() (*Return, error) {
	params, err := r.params()
	if err != nil {
		return nil, err
	}
	m, err := r.p.Refund(params)
	if err != nil {
		return nil, err
	}

	if err = r.p.ValidateAll(m); err != nil {
		return nil, err
	}

	return newReturn(m)
}
