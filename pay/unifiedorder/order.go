// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unifiedorder

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/issue9/wechat/pay"
)

const format = "20060102150405"

// Order 表示一条完整的订单数据
type Order struct {
	DeviceInfo     string        // 设备号
	SignType       string        // 签名类型
	Body           string        // 商品描述
	Attach         string        // 附加数据
	TradeNO        string        // 商户订单号
	FeeType        string        // 货币类型,CNY
	TotalFee       int           // 总金额
	SpbillCreateIP string        // 终端IP
	Start          time.Time     // 交易起始时间
	ExpireIn       time.Duration // 交易结束时间
	Tag            string        // 商品标记
	NotifyURL      string        // 通知地址
	TradeType      string        // 交易类型
	ProductID      string        // 商品ID
	LimitPay       string        // 指定支付方式
	OpenID         string        // 用户标识

	conf  *pay.Config
	goods []*Good
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

// NewOrder 声明一个新的 Order 实例
func NewOrder(conf *pay.Config) *Order {
	return &Order{
		conf:  conf,
		goods: []*Good{},
	}
}

// New 根据现在有的参数，生成一个新的订单内容。
// 会尽量重用现在有的参数。
func (o *Order) New() *Order {
	o.TradeNO = ""
	o.TotalFee = 0
	o.goods = o.goods[:0]
	return nil
}

// Goods 为当前订单添加一条订的物品记录
func (o *Order) Goods(goods ...*Good) {
	o.goods = append(o.goods, goods...)
}

// Params pay.Paramser 接口
func (o *Order) Params() (map[string]string, error) {
	detail, err := json.Marshal(o.goods)
	if err != nil {
		return nil, err
	}

	var start, end string
	if !o.Start.IsZero() {
		start = o.Start.Format(format)
		if o.ExpireIn > 0 {
			end = o.Start.Add(o.ExpireIn).Format(format)
		}
	}

	// 若存在商品信息，则价格由商品信息决定
	if len(o.goods) > 0 {
		o.TotalFee = 0
		for _, g := range o.goods {
			o.TotalFee += g.Quantity * g.Price
		}
	}

	return map[string]string{
		"appid":            o.conf.AppID,
		"mch_id":           o.conf.MchID,
		"device_info":      o.DeviceInfo,
		"nonce_str":        "", // 为空，由 pay.Post 自行计算
		"sign":             "",
		"sign_type":        o.conf.SignType,
		"body":             o.Body,
		"detail":           string(detail),
		"attach":           o.Attach,
		"out_trade_no":     o.TradeNO,
		"fee_type":         o.FeeType,
		"total_fee":        strconv.Itoa(o.TotalFee),
		"spbill_create_ip": o.SpbillCreateIP,
		"time_start":       start,
		"time_expire":      end,
		"tag":              o.Tag,
		"notify_url":       o.NotifyURL,
		"trade_type":       o.TradeType,
		"product_id":       o.ProductID,
		"limit_pay":        o.LimitPay,
		"openid":           o.OpenID,
	}, nil
}

// Pay 下单
//
// Example:
//  o := NewOrder(mch)
//  o.Pay()
//
//  o = o.New() // 新的参数 o，会重用大部分 o 中的数据。
//  o.Body = "body"
//  o.Pay()
func (o *Order) Pay() (*Return, error) {
	ret := &Return{}
	err := pay.Post(o.conf, pay.UnifiedOrderURL, o, ret)
	return ret, err
}
