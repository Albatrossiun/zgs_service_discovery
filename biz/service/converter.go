package service

import (
	"github.com/Albatrossiun/zgs_service_discovery/biz/model/zgs_service_discovery"
)

func AgentsObjToAgentInfo(obj AgentsObj) zgs_service_discovery.AgentInfo {
	return zgs_service_discovery.AgentInfo{
		UUID:   obj.UUid,
		IP:     obj.Ip,
		Port:   obj.Port,
		Status: obj.Status,
		Group:  obj.Group,
		Ext:    obj.Ext,
	}
}
