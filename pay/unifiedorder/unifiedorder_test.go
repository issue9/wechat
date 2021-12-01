// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unifiedorder

import (
	"testing"

	"github.com/issue9/assert/v2"
)

func TestOrder_limitPay(t *testing.T) {
	a := assert.New(t, false)

	o := &Order{}
	a.Equal(o.limitPay(), limitPayNoCredit)

	o.Credit = true
	a.Equal(o.limitPay(), "")
}

func TestOrder_totalFee(t *testing.T) {
	a := assert.New(t, false)

	o := &Order{
		TotalFee: 500,
	}
	fee, err := o.totalFee()
	a.NotError(err).Equal(fee, 500)

	// 刚好相等
	o.Goods(&Good{Price: 50, Quantity: 5}, &Good{Price: 50, Quantity: 5})
	fee, err = o.totalFee()
	a.NotError(err).Equal(fee, 500)

	o.Goods(&Good{Price: 50, Quantity: 5})
	fee, err = o.totalFee()
	a.Error(err).Equal(fee, 0)
}
