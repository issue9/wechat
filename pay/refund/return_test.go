// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package refund

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/wechat/pay"
)

var _ pay.Returner = &Return{}

func TestReturn_From(t *testing.T) {
	a := assert.New(t)

	xml := `<xml>
   <return_code><![CDATA[SUCCESS]]></return_code>
   <return_msg><![CDATA[OK]]></return_msg>
   <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
   <mch_id><![CDATA[10000100]]></mch_id>
   <nonce_str><![CDATA[NfsMFbUFpdbEhPXP]]></nonce_str>
   <sign><![CDATA[B7274EB9F8925EB93100DD2085FA56C0]]></sign>
   <result_code><![CDATA[SUCCESS]]></result_code>
   <transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>
   <out_trade_no><![CDATA[1415757673]]></out_trade_no>
   <out_refund_no><![CDATA[1415701182]]></out_refund_no>
   <refund_id><![CDATA[2008450740201411110000174436]]></refund_id>
   <refund_channel><![CDATA[]]></refund_channel>
   <refund_fee>1</refund_fee> 
</xml>`

	buf := bytes.NewBufferString(xml)
	ret := &Return{}
	a.NotError(ret.From(buf))
	a.Equal(ret.AppID, "wx2421b1c4370ec43b")
	a.Equal(ret.Sign, "B7274EB9F8925EB93100DD2085FA56C0")
	a.Equal(ret.TotalFee, 0)
	a.True(ret.Success())
}
