// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

// 返回状态的值
const (
	Success = "SUCCESS"
	Fail    = "FAIL"
)

// DateFormat 日期格式
const DateFormat = "20060102150405"

// 接口地址
const (
	UnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	OrderQueryURL   = "https://api.mch.weixin.qq.com/pay/orderquery"
	CloseOrderURL   = "https://api.mch.weixin.qq.com/pay/closeorder"
	RefundURL       = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	RefundQueryURL  = "https://api.mch.weixin.qq.com/pay/refundquery"
	DownloadBillURL = "https://api.mch.weixin.qq.com/pay/downloadbill"
	ReportURL       = "https://api.mch.weixin.qq.com/payitil/report"
)

// 交易类型
const (
	TradeTypeJSAPI  = "JSAPI"  // 公众号
	TradeTypeNative = "NATIVE" // 扫码
	TradeTypeApp    = "APP"
)

// 签名的类型
const (
	SignTypeMD5        = "MD5"
	SignTypeHmacSha256 = "HMAC-SHA256"
)

// 退款奖金来源
const (
	RefundSourceRechargeFunds  = "REFUND_SOURCE_UNSETTLED_FUNDS" // 可用余额退款
	RefundSourceUnsettledFunds = "REFUND_SOURCE_UNSETTLED_FUNDS" // 未结算资金退款（默认使用未结算资金退款）
)

// 退款渠道
const (
	RefundChannelOriginal = "ORIGINAL" // 原路退款
	RefundChannelBalance  = "BALANCE"  // 退回到余额
)

// 代金券类型
const (
	CouponTypeCash   = "CASH"    // 充值代金券
	CouponTypeNoCash = "NO_CASH" // 非充值代金券
)
