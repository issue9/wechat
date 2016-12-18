// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pay

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

func TestEs(t *testing.T) {
	a := assert.New(t)

	v := []byte("<![CDATA[xxx]]>")
	buf := new(bytes.Buffer)
	a.NotError(xml.EscapeText(buf, v))
	t.Log(buf.String())
}
