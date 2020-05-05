package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)

//處理和客戶端的通訊
func process(conn net.Conn) {
	// 延時關閉
	defer conn.Close()
	// 這裡需要創建一個processor總控實例
	pro := &Processor{
		Conn: conn,
	}
	err := pro.processControl()
	if err != nil {
		fmt.Printf("pro.processControl() err = %v\n", err)
		return
	}
}

// 這裡編寫一個函數完成對UserDao的初始化
func initUserDao() {
	// 這裡需要注意初始化順序，因為UserDao需要依賴pool
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	// 當服務器啟動時，就初始化redis連接池
	initPool("0.0.0.0:6379", 16, 0, 300*time.Second)
	initUserDao()
	//提示訊息
	fmt.Println("new set服務器在8889端口監聽...")
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
