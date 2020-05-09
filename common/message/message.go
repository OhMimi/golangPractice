package message

const (
	LoginMesType       = "LoginMes"
	LoginResMesType    = "LoginResMes"
	RegisterMesType    = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
)

type Message struct {
	Type string `json:"type"` // 消息類型
	Data string `json:"data"` // 消息的訊息
}

type LoginMes struct {
	UserID   int    `json:"userID"`   // 用戶ID
	UserPwd  string `json:"userPwd"`  // 用戶密碼
	UserName string `json:"userName"` // 用戶名
}

type LoginResMes struct {
	Code  int    `json:"code"`  // 狀態碼 500 -> 該用戶未註冊  200 -> 登入成功
	Error string `json:"error"` // 錯誤訊息
}

type RegisterMes struct {
	User User // 類型就是User結構體
}

type RegisterResMes struct {
	Code  int    `json:"code"`  // 狀態碼 400 -> 該用戶已存在  200 -> 登入成功
	Error string `json:"error"` // 錯誤訊息
}
