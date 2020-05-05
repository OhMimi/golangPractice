package model

import (
	"errors"
)

// 根據業務邏輯需要，自定義一些錯誤
var (
	ERROR_USER_NOTEXIST = errors.New("該用戶不存在")
	ERROR_USER_EXISTED  = errors.New("該用戶已存在")
	ERROR_USER_PWDERROR = errors.New("該用戶密碼不正確")
)
