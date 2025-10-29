package svc

import (
	"jxzy/apis/api_knowledge/internal/config"
	"jxzy/bll/bll_context/bll_context"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	BllContextRpc bll_context.BllContextServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		BllContextRpc: bll_context.NewBllContextServiceClient(zrpc.MustNewClient(c.BllContextRpc).Conn()),
	}
}
