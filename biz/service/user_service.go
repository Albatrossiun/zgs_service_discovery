package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/domain"
	"github.com/Albatrossiun/zgs_service_discovery/biz/model/zgs_service_discovery"
	"unsafe"
)

type UserService struct {
	userDomain *domain.UserDomain
}

func NewUserService() *UserService {
	return &UserService{
		userDomain: domain.NewUserDomain(),
	}
}

type AgentsObj struct {
	UUid        string `json:"uuid"`
	Ip          string `json:"ip"`
	Port        string `json:"port"`
	Status      string `json:"status"`
	Group       string `json:"group"`
	Ext         string `json:"ext"`
	OfflineTime int    `json:"offline_time"`
}

func (u *UserService) Regist(ctx context.Context, req zgs_service_discovery.RegistRequest) zgs_service_discovery.RegistResponse {
	agentsObj := &AgentsObj{
		UUid:   req.UUID,
		Ip:     req.IP,
		Port:   req.Port,
		Status: "online",
	}
	ipAndPortJson, err := json.Marshal(agentsObj)
	if err != nil {
		fmt.Println("service Regist Marshal err=", err)
		resp := zgs_service_discovery.RegistResponse{
			Code:    500,
			Message: err.Error(),
		}
		return resp
	}
	// TODO 幂等性校验
	err = u.userDomain.Regist("service_"+req.UUID, string(ipAndPortJson))
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
	agentsJsonList, err := u.userDomain.ListAgents()
	if err != nil {
		fmt.Println("service ListAgents() err=", err)
		return zgs_service_discovery.ListAgentsInfoResponse{
			Total: 0,
		}
	}
	if len(agentsJsonList) == 0 {
		fmt.Println("service ListAgents() agentsJsonList is empty", err)
		return zgs_service_discovery.ListAgentsInfoResponse{
			Total: 0,
		}
	}

	// set
	agentsStatus := req.Status
	type void struct{}
	var member void

	set := make(map[string]void) // New empty set
	for _, statusStr := range agentsStatus {
		set[statusStr] = member
	}

	var agentsObjList []AgentsObj
	for _, agentJson := range agentsJsonList {
		var agentsObj AgentsObj
		err = json.Unmarshal(*(*[]byte)(unsafe.Pointer(&agentJson)), &agentsObj)
		if err != nil {
			fmt.Println("service ListAgents() Unmarshal err=", err)
			return zgs_service_discovery.ListAgentsInfoResponse{}
		}

		if _, exists := set[agentsObj.Status]; !exists {
			agentsObjList = append(agentsObjList, agentsObj)
		}
	}

	var agents []*zgs_service_discovery.AgentInfo
	for _, obj := range agentsObjList {
		agent := AgentsObjToAgentInfo(obj)
		agents = append(agents, &agent)
	}

	total := len(agents)
	return zgs_service_discovery.ListAgentsInfoResponse{
		Total:  int32(total),
		Agents: agents,
	}
}
