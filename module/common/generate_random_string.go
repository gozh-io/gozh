package common

import (
	"math/rand"
	"strings"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// 生成len长度的随机字符串
func GenerateRandomString(len int) string {
	rand.Seed(time.Now().UnixNano())
	randomString := make([]byte, len)
	for i := range randomString {
		randomString[i] = letterBytes[rand.Intn(strings.Count(letterBytes, "")-1)]
	}
	return string(randomString)
}
