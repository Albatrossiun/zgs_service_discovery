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
		UUID: "cdsx",
		IP:   "124.222.8.21",
		Port: "9911",
	}
	resp := srv.Regist(ctx, req)
	fmt.Println(resp)
}

func TestListAgents(t *testing.T) {
	var status []string
	status = append(status, "online")
	req := zgs_service_discovery.ListAgentsByGroupAndStatusRequest{
		Status: status,
	}
	resp := srv.ListAgentsByGroupAndStatus(ctx, req)
	fmt.Println(resp)
}

func TestUpdateOnlineAgentsGroup(t *testing.T) {
	var uuids []string
	uuids = append(uuids, "c7dc5d07-eead-4228-a711-40e5302e705d", "5345da4e-ecfa-4f1d-afeb-ac88bd441959", "ad9e32da-e2f2-4a54-b189-07f10e6780eb")

	req := zgs_service_discovery.UpdateOnlineAgentsGroupRequest{
		Group: "groupA",
		Uuids: uuids,
	}
	resp := srv.UpdateOnlineAgentsGroup(ctx, req)
	fmt.Println(resp)
}
