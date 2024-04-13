package stringutils

import "strings"

// ContainsString 判断字符串数组是否包含指定的字符串
func ContainsString(arr []string, target string) bool {
	for _, s := range arr {
		if strings.Contains(s, target) {
			return true
		}
	}
	return false
}
