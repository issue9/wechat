// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

import (
	"encoding/xml"
	"io"
)

// Config 支付接口中的一些基本配置
type Config struct {
	MchID    string
	AppID    string
	APIKey   string
	SignType string // 签名类型，可以是 MD5 或是 HMAC-SHA256
	NonceStr string // 随机字符串
}

// Returner 用于描述从服务端返回的数据结构
type Returner interface {
	// 是否正常返回
	OK() bool

	// 是否交易完成
	Success() bool

	// 从 io.Reader 初始化当前的 Returner
	From(io.Reader) error
}

// Paramser 描述了向服务器提交数据的基本接口。
// 需要能将对象转换成 map[string]string 形式的接口
type Paramser interface {
	// 将提交的数据转换成 map[string]string 格式。
	Params() (map[string]string, error)
}

// Params 是 Paramser 接口的默认实现
type Params map[string]string

func (p Params) Params() (map[string]string, error) {
	return p, nil
}

// Return 是 Returner 接口的默认实现
type Return map[string]string

func (ret Return) OK() bool {
	return ret["return_code"] == Success
}

func (ret Return) Success() bool {
	return ret.OK() && ret["result_code"] == Success
}

func (ret Return) From(r io.Reader) error {
	d := xml.NewDecoder(r)
	for token, err := d.Token(); true; token, err = d.Token() {
		if err != nil {
			return err
		}

		var key, val string
		switch t := token.(type) {
		case xml.StartElement:
			key = t.Name.Local
		case xml.CharData:
			val = string(t)
		}
		ret[key] = val
	}

	return nil
}
