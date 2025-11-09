// Placeholder for stringutil.go
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"unicode"
)

// MD5 计算
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// UpperFirst 首字母大写
func UpperFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// LowerFirst 首字母小写
func LowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// SplitTrim 切分并去除空格
func SplitTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var out []string
	for _, v := range parts {
		if str := strings.TrimSpace(v); str != "" {
			out = append(out, str)
		}
	}
	return out
}
