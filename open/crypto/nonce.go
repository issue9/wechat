// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package crypto

import "github.com/issue9/rands/v2"

func nonce() []byte {
	return rands.Bytes(16, 17, rands.AlphaNumber())
}

func nonceString() string {
	return rands.String(16, 17, rands.AlphaNumber())
}
