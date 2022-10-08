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

// Regist create + update
func (u *UserRepository) Regist(uuid, agentsObjJson string) error {
	// 通过go向redis写入数据
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("Set", uuid, agentsObjJson)
	if err != nil {
		fmt.Println("repo Regist err = ", err)
		return err
	}
	return nil
}

// MRegist mCreate + mUpdate
func (u *UserRepository) MRegist(uuid, agentsObjJson string) error {
	// 通过go向redis写入数据
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("Set", uuid, agentsObjJson)
	if err != nil {
		fmt.Println("repo Regist err = ", err)
		return err
	}
	return nil
}

func (u *UserRepository) ListAgents() ([]string, error) {
	// 读取数据
	conn := pool.Get()
	defer conn.Close()
	uuid := "service_*"
	keys, err := conn.Do("Keys", uuid)
	if err != nil {
		fmt.Println("repo ListAgents Get Do err = ", err)
		return nil, err
	}

	if keys == nil || len(keys.([]interface{})) == 0 {
		fmt.Println("repo ListAgents is empty")
		return nil, nil
	}
	jsonList, err := redis.Strings(conn.Do("mget", keys.([]interface{})...))
	if err != nil {
		fmt.Println("repo ListAgents Get Do err = ", err)
		return nil, err
	}
	//for _, value := range jsonList {
	//  fmt.Println(value)
	//}
	return jsonList, nil
}

func (u *UserRepository) DeleteAgents(uuid string) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("Del", uuid)
	if err != nil {
		fmt.Println("repo Regist delete agent err = ", err)
		return err
	}
	return nil
}

func (u *UserRepository) GetAgentsByUUids(uuids []string) ([]string, error) {
	// 读取数据
	conn := pool.Get()
	defer conn.Close()
	var keys []interface{}
	for _, uid := range uuids {
		uidstr := "service_" + uid
		keys = append(keys, uidstr)
	}
	jsonList, err := redis.Strings(conn.Do("mget", keys...))
	if err != nil {
		fmt.Println("repo ListAgents MGet Do err = ", err)
		return nil, err
	}
	return jsonList, nil
}
