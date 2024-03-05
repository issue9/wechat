// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package common

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/issue9/assert/v4"
)

var _ error = &Result{}

func TestNewResult(t *testing.T) {
	a := assert.New(t, false)

	// HTTP 错误代码
	r := NewResult(http.StatusBadRequest)
	a.Equal(r.Message, http.StatusText(http.StatusBadRequest))

	// 微信错误代码
	r = NewResult(40002)
	a.Equal(r.Message, messages[40002])

	// 不存在
	r = NewResult(-100)
	a.Equal(r.Code, 601).Equal(r.Message, messages[601]+strconv.Itoa(-100))
}

func TestFrom(t *testing.T) {
	a := assert.New(t, false)

	r := From([]byte("d"))
	a.Equal(r.Code, 600)

	r = From([]byte(`{
		"errcode":40002,
		"errmsg": "error message"
	}`))
	a.Equal(r.Code, 40002).Equal(r.Message, "error message")
}

func TestResult_IsOK(t *testing.T) {
	a := assert.New(t, false)

	r := NewResult(0)
	a.True(r.IsOK())

	r = NewResult(http.StatusOK)
	a.True(r.IsOK())

	r = NewResult(-1)
	a.False(r.IsOK())

	r = NewResult(http.StatusInternalServerError)
	a.False(r.IsOK())
}
