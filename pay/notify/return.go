// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package notify

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
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

// From 从 r 中读取数据到 *Return 中。
func newReturn(params map[string]string) (*Return, error) {
	ret := &Return{}
	if err := internal.Map2XMLObj(params, ret); err != nil {
		return nil, err
	}

	for name, val := range params {
		switch {
		case strings.HasPrefix(name, "coupon_id_"):
			index, err := getCouponIndex(name, "coupon_id_")
			if err != nil {
				return nil, err
			}

			id, err := strconv.Atoi(val)
			if err != nil {
				return err
			}

			if index >= len(ret.Coupons) { // 不存在
				ret.Coupons = append(ret.Coupons, &pay.Coupon{
					ID: id,
				})
				break
			}
			ret.Coupons[index].ID = id
		case strings.HasPrefix(name, "coupon_type_"):
			index, err := getCouponIndex(name, "coupon_type_")
			if err != nil {
				return nil, err
			}

			if index >= len(ret.Coupons) { // 不存在
				ret.Coupons = append(ret.Coupons, &pay.Coupon{
					Type: val,
				})
				break
			}
			ret.Coupons[index].Type = val
		case strings.HasPrefix(name, "coupon_fee_"):
			index, err := getCouponIndex(name, "coupon_fee_")
			if err != nil {
				return nil, err
			}

			fee, err := strconv.Atoi(string(val))
			if err != nil {
				return nil, err
			}

			if index >= len(ret.Coupons) { // 不存在
				ret.Coupons = append(ret.Coupons, &pay.Coupon{
					Fee: fee,
				})
				break
			}
			ret.Coupons[index].Fee = fee
		} // ned switch
	} // end for

	if ret.CouponCount != len(ret.Coupons) {
		return nil, fmt.Errorf("返回的代金券数量[%v]和实际的数量[%v]不相符", ret.CouponCount, len(ret.Coupons))
	}

	// 转换时间值
	end, err := time.Parse("20060102150405", ret.TimeEnd)
	if err != nil {
		return nil, err
	}
	ret.end = end

	return ret, nil
}

// 将 Return 各个字段以 xml 标签中的值进行索引，方便查找。
func values(ret *Return) (map[string]reflect.Value, error) {
	v := reflect.ValueOf(ret).Elem()
	t := v.Type()
	values := make(map[string]reflect.Value, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("xml")
		values[tag] = v.Field(i)
	}

	return values, nil
}

// 获取代金券的索引值
// 比如从 coupon_type_1 获取 1
func getCouponIndex(name, prefix string) (int, error) {
	str := strings.TrimPrefix(name, prefix)
	index, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return index, nil
}

// Read 从 r 读取内容，并尝试转换成 Return 实例
func Read(r io.Reader) (*Return, error) {
	ret := &Return{Coupons: []*pay.Coupon{}}
	err := ret.From(r)
	return ret, err
}
