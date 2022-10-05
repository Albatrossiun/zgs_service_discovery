package dao

import (
	"context"
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/storage"
	"testing"
)

var (
	ctx context.Context
)

func init() {
	err := storage.InitRedisPool()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	InitRedis()
	InitUserRepository()
}

func TestDeleteAgents(t *testing.T) {
	uuid := "service_vbnm"
	err := userRepository.DeleteAgents(uuid)
	if err != nil {
		fmt.Println("TestDeleteAgents() err = ", err)
	}
}

func TestGetAgentsByUUids(t *testing.T) {
	uuid := "d575052a-b2e9-4e84-a5c4-e4c2e9294843"
	var str []string
	str = append(str, uuid)
	jsonlist, err := userRepository.GetAgentsByUUids(str)
	if err != nil {
		fmt.Println("TestDeleteAgents() err = ", err)
	}
	fmt.Println(jsonlist)
}
