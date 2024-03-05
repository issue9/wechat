// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package internal

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
