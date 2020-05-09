package processer

import (
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
}

func (up *UserProcess) Login(userID int, userPwd string) (err error) {
	// fmt.Printf("userID = %d userPwd = %s\n", userID, userPwd)

	// 1. 連接到服務器
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Printf("net.Dial err = %v\n", err)
		return
	}

	// 延時關閉
	defer conn.Close()

	// 2. 準備通過conn發送消息給服務器
	var mes message.Message
	mes.Type = message.LoginMesType

	// 3. 創建一個LoginMes結構體
	var loginMes message.LoginMes
	loginMes.UserID = userID
	loginMes.UserPwd = userPwd

	// 4. 將loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Printf("json.Marshal err = %v\n", err)
		return
	}

	// 5. 將data賦給mes.Data字段
	mes.Data = string(data)

	// 6. 將mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Printf("json.Marshal err = %v\n", err)
		return
	}

	// 7. 發送data長度發送給服務器

	// 創建一個Transfer實例
	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("tf.WritePkg() err = %v\n", err)
		return
	}

	// 8. 這裡需要處理服務器端回傳的消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Printf("readPkg(conn) err = %v\n", err)
		return
	}

	// 將mes的Data的部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Printf("json.Unmarshal([]byte(mes.Data)) err = %v\n", err)
		return
	}

	if loginResMes.Code == 200 {
		// fmt.Println("登入成功")
		// 這裡需要起一個客戶端的協程
		// 該協程保持與服務端的通訊，如果服務器有數據推送給客戶端
		// 則接收並顯示在客戶端的終端
		go serverProcessMes(conn)
		// 1. 顯示成功登入後的菜單
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}

	return
}

func (up *UserProcess) Register(userID int, userPwd, userName string) (err error) {
	// 1. 連接到服務器
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Printf("net.Dial err = %v\n", err)
		return
	}
	defer conn.Close()
	// 2. 準備通過conn發送消息給服務端
	var mes message.Message
	mes.Type = message.RegisterMesType
	// 3. 創建一個RegisterMes結構體
	var registerMes message.RegisterMes
	registerMes.User.UserID = userID
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4. 將loginMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Printf("json.Marshal err = %v\n", err)
		return
	}

	// 5. 將data賦給mes.Data字段
	mes.Data = string(data)

	// 6. 將mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Printf("json.Marshal err = %v\n", err)
		return
	}

	// 創建一個Transfer實例
	tf := &utils.Transfer{
		Conn: conn,
	}
	// 7. 發送data長度發送給服務器
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("tf.WritePkg() err = %v\n", err)
		return
	}

	// 8. 這裡需要處理服務器端回傳的消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Printf("readPkg() err = %v\n", err)
		return
	}
	// 將mes的Data的部分反序列化成RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if err != nil {
		fmt.Printf("json.Unmarshal([]byte(mes.Data)) err = %v\n", err)
		return
	}

	if registerResMes.Code == 200 {
		fmt.Println("註冊成功，請重新登入")
	} else {
		fmt.Println(registerResMes.Error)
	}
	os.Exit(0)
	return
}
