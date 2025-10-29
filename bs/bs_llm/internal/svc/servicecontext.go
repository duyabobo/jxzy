package svc

import (
	"context"
	"jxzy/bs/bs_llm/internal/config"
	"jxzy/bs/bs_llm/internal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	LlmSceneModel      model.LlmSceneModel
	LlmCompletionModel model.LlmCompletionModel
	logger             logx.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	var sceneModel model.LlmSceneModel
	var completionModel model.LlmCompletionModel

	logger := logx.WithContext(context.Background())

	// 初始化数据库连接
	if c.MySQL.DataSource != "" {
		conn := sqlx.NewMysql(c.MySQL.DataSource)
		sceneModel = model.NewLlmSceneModel(conn)
		completionModel = model.NewLlmCompletionModel(conn)
		logger.Info("Successfully connected to MySQL")
	}

	return &ServiceContext{
		Config:             c,
		LlmSceneModel:      sceneModel,
		LlmCompletionModel: completionModel,
		logger:             logger,
	}
}
