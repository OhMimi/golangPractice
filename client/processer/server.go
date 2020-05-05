package processer

import (
	"chatroom/client/utils"
	"fmt"
	"net"
	"os"
)

// 顯示使用者菜單
func ShowMenu() {
	fmt.Println("---恭喜xxx登入成功---")
	fmt.Println("1.顯示在線用戶列表")
	fmt.Println("2.發送消息")
	fmt.Println("3.訊息列表")
	fmt.Println("4.退出系統")
	fmt.Println("請選擇(1-4):")

	var key int
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("1.")
	case 2:
		fmt.Println("2.")
	case 3:
		fmt.Println("3.")
	case 4:
		fmt.Println("你選擇退出系統...")
		os.Exit(0)
	default:
		fmt.Println("輸入有誤，請重新輸入...")
	}
}

// 保持與服務端的通訊
func serverProcessMes(conn net.Conn) {
	// 創建一個Transfer實例，讓他不停地讀取服務器發送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客戶端正在等待讀取服務器發送的消息...")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Printf("readPkg(conn) err = %v\n", err)
			return
		}
		fmt.Printf("mes = %v\n", mes)
	}

}
