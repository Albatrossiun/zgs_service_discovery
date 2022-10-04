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
					if agent.Status == "online" {
						url := "http://" + agent.Ip + ":" + agent.Port + "/status"
						resp, err := http.Get(url)
						fmt.Println(err)
						if err == nil && resp.Status == "200 OK" {
							continue
						} else {
							// 修改redis的对应status = offline
							agentsObj := &service.AgentsObj{
								UUid:        agent.UUid,
								Ip:          agent.Ip,
								Port:        agent.Port,
								Status:      "offline",
								OfflineTime: time.Now().Unix(),
							}
							err = updateAgentStatusByStruct(agentsObj)
							if err != nil {
								fmt.Println("StartHeartbeat() updateAgentStatusByStruct err = ", err)
							}
						}
					} else if agent.Status == "offline" {
						url := "http://" + agent.Ip + ":" + agent.Port + "/status"
						resp, err := http.Get(url)
						fmt.Println(err)
						if err == nil && resp.Status == "200 OK" {
							// 修改redis的对应status = online
							agentsObj := &service.AgentsObj{
								UUid:        agent.UUid,
								Ip:          agent.Ip,
								Port:        agent.Port,
								Status:      "online",
								OfflineTime: 0,
							}
							err = updateAgentStatusByStruct(agentsObj)
							if err != nil {
								fmt.Println("StartHeartbeat() updateAgentStatusByStruct err = ", err)
							}
						} else {
							if time.Now().Unix()-agent.OfflineTime > 10000 {
								// 在redis中删除数据
								err = userDomain.DeleteAgents("service_" + agent.UUid)
								if err != nil {
									fmt.Println("StartHeartbeat() DeleteAgents err = ", err)
								}
							}
						}
					} else {
						fmt.Println("StartHeartbeat() status error，delete the data directly")
						// 在redis中删除数据
						err = userDomain.DeleteAgents("service_" + agent.UUid)
						if err != nil {
							fmt.Println("StartHeartbeat() DeleteAgents err = ", err)
						}
					}
				}
			}

			time.Sleep(5 * time.Second) // 5秒
		}
	}()
}

func updateAgentStatusByStruct(obj *service.AgentsObj) error {
	agentsObjJson, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("StartHeartbeat() Marshal err = ", err)
		return err
	}

	err = userDomain.Regist("service_"+obj.UUid, string(agentsObjJson))
	if err != nil {
		fmt.Println("StartHeartbeat() Regist err = ", err)
		return err
	}
	return nil
}
