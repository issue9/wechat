// SPDX-FileCopyrightText: 2016-2024 caixw
//
// SPDX-License-Identifier: MIT

package internal

import (
	"encoding/xml"
	"errors"
	"io"
	"reflect"
	"strconv"
)

// MapFromXMLReader 从 io.Reader 读取内容，并填充到 map 中
func MapFromXMLReader(r io.Reader) (map[string]string, error) {
	ret := make(map[string]string, 10)
	d := xml.NewDecoder(r)
	for token, err := d.Token(); true; token, err = d.Token() {
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		elem, ok := token.(xml.StartElement)
		if !ok {
			continue
		}

		name := elem.Name.Local
		if name == "xml" {
			continue
		}

		token, err = d.Token()
		if err != nil { // 此处若是 io.EOF，也是属于非正常结束
			return nil, err
		}
		bs, ok := token.(xml.CharData)
		if !ok {
			return nil, errors.New("无法转换成 xml.CharData")
		}
		ret[name] = string(bs)
	}

	return ret, nil
}

// Map2XMLObj 将 map 转换到 v
func Map2XMLObj(maps map[string]string, v interface{}) error {
	values := values(v)

	for k, v := range maps {
		val, found := values[k]
		if !found {
			continue
		}

		switch val.Kind() {
		case reflect.String:
			val.SetString(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			val.SetInt(x)
		} // end switch
	} // end for

	return nil
}

// 将 obj 各个字段以 xml 标签中的值进行索引，方便查找。
func values(obj interface{}) map[string]reflect.Value {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()
	values := make(map[string]reflect.Value, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("xml")
		values[tag] = v.Field(i)
	}

	return values
}
