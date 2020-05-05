package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//  這裡將這些方法關連到結構體
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte // 傳輸時用的緩衝
}

func (t *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("讀取客戶端發送的數據...")
	_, err = t.Conn.Read(t.Buf[0:4])
	if err != nil {
		fmt.Printf("conn.Read(b) err = %v\n", err)
		return
	}
	// fmt.Printf("讀到的t.Buf = %v\n", t.Buf[0:4])

	// 根據t.Buf[0:4]轉換成一個uint32的類型
	var pkgLen = binary.BigEndian.Uint32(t.Buf[0:4])

	// 根據pkgLen來讀取內容
	n, err := t.Conn.Read(t.Buf[0:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Printf("conn.Read(t.Buf[0:pkgLen]) err = %v\n", err)
		return
	}

	// 將pkgLen進行反序列化成 -> message.Message
	err = json.Unmarshal(t.Buf[0:pkgLen], &mes)
	if err != nil {
		fmt.Printf("json.Unmarshal(t.Buf[0:pkgLen]) err = %v\n", err)
		return
	}
	// fmt.Println("mes=", mes)
	return
}

func (t *Transfer) WritePkg(data []byte) (err error) {
	// 先發一個長度給對方
	var pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(t.Buf[0:4], pkgLen)
	// 發送長度
	n, err := t.Conn.Write(t.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Printf("conn.Write(t.Buf[0:4]) err = %v\n", err)
		return
	}

	// 發送數據本身
	n, err = t.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Printf("conn.Write(data)) err = %v\n", err)
		return
	}
	return
}
