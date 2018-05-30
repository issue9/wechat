// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package crypto 消息的加解密功能
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/issue9/wechat/common"
)

const messageFormat = `<xml>
<Encrypt><![CDATA[%s]]></Encrypt>
<MsgSignature><![CDATA[%s]]></MsgSignature>
<TimeStamp>%s</TimeStamp>
<Nonce><![CDATA[%s]]></Nonce>
</xml>`

type receiver struct {
	Root       xml.Name `xml:"xml"`
	ToUserName string   `xml:"ToUserName"`
	Encrypt    string   `xml:"Encrypt"`
}

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

// Encrypt 加密 XML 内容
//
// 返回加密后的 XML 结构体内容以及签名内容
func (c *Crypto) Encrypt(xmltext []byte, timestamp, nonce string) ([]byte, string, error) {
	entext, err := c.encrypt(xmltext)
	if err != nil {
		return nil, "", err
	}

	if timestamp == "" {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	sign := sha1Sign(c.token, timestamp, nonce)
	return []byte(fmt.Sprintf(messageFormat, entext, sign, timestamp, nonce)), sign, nil
}

// Decrypt 解密 XML 内容
func (c *Crypto) Decrypt(body []byte, sign, timestamp, nonce string) ([]byte, error) {
	r := &receiver{}

	if timestamp == "" {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	if sha1Sign(c.token, timestamp, nonce) != sign {
		return nil, errors.New("签名不同")
	}

	if err := xml.Unmarshal(body, r); err != nil {
		return nil, err
	}

	return c.decrypt([]byte(r.Encrypt))
}

// base64Encoding(AES_Encrypt[random(16B) + msg_len(4B) + rawXMLMsg + appId])
func (c *Crypto) encrypt(xmltext []byte) ([]byte, error) {
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

func (c *Crypto) decrypt(text []byte) ([]byte, error) {
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

	plaintext = common.PKCS7UnPadding(plaintext)

	size := decodeNetworkByteOrder(plaintext[16:20])
	return plaintext[20 : 20+int(size)], nil
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

func sha1Sign(token, timestamp, nonce string) (signature string) {
	strs := sort.StringSlice{token, timestamp, nonce}
	strs.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce))
	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}
