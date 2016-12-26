// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

import (
	"strconv"
	"strings"
)

// Coupon 代金券
type Coupon struct {
	ID   int    // 代金券ID
	Type string // 代金券类型，CASH--充值代金券、NO_CASH--非充值代金券
	Fee  int    // 单个代金券支付金额
}

// GetCoupons 从 params 获取所有的代金券信息
func GetCoupons(params map[string]string) ([]*Coupon, error) {
	coupons := []*Coupon{}

	for name, val := range params {
		switch {
		case strings.HasPrefix(name, "coupon_id_"):
			index, err := getCouponIndex(name, "coupon_id_")
			if err != nil {
				return nil, err
			}

			id, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}

			if index >= len(coupons) { // 不存在
				coupons = append(coupons, &Coupon{
					ID: id,
				})
				break
			}
			coupons[index].ID = id
		case strings.HasPrefix(name, "coupon_type_"):
			index, err := getCouponIndex(name, "coupon_type_")
			if err != nil {
				return nil, err
			}

			if index >= len(coupons) { // 不存在
				coupons = append(coupons, &Coupon{
					Type: val,
				})
				break
			}
			coupons[index].Type = val
		case strings.HasPrefix(name, "coupon_fee_"):
			index, err := getCouponIndex(name, "coupon_fee_")
			if err != nil {
				return nil, err
			}

			fee, err := strconv.Atoi(string(val))
			if err != nil {
				return nil, err
			}

			if index >= len(coupons) { // 不存在
				coupons = append(coupons, &Coupon{
					Fee: fee,
				})
				break
			}
			coupons[index].Fee = fee
		} // ned switch
	} // end for

	return coupons, nil
}

// 获取代金券的索引值，比如从 coupon_type_1 获取 1
func getCouponIndex(name, prefix string) (int, error) {
	str := strings.TrimPrefix(name, prefix)
	index, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return index, nil
}
