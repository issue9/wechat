// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

import (
	"testing"

	"github.com/issue9/assert"
)

func TestGetCouponIndex(t *testing.T) {
	a := assert.New(t)

	test := func(name, prefix string, wont int) {
		num, err := getCouponIndex(name, prefix)
		a.NotError(err).Equal(num, wont)
	}
	test("coupon_id_123", "coupon_id_", 123)
	test("coupon_id_0", "coupon_id_", 0)
	test("coupon_id_-1", "coupon_id_", -1)

	num, err := getCouponIndex("coupon_id_w123", "coupon_id_")
	a.Error(err).Equal(num, 0)
}

func TestGetCoupons(t *testing.T) {
	a := assert.New(t)
	maps := map[string]string{
		"coupon_id_0":   "0",
		"coupon_type_0": "cash",
		"coupon_fee_0":  "10",

		"coupon_type_1": "cash",
		"coupon_id_1":   "1",
		"coupon_fee_1":  "10",

		"coupon_fee_2":  "10",
		"coupon_id_2":   "2",
		"coupon_type_2": "cash",
	}

	coupons, err := GetCoupons(maps)
	a.NotError(err).NotNil(coupons)
	a.Equal(len(coupons), 3)
	a.Equal(coupons[0].ID, 0).
		Equal(coupons[0].Fee, 10).
		Equal(coupons[0].Type, "cash")

	// 确定顺序是否正确
	a.Equal(coupons[1].ID, 1)
	a.Equal(coupons[2].ID, 2)
}
