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

// Data 表示发送模板时的数据内容
type Data map[string]KV

type KV struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
