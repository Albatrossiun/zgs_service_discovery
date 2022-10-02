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

func (u *UserDomain) Regist(uuid, ipAndPort string) error {
	err := u.userRepository.Regist(uuid, ipAndPort)
	if err != nil {
		fmt.Println("domain Regist err=", err)
		return err
	}
	return nil
}

func (u *UserDomain) ListAgents(uuid string) (string, error) {
	ipAndPort, err := u.userRepository.ListAgents(uuid)
	if err != nil {
		fmt.Println("domain ListAgents err=", err)
		return "", err
	}
	return ipAndPort, nil
}
