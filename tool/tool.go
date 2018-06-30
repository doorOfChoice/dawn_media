package tool

import (
	"strconv"
	"crypto/md5"
	"encoding/hex"
)

func GetInt(v string) int {
	num, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return num
}

func GetInts(vs []string) []int {
	rs := make([]int, 0)
	for _, v := range vs {
		rs = append(rs, GetInt(v))
	}
	return rs
}

func GetInt64(v string) int64 {
	return int64(GetInt(v))
}

/**
普通MD5
 */
func Md5Encode(s string) string	 {
	return Md5EncodeWithSalt(s, "")
}

/**
加盐MD5操作
 */
func Md5EncodeWithSalt(s, salt string) string {
	m := md5.New()
	if salt != "" {
		m.Write([]byte(salt))
	}
	return hex.EncodeToString(m.Sum([]byte(s)))
}