// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package weapp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 几个错误信息
var (
	ErrInvalidSessionKey = errors.New("无效的 sessionkey")
	ErrInvalidInitVector = errors.New("无效的 Initization vector")
)

// Decode 解析加密数据
func Decode(appid, sessionkey, data, iv string) ([]byte, error) {
	if len(sessionkey) != 24 {
		return nil, ErrInvalidSessionKey
	}
	aeskey, err := base64.StdEncoding.DecodeString(sessionkey)
	if err != nil {
		return nil, err
	}

	if len(iv) != 24 {
		return nil, ErrInvalidInitVector
	}
	aesiv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}

	text := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, aesiv)
	mode.CryptBlocks(text, cipherText)

	return pkcs7UnPadding(text), nil
}

func pkcs7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return plantText[:(length - unPadding)]
}
