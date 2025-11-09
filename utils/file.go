// Placeholder for file.go
package utils

import (
	"encoding/json"
	"os"
)

// WriteFile 写入文件
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// ReadFile 读取文件
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteJSONFile 写入 JSON 文件
func WriteJSONFile(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return WriteFile(path, data)
}

// ReadJSONFile 读取 JSON 文件
func ReadJSONFile(path string, v interface{}) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
