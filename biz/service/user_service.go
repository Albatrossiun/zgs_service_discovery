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
		fmt.Println("service Regist Marshal err = ", err)
		return zgs_service_discovery.RegistResponse{
			Code:    500,
			Message: err.Error(),
		}
	}

	// 获取redis所有的agents
	agentsObjList, err := u.ListAgentsWithUnmarshal()
	if err != nil {
		fmt.Println("service u.GetAllAgentsInRedis() err = ", err)
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
		fmt.Println("service Regist err = ", err)
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
		fmt.Println("service ListAgents() err = ", err)
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

	// setStatus
	agentsStatus := req.Status
	agentsGroup := req.Group
	type void struct{}
	var member void

	setStatus := make(map[string]void)
	for _, statusStr := range agentsStatus {
		setStatus[statusStr] = member
	}

	setGroup := make(map[string]void)
	for _, groupStr := range agentsGroup {
		setGroup[groupStr] = member
	}

	var agentsList []AgentsObj
	for _, obj := range agentsObjList {
		if len(req.Group) == 0 && len(req.Status) == 0 {
			agentsList = append(agentsList, obj)
		} else if len(req.Group) == 0 {
			if _, exists := setStatus[obj.Status]; exists {
				agentsList = append(agentsList, obj)
			}
		} else if len(req.Status) == 0 {
			if _, exists := setGroup[obj.Group]; exists {
				agentsList = append(agentsList, obj)
			}
		} else {
			_, exists1 := setGroup[obj.Group]
			_, exists2 := setStatus[obj.Status]
			if exists1 && exists2 {
				agentsList = append(agentsList, obj)
			}
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
		fmt.Println("service GetAllAgentsInRedis() err = ", err)
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
			fmt.Println("service ListAgents() Unmarshal err = ", err)
			return nil, err
		}
		agentsObjList = append(agentsObjList, agentsObj)
	}

	return agentsObjList, nil
}

func (u *UserService) UpdateAgentExt(ctx context.Context, req zgs_service_discovery.UpdateAgentExtRequest) zgs_service_discovery.UpdateAgentExtResponse {
	agentsObjList, err := u.ListAgentsWithUnmarshal()
	if err != nil {
		fmt.Println("service u.GetAllAgentsInRedis() err = ", err)
		return zgs_service_discovery.UpdateAgentExtResponse{
			Code:    500,
			Message: err.Error(),
		}
	}
	for _, obj := range agentsObjList {
		if obj.UUid == req.UUID && obj.Ext != req.Ext {
			if req.Ext != "busy" && req.Ext != "free" && req.Ext != "" {
				fmt.Println("service u.GetAllAgentsInRedis() req.Ext is false")
				return zgs_service_discovery.UpdateAgentExtResponse{
					Code:    500,
					Message: "the value of agent ext can only be BUSY, FREE or empty",
				}
			}
			agentsObj := &AgentsObj{
				UUid: req.UUID,
				Ext:  req.Ext,
			}
			agentsObjJson, err := json.Marshal(agentsObj)
			if err != nil {
				fmt.Println("service UpdateAgentExt Marshal err = ", err)
				return zgs_service_discovery.UpdateAgentExtResponse{
					Code:    500,
					Message: err.Error(),
				}
			}
			err = u.userDomain.Regist("service_"+obj.UUid, string(agentsObjJson))
			if err != nil {
				fmt.Println("service Regist err = ", err)
				return zgs_service_discovery.UpdateAgentExtResponse{
					Code:    500,
					Message: err.Error(),
				}
			}
		}
	}
	return zgs_service_discovery.UpdateAgentExtResponse{Code: 200, Message: "update agent ext successfully"}
}

func (u *UserService) UpdateOnlineAgentsGroup(ctx context.Context, req zgs_service_discovery.UpdateOnlineAgentsGroupRequest) zgs_service_discovery.UpdateOnlineAgentsGroupResponse {
	agentsObjList, err := u.ListAgentsWithUnmarshal()
	if err != nil {
		fmt.Println("service u.GetAllAgentsInRedis() err = ", err)
		return zgs_service_discovery.UpdateOnlineAgentsGroupResponse{
			Code:    500,
			Message: err.Error(),
		}
	}

	uuids := make([]string, 0)
	uuidsMap := make(map[string]AgentsObj)
	for _, uid := range req.Uuids {
		uuidsMap[uid] = AgentsObj{}
	}

	for _, obj := range agentsObjList {
		if _, exists := uuidsMap[obj.UUid]; exists {
			if obj.Status != "online" && obj.Ext != "free" {
				return zgs_service_discovery.UpdateOnlineAgentsGroupResponse{
					Code:    500,
					Message: "there is an agent whose state is not ONLINE or ext is not FREE",
				}
			}
			if obj.Group == req.Group {
				return zgs_service_discovery.UpdateOnlineAgentsGroupResponse{
					Code:    200,
					Message: "this is consistent with the source data and does not need to be modified",
				}
			}
			uuids = append(uuids, obj.UUid)
		}
	}

	agentsObjJson, err := u.userDomain.GetAgentsByUUids(uuids)
	if err != nil {
		fmt.Println("service GetAgentsByUUids err = ", err)
		return zgs_service_discovery.UpdateOnlineAgentsGroupResponse{
			Code:    500,
			Message: err.Error(),
		}
	}
	for _, obj := range agentsObjJson {
		var objStruct AgentsObj
		err = json.Unmarshal([]byte(obj), &objStruct)
		if err != nil {
			return zgs_service_discovery.UpdateOnlineAgentsGroupResponse{
				Code:    500,
				Message: err.Error(),
			}
		}
		objStruct.Group = req.Group
		agentsObjString, err := json.Marshal(objStruct)
		err = u.userDomain.Regist("service_"+objStruct.UUid, string(agentsObjString))
		if err != nil {
			return zgs_service_discovery.UpdateOnlineAgentsGroupResponse{
				Code:    500,
				Message: err.Error(),
			}
		}
	}

	return zgs_service_discovery.UpdateOnlineAgentsGroupResponse{Code: 200, Message: "update online agents group successfully"}
}
