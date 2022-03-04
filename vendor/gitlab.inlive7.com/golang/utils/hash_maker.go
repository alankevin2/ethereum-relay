package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// @title CreateSHAHash
// @description 生成SHA Hash
// @param key 目標字串
// @param shaType SHA類別
// @return string SHA Hash字串
func CreateSHAHash(key string, shaType string) string {
	switch shaType {
	case "sha256":
		hasher := sha256.New()
		hasher.Write([]byte(key))
		return hex.EncodeToString(hasher.Sum(nil))
	case "sha512":
		hasher := sha512.New()
		hasher.Write([]byte(key))
		return hex.EncodeToString(hasher.Sum(nil))
	default:
		return ""
	}
}
