package model

type User struct {
	UserID   int    `json:"userID"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
