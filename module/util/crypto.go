package util

import (
	"crypto/sha256"
	"crypto/md5"
	"encoding/hex"
	"github.com/globalsign/mgo/bson"
)

// 密码的sha256加密
func cryptoSha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return string(bs)
}

// 获取字符串的MD5Hash
func GetMD5Hash(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// 获取user的id
func GetUserId(username string) bson.ObjectId {
	return bson.ObjectIdHex(GetMD5Hash(username)[4:28])
}