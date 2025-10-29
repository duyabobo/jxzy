package svc

import (
	"context"
	"jxzy/bll/bll_context/internal/config"
	"jxzy/bll/bll_context/internal/model"
	"jxzy/bs/bs_llm/bsllmservice"
	"jxzy/bs/bs_rag/bsragservice"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	ChatSessionModel    model.ChatSessionModel
	ChatSessionQasModel model.ChatSessionQasModel
	LLMRpc              bsllmservice.BsLlmService
	RagRpc              bsragservice.BsRagService
	logger              logx.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	var sessionModel model.ChatSessionModel
	var sessionQasModel model.ChatSessionQasModel
	var llmRpc bsllmservice.BsLlmService
	var ragRpc bsragservice.BsRagService

	logger := logx.WithContext(context.Background())

	// 初始化数据库连接
	if c.MySQL.DataSource != "" {
		conn := sqlx.NewMysql(c.MySQL.DataSource)
		sessionModel = model.NewChatSessionModel(conn)
		sessionQasModel = model.NewChatSessionQasModel(conn)
		logger.Info("Successfully connected to MySQL")
	}

	// 初始化LLM RPC客户端（直连方式）
	if c.BsLlmRpc.Target != "" {
		client, err := zrpc.NewClient(c.BsLlmRpc)
		if err == nil {
			llmRpc = bsllmservice.NewBsLlmService(client)
			logger.Info("Successfully connected to LLM RPC via direct connection")
		} else {
			logger.Errorf("Failed to connect to LLM RPC: %v", err)
		}
	} else if len(c.LLMRpc.Etcd.Hosts) > 0 {
		// 兼容旧的etcd配置
		client, err := zrpc.NewClient(c.LLMRpc)
		if err == nil {
			llmRpc = bsllmservice.NewBsLlmService(client)
			logger.Info("Successfully connected to LLM RPC via etcd")
		} else {
			logger.Errorf("Failed to connect to LLM RPC: %v", err)
		}
	}

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
		Config:              c,
		ChatSessionModel:    sessionModel,
		ChatSessionQasModel: sessionQasModel,
		LLMRpc:              llmRpc,
		RagRpc:              ragRpc,
		logger:              logger,
	}
}
