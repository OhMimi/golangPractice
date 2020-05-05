package main

import (
	"chatroom/common/message"
	"chatroom/server/processer"
	"chatroom/server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 編寫一個serverProcessMes函數
// 功能: 根據客戶端發送消息種類的不同，決定調用哪個函數來處理
func (p *Processor) serverProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		// 處理登入要求
		userProcesser := &processer.UserProcess{
			Conn: p.Conn,
		}
		err = userProcesser.ServerProcessLogin(mes)
	case message.RegisterMesType:
		// 處理註冊要求
	default:
		fmt.Println("消息類型不存在，無法處理...")
	}

	return
}

func (p *Processor) processControl() (err error) {
	// 讀取客戶端發送的消息
	for {
		// 這裡將讀取數據包，直接封裝成一個函數readPkg()，返回Message, Error
		// 這裡完成一個Transfer實例來完成讀封包的動作
		tf := &utils.Transfer{
			Conn: p.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客戶端退出，服務器端也退出...")
			} else {
				fmt.Printf("readPkg err = %v\n", err)
			}
			return err
		}

		// fmt.Printf("mes = %v\n", mes)
		err = p.serverProcessMes(&mes)
		if err != nil {
			fmt.Printf("serverProcessMes(conn, &mes) err = %v\n", err)
			return err
		}
	}
}
