package processer

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

// 編寫一個serverProcessLogin函數，專門處理登入請求
func (up *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 先從mes中取出mes.Data，並直接反序列化為LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Printf("json.Unmarshal(mes.Data) err = %v\n", err)
		return
	}

	// 1. 先聲明一個回傳訊息 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 2. 再聲明一個回傳Data LoginResMes
	var loginResMes message.LoginResMes
	// 需要到redis中完成驗證
	// 使用model.UserDao到redis驗證
	user, err := model.MyUserDao.Login(loginMes.UserID, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXIST {
			loginResMes.Code = 500 // 500 表示用戶未註冊
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWDERROR {
			loginResMes.Code = 403 // 403 表示用戶密碼錯誤
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505 // 505 表示未知訊息
			loginResMes.Error = "服務器內部錯誤..."
		}
		// loginResMes.Code = 500 // 500 表示用戶未註冊
		// loginResMes.Error = "該用戶不存在，請註冊後再使用..."

	} else {
		loginResMes.Code = 200 // 200 表示登入成功
		fmt.Println("登入成功user = ", user)
	}

	// 3. 將loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Printf("json.Marshal(loginResMes) err = %v\n", err)
		return
	}

	// 4. 將data賦值給resMes.Data
	resMes.Data = string(data)

	// 5. 將resMes序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Printf("json.Marshal(resMes) err = %v\n", err)
		return
	}

	// 6. 發送Data，我們將其封裝到writePkg函數中
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	// if err != nil {
	// 	fmt.Printf("writePkg err = %v\n", err)
	// 	return
	// }

	return
}
