package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

//启动程序时，就初始化链接池
var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTime time.Duration ) {
	pool = &redis.Pool{
		MaxIdle: maxIdle,
		MaxActive: maxActive,
		IdleTimeout: idleTime,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
