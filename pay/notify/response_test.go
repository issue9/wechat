// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package notify

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/wechat/common"
)

func TestResponse(t *testing.T) {
	a := assert.New(t)
	str := `<xml><return_code><![CDATA[code]]></return_code><return_msg><![CDATA[msg]]></return_msg></xml>`

	resp := &Response{
		Code:    common.CData{Text: "code"},
		Message: common.CData{Text: "msg"},
	}
	data, err := xml.Marshal(resp)
	a.NotError(err).Equal(string(data), str)

	resp.XMLName.Local = "local"
	data, err = xml.Marshal(resp)
	a.NotError(err).Equal(string(data), str)
}
