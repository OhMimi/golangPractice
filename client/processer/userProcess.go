package processer

import (
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
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

	// 7. 到這個時候 data 就是我們要發送的訊息
	// 7.1 先把 data長度發送給服務器
	// 先獲取到 data的長度 -> 轉成一個表示長度的byte切片
	var pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	// 發送長度
	n, err := conn.Write(buf[0:4])
	if n != 4 || err != nil {
		fmt.Printf("conn.Write(bytes) err = %v\n", err)
		return
	}

	fmt.Printf("客戶端，發送消息的長度 = %d 內容 = %s\n", len(data), string(data))

	// 發送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("conn.Write(data) err = %v\n", err)
		return
	}

	// 創建一個Transfer實例
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 這裡需要處理服務器端回傳的消息
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
