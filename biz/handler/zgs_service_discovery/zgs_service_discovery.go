// Code generated by hertz generator.

package zgs_service_discovery

import (
	"context"
	zgs_service_discovery "github.com/Albatrossiun/zgs_service_discovery/biz/model/zgs_service_discovery"
	"github.com/Albatrossiun/zgs_service_discovery/biz/service"
	"github.com/cloudwego/hertz/pkg/app"
)

var (
	srv *service.UserService
)

func InitUserServiceHandler() {
	srv = service.NewUserService()
}

// Regist .
// @router /regist [POST]
func Regist(ctx context.Context, c *app.RequestContext) {
	var err error
	var req zgs_service_discovery.RegistRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	resp := srv.Regist(ctx, req)

	c.JSON(200, resp)
}

// ListAgents .
// @router /list_agents [POST]
func ListAgents(ctx context.Context, c *app.RequestContext) {
	var err error
	var req zgs_service_discovery.ListAgentsByGroupAndStatusRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	resp := srv.ListAgentsByGroupAndStatus(ctx, req)

	c.JSON(200, resp)
}
