// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package template

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/issue9/wechat/mp/common/result"
	"github.com/issue9/wechat/mp/common/token"
)

// List 模板列表
type List struct {
	result.Result
	List []*Template `json:"template_list"`
}

// Template 模板
type Template struct {
	ID              string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

// Templates 获取模板列表
func Templates(srv token.Server) (*List, error) {
	url := token.URL(srv, "cgi-bin/template/get_all_private_template", nil)
	resp, err := http.Get(url)
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
