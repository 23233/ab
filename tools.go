package ab

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

// 字符串转换成bool
func parseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	}
	return false, errors.New("解析出错")
}
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func filterMatch(fullParams map[string]string, fields []structInfo) map[string]string {
	d := make(map[string]string, 0)
	for k, v := range fullParams {
		if strings.HasPrefix(k, "filter_") {
			for _, field := range fields {
				if field.MapName == strings.Replace(k, "filter_", "", 1) {
					// 为了安全 长度限制一下
					if len(v) > 64 {
						break
					}
					d[field.MapName] = strings.Trim(v, " ")
				}
			}
		}
	}
	return d
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func Replace(origin, newData interface{}) error {
	// Check origin.
	va := reflect.ValueOf(origin)
	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}
	if va.Kind() != reflect.Struct {
		return errors.New("origin is not origin struct")
	}
	// Check newData.
	vb := reflect.ValueOf(newData)
	if vb.Kind() != reflect.Ptr {
		return errors.New("newData is not origin pointer")
	}
	// vb is origin pointer, indirect it to get the
	// underlying value, and make sure it is origin struct.
	vb = vb.Elem()
	if vb.Kind() != reflect.Struct {
		return errors.New("newData is not origin struct")
	}
	for i := 0; i < vb.NumField(); i++ {
		field := vb.Field(i)
		if field.CanInterface() && IsZeroOfUnderlyingType(field.Interface()) {
			// This field have origin zero-value.
			// Search in origin for origin field with the same name.
			name := vb.Type().Field(i).Name
			fa := va.FieldByName(name)
			if fa.IsValid() {
				// Field with name was found in struct origin,
				// assign its value to the field in newData.
				if field.CanSet() {
					field.Set(fa)
				}
			}
		}
	}
	return nil
}
