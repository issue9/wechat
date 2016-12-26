// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package unifiedorder 执行微信的下单操作。
// 相对于 pay.Post，会便利很多，不需要每次都指定所有参数。
//  p := pay.New(...)
//  oconf := unifiedorder.New(p, pay.TradeTypeJSSDK, ...)
//
//  // 下单支付
//  o := oconf.NewOrder()
//  o.Body = "..."
//  o.Goods(...)
//  o.Pay(...)
//
//  // 另一笔支付操作
//  o = oconf.NewOrder()
//  o.Body = "..."
//  o.Goods(...)
//  o.Pay(...)
package unifiedorder

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/issue9/wechat/pay"
)

const limitPayNoCredit = "no_credit"

// Order 订单数据
//
// 订单数据中，包含了大量相同的数据，这些数据在每次
// 支付时都重新声明是极大的浪费，所以建议先声明一个
// 包含这些公用数据的 Order 实例，之后用 Order.NewOrder()
// 生成一个新的 Order 实例，这个新的实例就继承这些数据。比如：
//  conf := unifiedorder.New(...)
//  conf.DeviceInfo = "WEB"
//  conf.SignType = pay.SignTypeMD5
//
//  o1 := conf.NewOrder() // 继承了 conf.DeviceInfo 等数据
//  o1.Body = "..."
//  o1.Pay(...)
//
//  o2 := conf.NewOrder() // 依然从 conf 继承了 DeviceInfo 等数据
type Order struct {
	pay            *pay.Pay
	DeviceInfo     string        // 设备号
	SignType       string        // 签名类型
	FeeType        string        // 货币类型，默认 CNY
	SpbillCreateIP string        // 终端 IP
	ExpireIn       time.Duration // 交易结束时间
	NotifyURL      string        // 通知地址
	TradeType      string        // 交易类型
	Credit         bool          // 是否允许使用信用卡

	Body       string    // 商品描述
	Attach     string    // 附加数据
	OutTradeNO string    // 商户订单号
	TotalFee   int       // 总金额
	Start      time.Time // 交易起始时间
	Tag        string    // 商品标记
	ProductID  string    // 商品 ID
	OpenID     string    // 用户标识
	goods      []*Good
}

// Good 商品详情
type Good struct {
	ID           string `json:"goods_id"`
	WxpayGoodsID string `json:"wxpay_goods_id,omitempty"`
	Name         string `json:"goods_name"`
	Quantity     int    `json:"quantity"` // 数量
	Price        int    `json:"price"`    // 价格，单位：分
	Category     string `json:"goods_category,omitempty"`
	Body         string `json:"body,omitempty"`
}

// Return 表示统一下单功能的返回值类型。
type Return struct {
	TradeType string
	PrepayID  string
	CodeURL   string
}

// New 新的订单实例
func New(p *pay.Pay) *Order {
	return &Order{
		pay: p,
	}
}

// 获取支付类型
func (o *Order) limitPay() string {
	if !o.Credit {
		return limitPayNoCredit
	}
	return ""
}

// NewOrder 生成一条新的订单
func (o *Order) NewOrder() *Order {
	ret := &Order{}
	*ret = *o

	// 重置数据
	ret.Body = ""
	ret.Attach = ""
	ret.OutTradeNO = ""
	ret.TotalFee = 0
	ret.Start = time.Time{}
	ret.Tag = ""
	ret.ProductID = ""
	ret.OpenID = ""
	ret.goods = ret.goods[:0]

	return ret
}

// Goods 为当前订单添加一条或是多条物品记录
func (o *Order) Goods(goods ...*Good) {
	o.goods = append(o.goods, goods...)
}

// 获取订单的实际金额
func (o *Order) totalFee() (int, error) {
	if len(o.goods) == 0 {
		return o.TotalFee, nil
	}

	totalFee := 0
	for _, g := range o.goods {
		totalFee += g.Quantity * g.Price
	}

	if o.TotalFee > 0 && o.TotalFee != totalFee {
		return 0, errors.New("指定了 TotalFee，但与实际的 goods 计算值不相同")
	}

	return totalFee, nil
}

// 将当前实例转换成 map[string]string 格式
func (o *Order) params() (map[string]string, error) {
	detail, err := json.Marshal(o.goods)
	if err != nil {
		return nil, err
	}

	var start, end string
	if !o.Start.IsZero() {
		start = o.Start.Format(pay.DateFormat)
		if o.ExpireIn > 0 {
			end = o.Start.Add(o.ExpireIn).Format(pay.DateFormat)
		}
	}

	totalFee, err := o.totalFee()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"device_info":      o.DeviceInfo,
		"sign_type":        o.SignType,
		"body":             o.Body,
		"detail":           "<![CDATA[" + string(detail) + "]]>",
		"attach":           o.Attach,
		"out_trade_no":     o.OutTradeNO,
		"fee_type":         o.FeeType,
		"total_fee":        strconv.Itoa(totalFee),
		"spbill_create_ip": o.SpbillCreateIP,
		"time_start":       start,
		"time_expire":      end,
		"tag":              o.Tag,
		"notify_url":       o.NotifyURL,
		"trade_type":       o.TradeType,
		"product_id":       o.ProductID,
		"limit_pay":        o.limitPay(),
		"openid":           o.OpenID,
	}, nil
}

// Pay 下单
//
// Example:
//  uo := &unifiedorder.UnifiedOrder{...}
//
//  o = uo.NewOrder()
//  o.Body = "body"
//  o.Pay()
func (o *Order) Pay() (*Return, error) {
	params, err := o.params()
	if err != nil {
		return nil, err
	}
	m, err := o.pay.UnifiedOrder(params)
	if err != nil {
		return nil, err
	}

	if err = o.pay.ValidateAll(m); err != nil {
		return nil, err
	}

	return &Return{
		TradeType: m["trade_type"],
		PrepayID:  m["prepay_id"],
		CodeURL:   m["code_url"],
	}, nil
}
