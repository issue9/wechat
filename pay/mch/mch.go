// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mch

// Mch 支付接口中，涉及到商户的一些基本配置
type Mch struct {
	ID       string
	AppID    string
	APIKey   string
	SignType string
}
