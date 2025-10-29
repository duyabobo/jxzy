package svc

import (
	"jxzy/apis/api_chat/internal/config"
	"jxzy/apis/api_chat/internal/middleware"
	"jxzy/bll/bll_context/bllcontextservice"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	Cors          rest.Middleware
	BllContextRpc bllcontextservice.BllContextService
}

func NewServiceContext(c config.Config) *ServiceContext {
	var bllContextRpc bllcontextservice.BllContextService

	// 初始化bll_context RPC客户端（优先直连方式）
	if c.BllContextRpc.Target != "" {
		client, err := zrpc.NewClient(c.BllContextRpc)
		if err == nil {
			bllContextRpc = bllcontextservice.NewBllContextService(client)
		}
	} else if len(c.BllContextRpc.Etcd.Hosts) > 0 {
		// 兼容旧的etcd配置
		client, err := zrpc.NewClient(c.BllContextRpc)
		if err == nil {
			bllContextRpc = bllcontextservice.NewBllContextService(client)
		}
	}

	return &ServiceContext{
		Config:        c,
		Cors:          middleware.NewCorsMiddleware().Handle,
		BllContextRpc: bllContextRpc,
	}
}
