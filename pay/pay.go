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
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/issue9/wechat/pay/internal"
)

// 预定义的错误类型
var (
	ErrInvalidAppid = errors.New("返回的 appid 与当前的不匹配")
	ErrInvalidMchid = errors.New("返回的 mch_id 与当前的不匹配")
	ErrInvalidSign  = errors.New("不存在签名或是签名无法验证")
)

// Pay 支付的基本配置
type Pay struct {
	mchID  string
	appID  string
	apiKey string
	client *http.Client
}

// New 声明一个新的 *Pay 实例
func New(mchid, appid, apikey string, client *http.Client) *Pay {
	if client == nil {
		client = http.DefaultClient
	}
	return &Pay{
		mchID:  mchid,
		appID:  appid,
		apiKey: apikey,
		client: client,
	}
}

// MchID 获取商户 ID
func (p *Pay) MchID() string {
	return p.mchID
}

// AppID 获取 appid
func (p *Pay) AppID() string {
	return p.appID
}

// APIKey 获取 apikey
func (p *Pay) APIKey() string {
	return p.apiKey
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
// 比如：若已经指定了 appid，会不会使用 pay.AppID；
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
	if resp.StatusCode > 399 {
		return nil, fmt.Errorf("微信服务端返回[%v]状态码", resp.StatusCode)
	}
	defer resp.Body.Close()

	return internal.MapFromXMLReader(resp.Body)
}

// UnifiedOrder 执行统一下单
func (p *Pay) UnifiedOrder(params map[string]string) (map[string]string, error) {
	return p.Post(UnifiedOrderURL, params)
}

// OrderQuery 订单查询
func (p *Pay) OrderQuery(params map[string]string) (map[string]string, error) {
	return p.Post(OrderQueryURL, params)
}

// CloseOrder 关闭订单
func (p *Pay) CloseOrder(params map[string]string) (map[string]string, error) {
	return p.Post(CloseOrderURL, params)
}

// Refund 退款
func (p *Pay) Refund(params map[string]string) (map[string]string, error) {
	return p.Post(RefundURL, params)
}

// RefundQuery 退款查询
func (p *Pay) RefundQuery(params map[string]string) (map[string]string, error) {
	return p.Post(RefundQueryURL, params)
}

// DownloadBill 下载对账单
func (p *Pay) DownloadBill(params map[string]string) (map[string]string, error) {
	return p.Post(DownloadBillURL, params)
}

// Report 主动上报接口
func (p *Pay) Report(params map[string]string) (map[string]string, error) {
	return p.Post(ReportURL, params)
}

// ValidateReturn 验证从微信端返回的数据，仅验证 return_code
func (p *Pay) ValidateReturn(params map[string]string) error {
	if params["return_code"] != Success {
		return errors.New(params["return_msg"])
	}

	return nil
}

// ValidateResult 验证从微信端返回的数据，同时验证 return_code 和 result_code
func (p *Pay) ValidateResult(params map[string]string) error {
	if err := p.ValidateReturn(params); err != nil {
		return err
	}

	if params["result_code"] != Success {
		return errors.New(params["err_code"] + "-" + params["err_code_des"])
	}

	return nil
}

// ValidateSign  验证从微信端返回的数据，同时验证 ValidateResult 和 签名
func (p *Pay) ValidateSign(singType string, params map[string]string) error {
	if err := p.ValidateResult(params); err != nil {
		return err
	}

	sign := params["sign"]
	if sign == "" {
		return ErrInvalidSign
	}

	if sign != Sign(p.apiKey, singType, params) {
		return ErrInvalidSign
	}

	return nil
}

// ValidateAll 验证 ValidateSign 和 appid 及 mchid 是否匹配
func (p *Pay) ValidateAll(signType string, params map[string]string) error {
	if err := p.ValidateSign(signType, params); err != nil {
		return err
	}

	if params["mch_id"] != p.mchID {
		return ErrInvalidMchid
	}

	if params["appid"] != p.appID {
		return ErrInvalidAppid
	}

	return nil
}

// Sign 获取签名字符串
func (p *Pay) Sign(signType string, params map[string]string) string {
	return Sign(p.APIKey(), signType, params)
}

// 将 map 转换成 xml，并写入到 buf
func (p *Pay) map2XML(params map[string]string, buf *bytes.Buffer) error {
	if params["appid"] == "" {
		params["appid"] = p.appID
	}

	if params["mch_id"] == "" {
		params["mch_id"] = p.mchID
	}

	if params["nonce_str"] == "" {
		params["nonce_str"] = NonceString()
	}

	if params["sign"] == "" {
		params["sign"] = Sign(p.apiKey, params["sign_type"], params)
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
