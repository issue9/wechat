// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package template 模板消息管理
package template

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/issue9/wechat/mp/common/token"
)

// Send 发送模板信息
func Send(srv token.Server, to, id, url string, data Data) error {
	obj := &struct {
		To   string        `json:"touser"`
		ID   string        `json:"template_id"`
		URL  string        `json:"url"`
		Data map[string]KV `json:"data"`
	}{
		To:   to,
		ID:   id,
		URL:  url,
		Data: data,
	}

	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	/* 处理返回的信息 */

	url = token.URL(srv, "cgi-bin/message/template/send", nil)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rslt := &Result{}
	if err = json.Unmarshal(respData, rslt); err != nil {
		return err
	}
	if rslt.IsOK() {
		return nil
	}
	return rslt
}
