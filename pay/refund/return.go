// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package refund

import (
	"fmt"

	"github.com/issue9/wechat/internal/xxml"
	"github.com/issue9/wechat/pay"
)

// Return 退款的申请的返回值
type Return struct {
	DeviceInfo          string `xml:"device_info"`           // 设备号
	TransactionID       string `xml:"transaction_id"`        // 微信订单号
	OutTradeNO          string `xml:"out_trade_no"`          // 商户订单号
	OutRefundNO         string `xml:"out_refund_no"`         // 商户退款单号
	RefundID            string `xml:"refund_id"`             // 微信退款单号
	RefundChannel       string `xml:"refund_channel"`        // 退款渠道
	RefundFee           int    `xml:"refund_fee"`            // 退款金额
	SettlementRefundFee int    `xml:"settlement_refund_fee"` // 应结退款金额
	TotalFee            int    `xml:"total_fee"`             // 订单总金额
	SettlementTotalFee  int    `xml:"settlement_total_fee"`  // 应结订单金额
	FeeType             string `xml:"fee_type"`              // 订单金额货币类型
	CashFee             int    `xml:"cash_fee"`              // 现金支付金额
	CashFeeType         string `xml:"cash_fee_type"`         // 现金支付币种
	CashRefundFee       int    `xml:"cash_refund_fee"`       // 现金退款金额
	CouponRefundFee     int    `xml:"coupon_refund_fee"`     // 代金券退款总金额
	CouponRefundCount   int    `xml:"coupon_refund_count"`   // 退款代金券使用数量

	Coupons []*pay.Coupon
}

func newReturn(params map[string]string) (*Return, error) {
	ret := &Return{}
	if err := xxml.Map2XMLObj(params, ret); err != nil {
		return nil, err
	}

	coupons, err := pay.GetCoupons(params)
	if err != nil {
		return nil, err
	}

	if ret.CouponRefundCount != len(coupons) {
		return nil, fmt.Errorf("返回的代金券数量[%v]和实际的数量[%v]不相符", ret.CouponRefundCount, len(coupons))
	}

	ret.Coupons = coupons

	return ret, nil
}
