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
