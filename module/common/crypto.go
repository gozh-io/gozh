package common

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

//sha256加密
func Sha256Encode(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return string(bs)
}

// 获取字符串的MD5
func MD5Encode(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}
