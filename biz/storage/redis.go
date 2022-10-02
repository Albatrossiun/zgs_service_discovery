package storage

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	redisPool *redis.Pool
)

func InitRedisPool() error {
	redisPool = &redis.Pool{
		MaxIdle:     1,
		MaxActive:   20,
		IdleTimeout: 6 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "124.222.8.21:6379", redis.DialPassword("123abc"))
		},
	}
	conn := redisPool.Get()
	res, err := conn.Do("ping")
	if err != nil {
		fmt.Println(err)
		return err
	}
	if res.(string) != "PONG" {
		fmt.Println("res = ", res)
		return fmt.Errorf("ping failed")
	}
	return nil
}

func GetRedisPool() *redis.Pool {
	return redisPool
}
