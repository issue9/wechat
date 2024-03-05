// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package template

import "github.com/issue9/wechat/common"

// Result 表示消息推送之后的数据
type Result struct {
	common.Result
	MsgID int64 `json:"msgid"`
}

// Data 表示发送模板时的数据内容
type Data map[string]KV

type KV struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
