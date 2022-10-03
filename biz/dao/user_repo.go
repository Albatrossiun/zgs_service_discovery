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

func (u *UserRepository) Regist(uuid, agentsObjJson string) error {
	// 通过go向redis写入数据
	_, err := pool.Get().Do("Set", uuid, agentsObjJson)
	if err != nil {
		fmt.Println("repo Regist err=", err)
		return err
	}
	return nil
}

func (u *UserRepository) ListAgents() ([]string, error) {
	// 读取数据
	uuid := "service_*"
	keys, err := pool.Get().Do("Keys", uuid)
	if err != nil {
		fmt.Println("repo ListAgents Get Do err=", err)
		return nil, err
	}

	jsonList, err := redis.Strings(pool.Get().Do("mget", keys.([]interface{})...))
	if err != nil {
		fmt.Println("repo ListAgents Get Do err=", err)
		return nil, err
	}
	//for _, value := range jsonList {
	//  fmt.Println(value)
	//}
	return jsonList, nil
}
