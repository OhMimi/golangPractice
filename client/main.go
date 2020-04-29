package main

import (
	"fmt"
)

// 定義兩個全局變量
var userID int
var userPwd string

func main() {
	fmt.Println("main")
	// 接收用戶的選擇
	var key int
	// 判斷是否還繼續顯示菜單
	var loop = true

	for loop {
		fmt.Println("---歡迎登入多人聊天系統---")
		fmt.Println("1.登入聊天室")
		fmt.Println("2.註冊用戶")
		fmt.Println("3.退出系統")
		fmt.Println("請選擇(1-3):")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登入聊天室")
			loop = false
		case 2:
			fmt.Println("註冊用戶")
			loop = false
		case 3:
			fmt.Println("退出系統")
			loop = false
		default:
			fmt.Println("輸入有誤，請重新輸入...")
		}
	}

	//確認用戶的輸入，顯示新的提示信息
	if key == 1 {
		// 說明用戶要登入
		fmt.Println("請輸入用戶的ID:")
		fmt.Scanf("%d\n", &userID)
		fmt.Println("請輸入用戶的密碼:")
		fmt.Scanf("%s\n", &userPwd)
		// 先把登入函數寫到另一文件
		err := login(userID, userPwd)
		if err != nil {
			fmt.Printf("login err =%v\n", err)
			return
		} else {
			fmt.Println("login no err")
		}
	} else if key == 2 {
		fmt.Println("用戶註冊")
	}
}
