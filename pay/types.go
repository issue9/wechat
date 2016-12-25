// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

// Coupon 代金券
type Coupon struct {
	ID   int    // 代金券ID
	Type string // 代金券类型，CASH--充值代金券、NO_CASH--非充值代金券
	Fee  int    // 单个代金券支付金额
}
