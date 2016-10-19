// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package template

import (
	"github.com/issue9/wechat/mp/common/result"
)

// Result 表示消息推送之后的数据
type Result struct {
	result.Result
	MsgID int64 `json:"msgid"`
}

type SendTemplate struct {
	To   string         `json:"touser"`
	ID   string         `json:"template_id"`
	URL  string         `json:""url"`
	Data map[string]*KV `json:"data"`
}

type KV struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
