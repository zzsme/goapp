// Placeholder for convert.go
package utils

import (
	"fmt"
	"strconv"
)

// StringToInt 安全转换
func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func StringToUint64(s string) uint64 {
	u, _ := strconv.ParseUint(s, 10, 64)
	return u
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func FloatToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
