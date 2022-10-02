package storage

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var (
	redisConn redis.Conn
)

func InitRedisConn() {
	var err error
	redisConn, err = redis.Dial("tcp", "124.222.8.21:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(redisConn)
}

func GetRedisConn() redis.Conn {
	return redisConn
}
