// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package notify

import "encoding/xml"

// Response 向微信返馈的信息
type Response struct {
	XMLName xml.Name `xml:"xml"`
	Code    string   `xml:"return_code"`
	Message string   `xml:"return_msg"`
}
