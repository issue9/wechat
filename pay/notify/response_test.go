// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package notify

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert/v4"
	"github.com/issue9/wechat/internal"
)

func TestResponse(t *testing.T) {
	a := assert.New(t, false)
	str := `<xml><return_code><![CDATA[code]]></return_code><return_msg><![CDATA[msg]]></return_msg></xml>`

	resp := &Response{
		Code:    internal.CData{Text: "code"},
		Message: internal.CData{Text: "msg"},
	}
	data, err := xml.Marshal(resp)
	a.NotError(err).Equal(string(data), str)

	resp.XMLName.Local = "local"
	data, err = xml.Marshal(resp)
	a.NotError(err).Equal(string(data), str)
}
