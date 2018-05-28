// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package common

import "bytes"

// PKCS7UnPadding 解码
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return plantText[:(length - unPadding)]
}

// PKCS7Padding 编码
func PKCS7Padding(plantText []byte, blockSize int) []byte {
	padding := blockSize - len(plantText)%blockSize
	repeat := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plantText, repeat...)
}
