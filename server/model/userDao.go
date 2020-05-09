package model

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// 在服務器啟動後，就初始化一個UserDai=o實例
// 把它做成全局的變量，在需要和redis進行操作時，直接使用即可

var (
	MyUserDao *UserDao
)

// 定義一個UserDao結構體
// 完成對User結構體的各種操作
type UserDao struct {
	Pool *redis.Pool
}

// 使用工廠模式創建一個UserDao實例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		Pool: pool,
	}
	return
}

// 思考一下在UserDao應該提供哪些方法
// 1. 根據UserID返回一個User實例+err
func (ud *UserDao) getUserByID(conn redis.Conn, ID int) (user *User, err error) {
	// 通過給定ID到redis中查詢此用戶
	res, err := redis.String(conn.Do("HGET", "users", ID))
	if err != nil {
		// 出現此錯誤，代表在redis中沒有找到對應ID的用戶
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXIST
		} else {
			fmt.Printf("redis.String(conn.Do()) err = %v\n", err)
		}

		return
	}
	fmt.Println("res = ", res)
	// 這裡需要將res反序列化為User實例
	user = &User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Printf("json.Unmarshal() err = %v\n", err)
		return
	}
	return
}

// 完成登入驗證
// 1. Login完成登入驗證
// 2. 如果id跟pwd都正確，則返回一個user實例
// 2. 如果id跟pwd有錯誤，則返回一個錯誤訊息
func (ud *UserDao) Login(userID int, userPwd string) (user *User, err error) {
	// 先從UserDao取出一個連接
	conn := ud.Pool.Get()
	defer conn.Close()
	user, err = ud.getUserByID(conn, userID)
	if err != nil {
		// fmt.Printf("ud.getUserByID() err = %v\n", err)
		return
	}

	// 此時證明此用戶被獲取到，並驗證密碼
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWDERROR
	}
	return
}

func (ud *UserDao) Register(user *message.User) (err error) {
	// 先從UserDao取出一個連接
	conn := ud.Pool.Get()
	defer conn.Close()
	_, err = ud.getUserByID(conn, user.UserID)
	if err == nil {
		err = ERROR_USER_EXISTED
		return
	}

	// 此時證明此用戶尚未存在
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("json.Marshal err = %v\n", err)
		return
	}

	// 寫入redis
	_, err = conn.Do("HSET", "users", user.UserID, string(data))
	if err != nil {
		fmt.Printf("conn.Do err = %v\n", err)
		return
	}

	return
}
