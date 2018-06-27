package controller

import "time"

func UserReturnDto(ca, ua time.Time, username string, sex int) interface{} {
	var r struct {
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Username  string    `json:"username"`
		Sex       int       `json:"sex"`
	}
	r.CreatedAt = ca
	r.UpdatedAt = ua
	r.Username = username
	r.Sex = sex
	return r
}
