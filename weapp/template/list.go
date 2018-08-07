// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/issue9/wechat/mp/common/result"
	"github.com/issue9/wechat/mp/common/token"
)

// List 模板列表
type List struct {
	result.Result
	List []*Template `json:"list"`
}

// Template 模板
type Template struct {
	ID      string `json:"template_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Example string `json:"example"`
}

// Templates 获取模板列表
func Templates(srv token.Server, page, count int) (*List, error) {
	url := token.URL(srv, "cgi-bin/wxopen/template/list", nil)

	data, err := json.Marshal(&limit{Offset: page, Count: count})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	l := &List{}
	if err = json.Unmarshal(respData, l); err != nil {
		return nil, err
	}
	if l.IsOK() {
		return l, nil // l 也有可能是个错误
	}
	return nil, &l.Result
}
