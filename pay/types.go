// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

type Returner struct {
}

// ReturnError 表示微信返回内容的错误信息
type ReturnError struct {
	Code    string // 错误代码，可能为空
	Message string
}

func (e *ReturnError) Error() string {
	return e.Message
}

// Coupon 代金券
type Coupon struct {
	ID   int    // 代金券ID
	Type string // 代金券类型，CASH--充值代金券、NO_CASH--非充值代金券
	Fee  int    // 单个代金券支付金额
}
