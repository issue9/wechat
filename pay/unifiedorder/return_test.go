// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unifiedorder

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/wechat/pay"
)

var _ pay.Returner = &Return{}

func TestParseReturn(t *testing.T) {
	a := assert.New(t)

	str := `<xml>
<return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
<mch_id><![CDATA[10000100]]></mch_id>
<nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str>
<openid><![CDATA[oUpF8uMuAJO_M2pxb1Q9zNjWeS6o]]></openid>
<sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign>
<result_code><![CDATA[SUCCESS]]></result_code>
<prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id>
<trade_type><![CDATA[JSAPI]]></trade_type>
</xml>`

	ret := &Return{}
	a.NotError(xml.Unmarshal([]byte(str), ret))

	a.Equal(ret.Code, "SUCCESS").
		Equal(ret.Message, "OK").
		Equal(ret.AppID, "wx2421b1c4370ec43b").
		Equal(ret.MchID, "10000100").
		Equal(ret.NonceStr, "IITRi8Iabbblz1Jc").
		Equal(ret.Sign, "7921E432F65EB8ED0CE9755F0E86D72F").
		Equal(ret.PrepayID, "wx201411101639507cbf6ffd8b0779950874").
		Equal(ret.TradeType, "JSAPI")
}

func TestReturn_IsOK(t *testing.T) {
	a := assert.New(t)

	ret := &Return{
		Code: "123",
	}
	a.False(ret.OK())

	ret.Code = "success"
	a.True(ret.OK())

	ret.Code = "SUCCESS"
	a.True(ret.OK())

	ret.Code = "SuccesS"
	a.True(ret.OK())
}

func TestReturn_IsSuccess(t *testing.T) {
	a := assert.New(t)

	ret := &Return{
		Code:       "Success",
		ResultCode: "123",
	}
	a.False(ret.Success())

	ret.ResultCode = "success"
	a.True(ret.Success())

	ret.ResultCode = "SUCCESS"
	a.True(ret.Success())

	ret.ResultCode = "SuccesS"
	a.True(ret.Success())

	ret.Code = "123"
	a.False(ret.Success())
}
