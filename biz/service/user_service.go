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

type AgentsObj struct {
	UUid        string `json:"uuid"`
	Ip          string `json:"ip"`
	Port        string `json:"port"`
	Status      string `json:"status"`
	Group       string `json:"group"`
	Ext         string `json:"ext"`
	OfflineTime int64  `json:"offline_time"`
}

func (u *UserService) Regist(ctx context.Context, req zgs_service_discovery.RegistRequest) zgs_service_discovery.RegistResponse {
	agentsObj := &AgentsObj{
		UUid:        req.UUID,
		Ip:          req.IP,
		Port:        req.Port,
		Status:      "online",
		OfflineTime: 0,
	}
	agentsObjJson, err := json.Marshal(agentsObj)
	if err != nil {
		fmt.Println("service Regist Marshal err=", err)
		return zgs_service_discovery.RegistResponse{
			Code:    500,
			Message: err.Error(),
		}
	}

	// 获取redis所有的agents
	agentsObjList, err := u.ListAgentsWithUnmarshal()
	if err != nil {
		fmt.Println("service u.GetAllAgentsInRedis() err=", err)
		return zgs_service_discovery.RegistResponse{
			Code:    500,
			Message: err.Error(),
		}
	}
	for _, obj := range agentsObjList {
		if obj.UUid == req.UUID {
			fmt.Println("service Regist agent has existed")
			return zgs_service_discovery.RegistResponse{
				Code:    200,
				Message: "multiple registration",
			}
		}
	}

	err = u.userDomain.Regist("service_"+req.UUID, string(agentsObjJson))
	if err != nil {
		fmt.Println("service Regist err=", err)
		return zgs_service_discovery.RegistResponse{
			Code:    500,
			Message: err.Error(),
		}
	}
	return zgs_service_discovery.RegistResponse{Code: 200, Message: "registered successfully"}
}

func (u *UserService) ListAgentsByGroupAndStatus(ctx context.Context, req zgs_service_discovery.ListAgentsByGroupAndStatusRequest) zgs_service_discovery.ListAgentsByGroupAndStatusResponse {
	agentsObjList, err := u.ListAgentsWithUnmarshal()
	if err != nil {
		fmt.Println("service ListAgents() err=", err)
		return zgs_service_discovery.ListAgentsByGroupAndStatusResponse{
			Total: 0,
		}
	}
	if len(agentsObjList) == 0 {
		fmt.Println("service ListAgents() agentsJsonList is empty")
		return zgs_service_discovery.ListAgentsByGroupAndStatusResponse{
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

	var agentsList []AgentsObj
	for _, obj := range agentsObjList {
		if _, exists := set[obj.Status]; !exists {
			agentsList = append(agentsList, obj)
		}
	}

	var agents []*zgs_service_discovery.AgentInfo
	for _, obj := range agentsList {
		agent := AgentsObjToAgentInfo(obj)
		agents = append(agents, &agent)
	}

	total := len(agents)
	return zgs_service_discovery.ListAgentsByGroupAndStatusResponse{
		Total:  int32(total),
		Agents: agents,
	}
}

func (u *UserService) ListAgentsWithUnmarshal() ([]AgentsObj, error) {
	agentsJsonList, err := u.userDomain.ListAgents()
	if err != nil {
		fmt.Println("service GetAllAgentsInRedis() err=", err)
		return nil, err
	}
	if len(agentsJsonList) == 0 {
		fmt.Println("service GetAllAgentsInRedis() agentsJsonList is empty")
		return nil, err
	}

	var agentsObjList []AgentsObj
	for _, agentJson := range agentsJsonList {
		var agentsObj AgentsObj
		err = json.Unmarshal([]byte(agentJson), &agentsObj)
		if err != nil {
			fmt.Println("service ListAgents() Unmarshal err=", err)
			return nil, err
		}
		agentsObjList = append(agentsObjList, agentsObj)
	}

	return agentsObjList, nil
}
