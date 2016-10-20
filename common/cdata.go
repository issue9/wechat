// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package common

// CData 表示 xml 中的 CDATA 的内容，本身并不能直接以字段的形式
// 出现在结构体中，只能以一个结构包含住。
//
// 具体的使用方式可参考测试代码
type CData struct {
	Text string `xml:",cdata"`
}
