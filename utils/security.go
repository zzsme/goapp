// Placeholder for security.go
package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
