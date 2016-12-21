// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package notify

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/wechat/pay"
)

var _ pay.Returner = &Return{}

func TestReturn(t *testing.T) {
	a := assert.New(t)

	xml := `<xml>
  <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
  <attach><![CDATA[支付测试]]></attach>
  <bank_type><![CDATA[CFT]]></bank_type>
  <fee_type><![CDATA[CNY]]></fee_type>
  <is_subscribe><![CDATA[Y]]></is_subscribe>
  <mch_id><![CDATA[10000100]]></mch_id>
  <nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
  <openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
  <out_trade_no><![CDATA[1409811653]]></out_trade_no>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
  <sub_mch_id><![CDATA[10000100]]></sub_mch_id><!-- 不存在的项目 -->
  <time_end><![CDATA[20140903131540]]></time_end>
  <total_fee>1</total_fee>
  <trade_type><![CDATA[JSAPI]]></trade_type>
  <transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
  <coupon_count>2</coupon_count>
  <coupon_id_0>10001</coupon_id_0>
  <coupon_type_0>type</coupon_type_0>
  <coupon_fee_0>250</coupon_fee_0>
  <coupon_id_1>10002</coupon_id_1>
  <coupon_type_1>type1</coupon_type_1>
  <coupon_fee_1>251</coupon_fee_1>
</xml>`

	buf := bytes.NewBufferString(xml)
	ret, err := Read(buf)
	a.NotError(err).NotNil(ret)
	a.Equal(ret.AppID, "wx2421b1c4370ec43b")
	a.Equal(ret.Sign, "B552ED6B279343CB493C5DD0D78AB241")
	a.Equal(ret.TotalFee, 1)
	a.True(ret.End().Unix() > 0)
	a.True(ret.Subscribed())
	a.True(ret.Success())
}
