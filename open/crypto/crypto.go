// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package crypto 消息的加解密功能
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/issue9/wechat/common"
)

// Crypto 加解密功能
type Crypto struct {
	token string
	appid []byte
	key   []byte

	plainlen int
}

// New 声明一个 Crypto 实例
//
// encodingAesKey 不需要结尾的 = 字符
func New(appid, token, encodingAesKey string) (*Crypto, error) {
	if len(encodingAesKey) != 43 {
		return nil, errors.New("无效的参数 encodingAesKey")
	}

	key, err := base64.StdEncoding.DecodeString(encodingAesKey + "=")
	if err != nil {
		return nil, err
	}

	return &Crypto{
		token:    token,
		appid:    []byte(appid),
		key:      key,
		plainlen: 16 + 4 + len(appid),
	}, nil
}

// Encrypt 加密内容 AES_Encrypt[random(16B) + msg_len(4B) + rawXMLMsg + appId]
func (c *Crypto) Encrypt(xmltext []byte) ([]byte, error) {
	text := make([]byte, 0, c.plainlen+len(xmltext))
	text = append(text, nonce()...)
	text = append(text, encodeNetworkByteOrder(uint32(len(xmltext)))...)
	text = append(text, xmltext...)
	text = append(text, c.appid...)
	text = common.PKCS7Padding(text, aes.BlockSize)

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, c.key[:aes.BlockSize])
	mode.CryptBlocks(text, text)

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(text)))
	base64.StdEncoding.Encode(dst, text)
	return dst, nil
}

// Decrypt 解密内容
func (c *Crypto) Decrypt(text []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(text)))
	n, err := base64.StdEncoding.Decode(dst, text)
	if err != nil {
		return nil, err
	}
	dst = dst[:n]

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, c.key[:aes.BlockSize])
	plaintext := make([]byte, len(dst))
	mode.CryptBlocks(plaintext, dst)

	return common.PKCS7UnPadding(plaintext), nil
}

// 编码成网络字节（大端）
func encodeNetworkByteOrder(n uint32) []byte {
	ret := make([]byte, 4)

	ret[0] = byte(n >> 24)
	ret[1] = byte(n >> 16)
	ret[2] = byte(n >> 8)
	ret[3] = byte(n)

	return ret
}

// 解码网格字节(大端)
func decodeNetworkByteOrder(b []byte) (n uint32) {
	return uint32(b[0])<<24 |
		uint32(b[1])<<16 |
		uint32(b[2])<<8 |
		uint32(b[3])
}
