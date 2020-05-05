package main

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// 定義一個全局的pool
var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     // 最大空閒連接數
		MaxActive:   maxActive,   // 表示和數據庫的最大連接數，0表示無限制
		IdleTimeout: idleTimeout, // 最大空閒時間
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
