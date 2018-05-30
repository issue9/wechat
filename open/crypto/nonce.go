// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package crypto

import (
	"github.com/issue9/rands"
)

var randstr = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func nonce() []byte {
	return rands.Bytes(16, 17, randstr)
}

func nonceString() string {
	return rands.String(16, 17, randstr)
}
