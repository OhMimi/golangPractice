package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	var buf = make([]byte, 8096)
	fmt.Println("讀取服務端發送的數據...")
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
