// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package internal

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert/v4"
)

type CDataTester struct {
	XMLName xml.Name `xml:"xml"`
	Data1   CData    `xml:"data1"`
}

func TestMarshalCData(t *testing.T) {
	a := assert.New(t, false)

	obj := &CDataTester{
		Data1: CData{"abc<def"},
	}

	bs, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(bs), `<xml><data1><![CDATA[abc<def]]></data1></xml>`)

}

func TestUnmarshalCData(t *testing.T) {
	a := assert.New(t, false)

	data := []byte(`
	<xml>
	<data1><![CDATA[adf<ddd]]></data1>
	</xml>
	`)
	o := &CDataTester{}
	a.NotError(xml.Unmarshal(data, o)).Equal(o.Data1.Text, "adf<ddd")
}
