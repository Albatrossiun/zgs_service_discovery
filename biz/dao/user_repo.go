package dao

import (
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/storage"
	"github.com/gomodule/redigo/redis"
)

var (
	pool           *redis.Pool
	userRepository UserRepository
)

type UserRepository struct{}

func InitRedis() {
	pool = storage.GetRedisPool()
}

func GetUserRepository() UserRepository {
	return userRepository
}

func InitUserRepository() {
	userRepository = UserRepository{}
}

func (u *UserRepository) Regist(uuid, ipAndPort string) error {
	// 通过go向redis写入数据
	_, err := pool.Get().Do("Set", uuid, ipAndPort)
	if err != nil {
		fmt.Println("repo Regist err=", err)
		return err
	}
	return nil
}

func (u *UserRepository) ListAgents(uuid string) (string, error) {
	// 读取数据
	r, err := redis.String(pool.Get().Do("Get", uuid))
	if err != nil {
		fmt.Println("repo ListAgents err=", err)
		return "", err
	}
	return r, nil
}
