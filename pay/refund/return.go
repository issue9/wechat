// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package refund

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/issue9/wechat/pay"
)

// Return 退款的申请的返回值
type Return struct {
	Code    string `xml:"return_code"`
	Message string `xml:"return_msg"`

	ResultCode          string `xml:"result_code"`           // 业务结果
	ErrCode             string `xml:"err_code"`              // 错误代码
	ErrCodeDesc         string `xml:"err_code_des"`          // 错误代码描述
	AppID               string `xml:"appid"`                 // 公众账号ID
	MchID               string `xml:"mch_id"`                // 商户号
	DeviceInfo          string `xml:"device_info"`           // 设备号
	NonceStr            string `xml:"nonce_str"`             // 随机字符串
	Sign                string `xml:"sign"`                  // 签名
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

// OK 通信是否正常，即 Result.Code 是否为 SUCCESS
func (ret *Return) OK() bool {
	return strings.ToUpper(ret.Code) == pay.Success
}

// Success 交易是否正常，即 (Result.Code == SUCCESS) && (Result.ResultCode == SUCCESS)
func (ret *Return) Success() bool {
	return ret.OK() && strings.ToUpper(ret.ResultCode) == pay.Success
}

// From 从 r 中读取数据到 *Return 中。
func (ret *Return) From(r io.Reader) error {
	d := xml.NewDecoder(r)
	values, err := values(ret)
	if err != nil {
		return err
	}

	for token, err := d.Token(); true; token, err = d.Token() {
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		elem, ok := token.(xml.StartElement)
		if !ok || elem.Name.Local == "xml" { // 忽略非 StartElement 和 xml 标签
			continue
		}
		name := elem.Name.Local // xml 元素的名称

		token, err = d.Token()
		if err != nil { // 此处若 err == io.EOF，必须是格式错误，不用专门判断
			return err
		}
		bs, ok := token.(xml.CharData)
		if !ok {
			return fmt.Errorf("无法转换成 xml.CharData")
		}
		val := string(bs) // xml 元素的值
		switch {
		case strings.HasPrefix(name, "coupon_id_"):
			index, err := getCouponIndex(name, "coupon_id_")
			if err != nil {
				return err
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
				return err
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
				return err
			}

			fee, err := strconv.Atoi(string(val))
			if err != nil {
				return err
			}

			if index >= len(ret.Coupons) { // 不存在
				ret.Coupons = append(ret.Coupons, &pay.Coupon{
					Fee: fee,
				})
				break
			}
			ret.Coupons[index].Fee = fee
		default:
			item, found := values[name]
			if !found { // 不存在的字段
				continue
			}

			if item.Kind() == reflect.String {
				item.SetString(val)
			} else if item.Kind() == reflect.Int {
				i, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					return err
				}
				item.SetInt(i)
			}
		} // ned switch
	} // end for

	if ret.CouponRefundCount != len(ret.Coupons) {
		return fmt.Errorf("返回的代金券数量[%v]和实际的数量[%v]不相符", ret.CouponRefundCount, len(ret.Coupons))
	}

	return nil
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
