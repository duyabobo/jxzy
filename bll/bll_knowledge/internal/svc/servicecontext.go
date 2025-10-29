package svc

import (
	"context"
	"jxzy/bll/bll_knowledge/internal/config"
	"jxzy/bs/bs_rag/bsragservice"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	RagRpc bsragservice.BsRagService
	logger logx.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	var ragRpc bsragservice.BsRagService
	logger := logx.WithContext(context.Background())

	// 初始化RAG RPC客户端
	if c.BsRagRpc.Target != "" {
		client, err := zrpc.NewClient(c.BsRagRpc)
		if err == nil {
			ragRpc = bsragservice.NewBsRagService(client)
			logger.Info("Successfully connected to RAG RPC via direct connection")
		} else {
			logger.Errorf("Failed to connect to RAG RPC: %v", err)
		}
	}

	return &ServiceContext{
		Config: c,
		RagRpc: ragRpc,
		logger: logger,
	}
}
