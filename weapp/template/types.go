// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package template

import "github.com/issue9/wechat/mp/common/result"

type limit struct {
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

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
