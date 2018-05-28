// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package weapp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/issue9/wechat/common"
)

// 几个错误信息
var (
	ErrInvalidSessionKey = errors.New("无效的 sessionkey")
	ErrInvalidInitVector = errors.New("无效的 Initization vector")
)

// 加密数据解密之后的部分内容
type encyData struct {
	Watermark *Watermark `json:"watermark"`
}

// Watermark 水印部分的内容
type Watermark struct {
	Appid     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

// Decode 解析加密数据
//
// 返回的 Watermark 可用以验证数据是否正确
func Decode(appid, sessionkey, data, iv string) ([]byte, *Watermark, error) {
	if len(sessionkey) != 24 {
		return nil, nil, ErrInvalidSessionKey
	}
	aeskey, err := base64.StdEncoding.DecodeString(sessionkey)
	if err != nil {
		return nil, nil, err
	}

	if len(iv) != 24 {
		return nil, nil, ErrInvalidInitVector
	}
	aesiv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, nil, err
	}

	mode := cipher.NewCBCDecrypter(block, aesiv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText = common.PKCS7UnPadding(cipherText)
	obj := &encyData{}
	if err = json.Unmarshal(cipherText, obj); err != nil {
		return nil, nil, err
	}

	return cipherText, obj.Watermark, nil
}
