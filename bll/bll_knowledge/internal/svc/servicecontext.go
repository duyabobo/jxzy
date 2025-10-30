package svc

import (
	"context"
	"jxzy/bll/bll_knowledge/internal/config"
	"jxzy/bll/bll_knowledge/internal/model"
    "jxzy/bs/bs_llm/bsllmservice"
	"jxzy/bs/bs_rag/bsragservice"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	RagRpc                bsragservice.BsRagService
    LlmRpc                bsllmservice.BsLlmService
	KnowledgeFileModel    model.KnowledgeFileModel
	KnowledgeSegmentModel model.KnowledgeSegmentModel
    KnowledgeSummarySentenceModel model.KnowledgeSummarySentenceModel
	logger                logx.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	var ragRpc bsragservice.BsRagService
    var llmRpc bsllmservice.BsLlmService
	var knowledgeFileModel model.KnowledgeFileModel
	var knowledgeSegmentModel model.KnowledgeSegmentModel
    var knowledgeSummarySentenceModel model.KnowledgeSummarySentenceModel
	logger := logx.WithContext(context.Background())

	// 初始化数据库连接
    if c.MySQL.DataSource != "" {
		conn := sqlx.NewMysql(c.MySQL.DataSource)
		knowledgeFileModel = model.NewKnowledgeFileModel(conn)
		knowledgeSegmentModel = model.NewKnowledgeSegmentModel(conn)
        knowledgeSummarySentenceModel = model.NewKnowledgeSummarySentenceModel(conn, c.Cache)
		logger.Info("Successfully connected to MySQL")
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

    // 初始化LLM RPC客户端
    if c.BsLlmRpc.Target != "" {
        client, err := zrpc.NewClient(c.BsLlmRpc)
        if err == nil {
            llmRpc = bsllmservice.NewBsLlmService(client)
            logger.Info("Successfully connected to LLM RPC via direct connection")
        } else {
            logger.Errorf("Failed to connect to LLM RPC: %v", err)
        }
    }

	return &ServiceContext{
		Config:                c,
		RagRpc:                ragRpc,
        LlmRpc:                llmRpc,
		KnowledgeFileModel:    knowledgeFileModel,
		KnowledgeSegmentModel: knowledgeSegmentModel,
        KnowledgeSummarySentenceModel: knowledgeSummarySentenceModel,
		logger:                logger,
	}
}
