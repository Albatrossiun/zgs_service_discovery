package backstage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/domain"
	"github.com/Albatrossiun/zgs_service_discovery/biz/service"
	"net/http"
	"time"
)

var (
	ctx         context.Context
	userService *service.UserService
	userDomain  *domain.UserDomain
)

func StartHeartbeat() {
	go func() {
		// 初始化
		userService = service.NewUserService()
		userDomain = domain.NewUserDomain()
		ctx = context.Background()

		// 定时任务
		for {
			agentsObjList, err := userService.ListAgentsWithUnmarshal()
			if err != nil {
				fmt.Println("StartHeartbeat() err = ", err)
				return
			}
			if len(agentsObjList) == 0 {
				fmt.Println("StartHeartbeat() agentsObjList is empty")
			} else {
				for _, agent := range agentsObjList {
					if agent.Group == "online" {
						url := agent.Ip + ":" + agent.Port + "/status"
						resp, err := http.Get(url)
						if err != nil {
							fmt.Println("StartHeartbeat() http.Get(url) err = ", err)
							continue
						}
						if resp.Status == "200" {
							continue
						} else {
							// 修改redis的对应status = offline
							agentsObj := &service.AgentsObj{
								UUid:   agent.UUid,
								Ip:     agent.Ip,
								Port:   agent.Port,
								Status: "offline",
							}
							err = updateAgentStatusByStruct(agentsObj)
							if err != nil {
								fmt.Println("StartHeartbeat() updateAgentStatusByStruct err = ", err)
							}
						}
					} else if agent.Group == "offline" {
						url := agent.Ip + ":" + agent.Port + "/status"
						resp, err := http.Get(url)
						if err != nil {
							fmt.Println("StartHeartbeat() http.Get(url) err = ", err)
							continue
						}
						if resp.Status == "200" {
							// 修改redis的对应status = online
							agentsObj := &service.AgentsObj{
								UUid:   agent.UUid,
								Ip:     agent.Ip,
								Port:   agent.Port,
								Status: "online",
							}
							err = updateAgentStatusByStruct(agentsObj)
							if err != nil {
								fmt.Println("StartHeartbeat() updateAgentStatusByStruct err = ", err)
							}
						} else {
							// 在redis中删除数据
							err = userDomain.DeleteAgents(agent.UUid)
							if err != nil {
								fmt.Println("StartHeartbeat() DeleteAgents err = ", err)
							}
						}
					} else {
						fmt.Println("StartHeartbeat() status error，delete the data directly")
						// 在redis中删除数据
						err = userDomain.DeleteAgents(agent.UUid)
						if err != nil {
							fmt.Println("StartHeartbeat() DeleteAgents err = ", err)
						}
					}
				}
			}

			time.Sleep(5)
		}
	}()
}

func updateAgentStatusByStruct(obj *service.AgentsObj) error {
	agentsObjJson, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("StartHeartbeat() Marshal err = ", err)
		return err
	}
	err = userDomain.Regist(obj.UUid, string(agentsObjJson))
	if err != nil {
		fmt.Println("StartHeartbeat() Regist err = ", err)
		return err
	}
	return nil
}
