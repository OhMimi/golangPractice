package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	var buf = make([]byte, 8096)
	fmt.Println("讀取客戶端發送的數據...")
	_, err = conn.Read(buf[0:4])
	if err != nil {
		fmt.Printf("conn.Read(b) err = %v\n", err)
		return
	}
	// fmt.Printf("讀到的buf = %v\n", buf[0:4])

	// 根據buf[0:4]轉換成一個uint32的類型
	var pkgLen = binary.BigEndian.Uint32(buf[0:4])

	// 根據pkgLen來讀取內容
	n, err := conn.Read(buf[0:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Printf("conn.Read(buf[0:pkgLen]) err = %v\n", err)
		return
	}

	// 將pkgLen進行反序列化成 -> message.Message
	err = json.Unmarshal(buf[0:pkgLen], &mes)
	if err != nil {
		fmt.Printf("json.Unmarshal(buf[0:pkgLen]) err = %v\n", err)
		return
	}
	// fmt.Println("mes=", mes)
	return
}

func writePkg(conn net.Conn, data []byte) (err error) {
	// 先發一個長度給對方
	var pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 發送長度
	n, err := conn.Write(buf[0:4])
	if n != 4 || err != nil {
		fmt.Printf("conn.Write(buf[0:4]) err = %v\n", err)
		return
	}

	// 發送數據本身
	n, err = conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Printf("conn.Write(data)) err = %v\n", err)
		return
	}
	return
}

// 編寫一個serverProcessLogin函數，專門處理登入請求
func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
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

	// 如果用戶id = 100 , 密碼 = 1234 認為合法，否則不合法
	if loginMes.UserID == 100 && loginMes.UserPwd == "1234" {
		//合法
		loginResMes.Code = 200 // 200 表示登入成功
	} else {
		//不合法
		loginResMes.Code = 500 // 500 表示用戶未註冊
		loginResMes.Error = "該用戶不存在，請註冊後再使用..."
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
	err = writePkg(conn, data)
	// if err != nil {
	// 	fmt.Printf("writePkg err = %v\n", err)
	// 	return
	// }

	return
}

// 編寫一個serverProcessMes函數
// 功能: 根據客戶端發送消息種類的不同，決定調用哪個函數來處理
func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		// 處理登入要求
		err = serverProcessLogin(conn, mes)
	case message.RegisterMesType:
		// 處理註冊要求
	default:
		fmt.Println("消息類型不存在，無法處理...")
	}

	return
}

//處理和客戶端的通訊
func process(conn net.Conn) {
	// 延時關閉
	defer conn.Close()
	// 讀取客戶端發送的消息
	for {
		// 這裡將讀取數據包，直接封裝成一個函數readPkg()，返回Message, Error
		mes, err := readPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客戶端退出，服務器端也退出...")
			} else {
				fmt.Printf("readPkg err = %v\n", err)
			}
			return
		}

		// fmt.Printf("mes = %v\n", mes)
		err = serverProcessMes(conn, &mes)
		if err != nil {
			fmt.Printf("serverProcessMes(conn, &mes) err = %v\n", err)
			return
		}
	}
}

func main() {
	//提示訊息
	fmt.Println("服務器在8889端口監聽...")
	listener, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Printf("net.Listen err = %v\n", err)
		return
	}

	defer listener.Close()

	// 一旦監聽成功，就等待客戶端來連接服務器
	for {
		fmt.Println("等待客戶端來連接服務器...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("listener.Accept err = %v\n", err)
		}

		// 一旦連接成功，就啟動一個協程跟客戶端保持通訊
		go process(conn)
	}
}
