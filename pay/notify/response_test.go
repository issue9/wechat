// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package notify

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

func TestResponse(t *testing.T) {
	a := assert.New(t)
	str := `<xml><return_code>code</return_code><return_msg>msg</return_msg></xml>`

	resp := &Response{
		Code:    "code",
		Message: "msg",
	}
	data, err := xml.Marshal(resp)
	a.NotError(err).Equal(string(data), str)

	resp.XMLName.Local = "local"
	data, err = xml.Marshal(resp)
	a.NotError(err).Equal(string(data), str)
}
