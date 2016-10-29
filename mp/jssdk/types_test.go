// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package jssdk

import (
	"testing"

	"github.com/issue9/assert"
)

func TestUserInfo_HeadImageURL(t *testing.T) {
	a := assert.New(t)

	info := &UserInfo{
		HeadImgURL: "https://test.com/abc.png/46",
	}

	url, err := info.HeadImageURL(0)
	a.NotError(err).Equal(url, "https://test.com/abc.png/0")

	url, err = info.HeadImageURL(46)
	a.NotError(err).Equal(url, "https://test.com/abc.png/46")

	// 无效的 size
	url, err = info.HeadImageURL(1924)
	a.Error(err).Equal(url, "")
}
