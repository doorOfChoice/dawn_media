package tool

import "strconv"

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
