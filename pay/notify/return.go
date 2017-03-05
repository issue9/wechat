// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package notify

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/issue9/wechat/pay"
	"github.com/issue9/wechat/pay/internal"
)

// Return 微信返回的信息结构
type Return struct {
	DeviceInfo         string `xml:"device_info"`          // 设备号
	OpenID             string `xml:"openid"`               // 用户标识
	IsSubscribe        string `xml:"is_subscribe"`         // 是否关注公众账号，Y-关注，N-未关注
	TradeType          string `xml:"trade_type"`           // 交易类型，JSAPI、NATIVE、APP
	BankType           string `xml:"bank_type"`            // 付款银行
	TotalFee           int    `xml:"total_fee"`            // 订单金额，单位为分
	SettlementTotalFee int    `xml:"settlement_total_fee"` // 应结订单金额
	FeeType            string `xml:"fee_type"`             // 货币种类，符合 ISO4217 标准的三位字母代码，默认：CNY
	CashFee            int    `xml:"cash_fee"`             // 现金支付金额
	CashFeeType        string `xml:"cash_fee_type"`        // 现金支付货币类型，符合 ISO4217 标准的三位字母代码
	CouponFee          int    `xml:"coupon_fee"`           // 总代金券金额
	CouponCount        int    `xml:"coupon_count"`         // 代金券使用数量
	TransactionID      string `xml:"transaction_id"`       // 微信支付订单号
	OutTradeNO         string `xml:"out_trade_no"`         // 商户订单号
	Attach             string `xml:"attach"`               // 商家数据包
	TimeEnd            string `xml:"time_end"`             // 支付完成时间，格式为yyyyMMddHHmmss

	Coupons []*pay.Coupon
	end     time.Time
}

// End 返回 TimeEnd 的 time.Time 格式数据
func (ret *Return) End() time.Time {
	return ret.end
}

// Subscribed 当前用户是否已经关注公众账号
func (ret *Return) Subscribed() bool {
	return ret.IsSubscribe == "Y"
}

// Read 从 r 读取内容，并尝试转换成 Return 实例
func Read(p *pay.Pay, r io.Reader) (*Return, error) {
	params, err := internal.MapFromXMLReader(r)
	if err != nil {
		return nil, err
	}

	if len(params) == 0 {
		return nil, errors.New("未读取到任何数据")
	}

	if err = p.ValidateAll(params["sign_type"], params); err != nil {
		return nil, err
	}

	ret := &Return{}
	err = internal.Map2XMLObj(params, ret)
	if err != nil {
		return nil, err
	}

	coupons, err := pay.GetCoupons(params)
	if err != nil {
		return nil, err
	}

	if ret.CouponCount != len(ret.Coupons) {
		return nil, fmt.Errorf("返回的代金券数量[%v]和实际的数量[%v]不相符", ret.CouponCount, len(ret.Coupons))
	}
	ret.Coupons = coupons

	// 转换时间值
	end, err := time.Parse(pay.DateFormat, ret.TimeEnd)
	if err != nil {
		return nil, err
	}
	ret.end = end.Add(-pay.TimeFixed) // 返回的为区八区，需要减去，才是 UTC

	return ret, nil
}
