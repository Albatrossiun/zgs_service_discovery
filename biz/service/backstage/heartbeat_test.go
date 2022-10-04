package backstage

import (
	"fmt"
	"github.com/Albatrossiun/zgs_service_discovery/biz/dao"
	"github.com/Albatrossiun/zgs_service_discovery/biz/handler/zgs_service_discovery"
	"github.com/Albatrossiun/zgs_service_discovery/biz/storage"
	"testing"
)

func init() {
	err := storage.InitRedisPool()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	dao.InitRedis()
	dao.InitUserRepository()
	zgs_service_discovery.InitUserServiceHandler()
}

func TestStartHeartbeat(t *testing.T) {
	StartHeartbeat()
}
