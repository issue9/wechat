// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// Pay 支付的基本配置
type Pay struct {
	MchID  string
	AppID  string
	APIKey string
	client *http.Client
}

// New 声明一个新的 *Pay 实例
func New(mchid, appid, apikey string, client *http.Client) *Pay {
	if client == nil {
		client = http.DefaultClient
	}
	return &Pay{
		MchID:  mchid,
		AppID:  appid,
		APIKey: apikey,
		client: client,
	}
}

// NewTLSPay 声明一个带证书的支付实例
func NewTLSPay(mchid, appid, apikey, certPath, keyPath, rootCAPath string) (*Pay, error) {
	client, err := newTLSClient(certPath, keyPath, rootCAPath)
	if err != nil {
		return nil, err
	}

	return New(mchid, appid, apikey, client), nil
}

// Post 发送请求，会优先使用 params 中的相关参数。
// 比如：若已经指定了 appid，会不会使用 conf.AppID；
// 若使用了 sign，则不会再计算 sign 值。
func (p *Pay) Post(url string, params map[string]string) (map[string]string, error) {
	buf := new(bytes.Buffer)
	if err := p.map2XML(params, buf); err != nil {
		return nil, err
	}
	resp, err := p.client.Post(url, "application/xml", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return mapFromReader(resp.Body)
}

// ValidateReturn 仅验证 return_code
func (p *Pay) ValidateReturn(params map[string]string) error {
	if params["return_code"] != Success {
		return &ReturnError{
			Code:    params["return_code"],
			Message: params["return_msg"],
		}
	}

	return nil
}

// ValidateResult 同时验证 return_code 和 result_code
func (p *Pay) ValidateResult(params map[string]string) error {
	if err := p.ValidateReturn(params); err != nil {
		return err
	}

	if params["result_code"] != Success {
		return &ReturnError{
			Code:    params["err_code"],
			Message: params["err_code_des"],
		}
	}

	return nil
}

// ValidateSign 同时验证 ValidateResult 和 签名
func (p *Pay) ValidateSign(params map[string]string) error {
	if err := p.ValidateResult(params); err != nil {
		return err
	}

	sign1 := params["sign"]
	if sign1 == "" {
		return errors.New("不存在 sign 字段")
	}

	if sign1 != Sign(p.APIKey, params) {
		return errors.New("签名验证无法通过")
	}

	return nil
}

// ValidateAll 验证 ValidateSign 和 appid 及 mchid 是否匹配
func (p *Pay) ValidateAll(params map[string]string) error {
	if err := p.ValidateSign(params); err != nil {
		return err
	}

	if params["mch_id"] != p.MchID {
		return errors.New("mch_id 不匹配")
	}

	if params["appid"] != p.AppID {
		return errors.New("appid 不匹配")
	}

	return nil
}

func (p *Pay) map2XML(params map[string]string, buf *bytes.Buffer) error {
	if params["appid"] == "" {
		params["appid"] = p.AppID
	}

	if params["mch_id"] == "" {
		params["mch_id"] = p.MchID
	}

	if params["nonce_str"] == "" {
		params["nonce_str"] = NonceString()
	}

	if params["sign"] == "" {
		params["sign"] = Sign(p.APIKey, params)
	}

	buf.WriteString("<xml>")
	for k, v := range params {
		if v == "" {
			continue
		}

		buf.WriteByte('<')
		buf.WriteString(k)
		buf.WriteByte('>')

		if err := xml.EscapeText(buf, []byte(v)); err != nil {
			return err
		}

		buf.WriteString("</")
		buf.WriteString(k)
		buf.WriteByte('>')
	}
	buf.WriteString("</xml>")

	return nil
}

func mapFromReader(r io.Reader) (map[string]string, error) {
	ret := make(map[string]string, 10)
	d := xml.NewDecoder(r)
	for token, err := d.Token(); true; token, err = d.Token() {
		if err != nil {
			return nil, err
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

	return ret, nil
}

// 获取一个带安全证书的 http.Client 实例
func newTLSClient(cert, key, root string) (*http.Client, error) {
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	r, err := ioutil.ReadFile(root)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(r)

	conf := &tls.Config{
		Certificates: []tls.Certificate{c},
		RootCAs:      pool,
	}

	return &http.Client{
		Transport: &http.Transport{TLSClientConfig: conf},
	}, nil
}
