package domain

import (
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/dao"
)

type UserDomain struct {
	userRepository dao.UserRepository
}

func NewUserDomain() *UserDomain {
	return &UserDomain{
		userRepository: dao.GetUserRepository(),
	}
}

func (u *UserDomain) Regist(uuid, agentsObjJson string) error {
	err := u.userRepository.Regist(uuid, agentsObjJson)
	if err != nil {
		fmt.Println("domain Regist err = ", err)
		return err
	}
	return nil
}

func (u *UserDomain) ListAgents() ([]string, error) {
	agentsList, err := u.userRepository.ListAgents()
	if err != nil {
		fmt.Println("domain ListAgents err = ", err)
		return nil, err
	}
	return agentsList, nil
}

func (u *UserDomain) DeleteAgents(uuid string) error {
	err := u.userRepository.DeleteAgents(uuid)
	if err != nil {
		fmt.Println("domain DeleteAgents err = ", err)
		return err
	}
	return nil
}

func (u *UserDomain) GetAgentsByUUids(uuids []string) ([]string, error) {
	agentsList, err := u.userRepository.GetAgentsByUUids(uuids)
	if err != nil {
		fmt.Println("domain GetAgentsByUUids err = ", err)
		return nil, err
	}
	return agentsList, nil
}
