package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/domain"
	"github.com/Albatrossiun/zgs_service_discovery/biz/model/zgs_service_discovery"
)

type UserService struct {
	userDomain *domain.UserDomain
}

func NewUserService() *UserService {
	return &UserService{
		userDomain: domain.NewUserDomain(),
	}
}

type IpAndPortObj struct {
	Ip   string `json:"ip"`
	Port string `json:"Port"`
}

func (u *UserService) Regist(ctx context.Context, req zgs_service_discovery.RegistRequest) zgs_service_discovery.RegistResponse {
	ipAndPortObj := &IpAndPortObj{
		Ip:   req.IP,
		Port: req.Port,
	}
	ipAndPortJson, err := json.Marshal(ipAndPortObj)
	if err != nil {
		fmt.Println("service Regist Marshal err=", err)
		resp := zgs_service_discovery.RegistResponse{
			Code:    500,
			Message: err.Error(),
		}
		return resp
	}
	err = u.userDomain.Regist(req.UUID, string(ipAndPortJson))
	if err != nil {
		fmt.Println("service Regist err=", err)
		resp := zgs_service_discovery.RegistResponse{
			Code:    500,
			Message: err.Error(),
		}
		return resp
	}
	return zgs_service_discovery.RegistResponse{Code: 200}
}

func (u *UserService) ListAgents(ctx context.Context, req zgs_service_discovery.ListAgentsInfoRequest) zgs_service_discovery.ListAgentsInfoResponse {
	return zgs_service_discovery.ListAgentsInfoResponse{}
}
