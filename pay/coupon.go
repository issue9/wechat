// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package pay

import (
	"sort"
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
	coupons := map[int]*Coupon{}

LOOP:
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

			c, found := coupons[index]
			if !found { // 不存在
				coupons[index] = &Coupon{
					ID: id,
				}
				continue LOOP
			}
			c.ID = id
		case strings.HasPrefix(name, "coupon_type_"):
			index, err := getCouponIndex(name, "coupon_type_")
			if err != nil {
				return nil, err
			}

			c, found := coupons[index]
			if !found { // 不存在
				coupons[index] = &Coupon{
					Type: val,
				}
				continue LOOP
			}
			c.Type = val
		case strings.HasPrefix(name, "coupon_fee_"):
			index, err := getCouponIndex(name, "coupon_fee_")
			if err != nil {
				return nil, err
			}

			fee, err := strconv.Atoi(string(val))
			if err != nil {
				return nil, err
			}

			c, found := coupons[index]
			if !found { // 不存在
				coupons[index] = &Coupon{
					Fee: fee,
				}
				continue LOOP
			}
			c.Fee = fee
		} // ned switch
	} // end for

	ret := make([]*Coupon, 0, len(coupons))
	for _, c := range coupons {
		ret = append(ret, c)
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].ID < ret[j].ID
	})

	return ret, nil
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
