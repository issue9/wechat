// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"encoding/xml"
	"io"
	"reflect"
	"strconv"
)

// MapFromReader 从 io.Reader 读取内容，并填充到 map 中
func MapFromReader(r io.Reader) (map[string]string, error) {
	ret := make(map[string]string, 10)
	d := xml.NewDecoder(r)
	for token, err := d.Token(); true; token, err = d.Token() {
		if err != nil {
			return nil, err
		}

		var key, val string
		switch t := token.(type) {
		case xml.StartElement:
			key = t.Name.Local
		case xml.CharData:
			val = string(t)
		}
		ret[key] = val
	}

	return ret, nil
}

// Map2XMLObj 将 map 转换到 v
func Map2XMLObj(maps map[string]string, v interface{}) error {
	values, err := values(v)
	if err != nil {
		return err
	}

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

// 将 Return 各个字段以 xml 标签中的值进行索引，方便查找。
func values(obj interface{}) (map[string]reflect.Value, error) {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()
	values := make(map[string]reflect.Value, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("xml")
		values[tag] = v.Field(i)
	}

	return values, nil
}
