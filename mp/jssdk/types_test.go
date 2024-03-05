// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package jssdk

import (
	"testing"

	"github.com/issue9/assert/v4"
)

func TestUserInfo_HeadImageURL(t *testing.T) {
	a := assert.New(t, false)

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
