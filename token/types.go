// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package token

import "time"

// Server 表示中控服务器接口
type Server interface {
	// 获取中控服务器缓存的 access_token。
	// 若 access_token 已经过期，服务器应该自动更新。
	Token() (*AccessToken, error)

	// 刷新中控服务器的 access_token。
	//
	// 在服务器端应该做好防止客户多次连续调用 Resresh 的可能。
	Refresh() (*AccessToken, error)
}

// AccessToken 用于描述从 https://api.weixin.qq.com/cgi-bin/token 正常返回的数据结构。
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Created     int64  // 该 access_token 的获取时间
}

// IsExpired 该 access_token 是否还在有效期之内。
func (at *AccessToken) IsExpired() bool {
	return time.Now().Unix() > (at.Created + at.ExpiresIn)
}
