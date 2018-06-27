package controller

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"strconv"
)

func getInt(v string) int {
	num, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return num
}

func getInt64(v string) int64 {
	return int64(getInt(v))
}

func result(code int, data interface{}, msg string) gin.H {
	return gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	}
}

func md5Encode(s string) string {
	return md5EncodeWithSalt(s, "")
}

func md5EncodeWithSalt(s, salt string) string {
	m := md5.New()
	if salt != "" {
		m.Write([]byte(salt))
	}
	return hex.EncodeToString(m.Sum([]byte(s)))
}
