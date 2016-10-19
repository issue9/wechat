// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"testing"

	"github.com/issue9/assert"
)

var _ Messager = &Text{}
var _ Messager = &Image{}
var _ Messager = &Voice{}
var _ Messager = &Video{}
var _ Messager = &ShortVideo{}
var _ Messager = &Location{}
var _ Messager = &Link{}

func TestGetMsgType(t *testing.T) {
	a := assert.New(t)

	data := []byte(`<xml>
	<MsgType>image</MsgType>
	</xml>`)

	typ, err := getMsgType(data)
	a.NotError(err).Equal(typ, "image")
}
