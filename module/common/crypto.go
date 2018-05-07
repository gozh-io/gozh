package common

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

//sha256加密
func Sha256Encode(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 获取字符串的MD5
func MD5Encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
