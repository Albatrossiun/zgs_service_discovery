package dao

import (
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/storage"
	"github.com/gomodule/redigo/redis"
)

var (
	conn           redis.Conn
	userRepository UserRepository
)

type UserRepository struct{}

func InitRedis() {
	conn = storage.GetRedisConn()
}

func GetUserRepository() UserRepository {
	return userRepository
}

func InitUserRepository() {
	userRepository = UserRepository{}
}

func (u *UserRepository) Regist(uuid, ipAndPort string) error {
	// 通过go向redis写入数据
	_, err := conn.Do("Set", uuid, ipAndPort)
	if err != nil {
		fmt.Println("repo Regist err=", err)
		return err
	}
	// 关闭连接
	defer conn.Close()
	return nil
}

func (u *UserRepository) ListAgents(uuid string) (string, error) {
	// 读取数据
	r, err := redis.String(conn.Do("Get", uuid))
	if err != nil {
		fmt.Println("repo ListAgents err=", err)
		return "", err
	}
	return r, nil
	//fmt.Println("Manipulate success, the ip and port of uuid is", r)
}
