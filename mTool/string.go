package mTool

import (
	"crypto/md5"
	"encoding/hex"
)

func IsStringEmpty(str string) bool {
	return str == ""
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}