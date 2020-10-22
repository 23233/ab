package ab

import (
	"github.com/pkg/errors"
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
					if len(v) >= 25 {
						break
					}
					d[field.MapName] = strings.Trim(v, " ")
				}
			}
		}
	}
	return d
}
