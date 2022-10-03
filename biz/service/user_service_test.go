package service

import (
	"context"
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/dao"
	"github.com/Albatrossiun/zgs_service_discovery/biz/model/zgs_service_discovery"
	"github.com/Albatrossiun/zgs_service_discovery/biz/storage"
	"testing"
)

var (
	srv *UserService
	ctx context.Context
)

func init() {
	err := storage.InitRedisPool()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	dao.InitRedis()
	dao.InitUserRepository()
	srv = NewUserService()
}

func TestRegist(t *testing.T) {
	var group []string
	group = append(group, "online")
	req := zgs_service_discovery.RegistRequest{
		UUID: "qwer",
		IP:   "124.222.8.21",
		Port: "9911",
	}
	resp := srv.Regist(ctx, req)
	fmt.Println(resp)
}

func TestListAgents(t *testing.T) {
	var group []string
	group = append(group, "online")
	req := zgs_service_discovery.ListAgentsInfoRequest{
		Group: group,
	}
	resp := srv.ListAgents(ctx, req)
	fmt.Println(resp)
}
