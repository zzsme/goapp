// Placeholder for array.go
package utils

// InArray 判断是否存在
func InArray[T comparable](item T, arr []T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

// RemoveDuplicate 去重
func RemoveDuplicate[T comparable](arr []T) []T {
	m := make(map[T]bool)
	var result []T
	for _, v := range arr {
		if !m[v] {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}
