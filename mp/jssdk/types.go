// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package jssdk

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// AccessToken 表示 jssdk 中返回的 access_token 结构体
//
// NOTE: 与 mp/common/token.AccessToken 同名，但结构不同。
type AccessToken struct {
	AccessToken  string        `json:"access_token"`
	ExpiresIn    time.Duration `json:"expires_in"`
	Created      time.Time     `json:"-"` // 该 access_token 的获取时间
	RefreshToken string        `json:"refresh_token"`
	OpenID       string        `json:"openid"`
	Scope        string        `json:"scope"`
}

// UserInfo 查询用户信息接口返回的数据
type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        string   `json:"sex"`        // 性别，1男，2女，0未知
	Province   string   `json:"province"`   // 省份
	City       string   `json:"city"`       // 城市
	Country    string   `json:"country"`    // 国家
	HeadImgURL string   `json:"headimgurl"` // 头像地址
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid,omitempty"`
}

// HeadImageURL 相对于 HeadImgURL 的好处是，可以指定图片的尺寸。
// 可以是 46，64，96，132，0(表示640)
func (u *UserInfo) HeadImageURL(size int) (string, error) {
	if size != 0 && size != 46 && size != 64 && size != 96 && size != 132 {
		return "", errors.New("取值不正确")
	}

	index := strings.LastIndexByte(u.HeadImgURL, '/')
	return u.HeadImgURL[:index] + "/" + strconv.Itoa(size), nil
}
