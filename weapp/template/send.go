// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package template 模板消息管理
package template

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/issue9/wechat/common/token"
)

// Send 发送模板信息
func Send(srv token.Server, to, tplid, page, formid string, data Data) error {
	obj := &struct {
		To   string        `json:"touser"`
		ID   string        `json:"template_id"`
		Page string        `json:"page"`
		Form string        `json:"form_id"`
		Data map[string]KV `json:"data"`
	}{
		To:   to,
		ID:   tplid,
		Page: to,
		Form: formid,
		Data: data,
	}

	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	/* 处理返回的信息 */

	url := token.URL(srv, "cgi-bin/message/wxopen/template/send", nil)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
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
