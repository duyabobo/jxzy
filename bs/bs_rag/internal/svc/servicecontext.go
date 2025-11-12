package svc

import (
	"context"
	"fmt"
	"jxzy/bs/bs_rag/internal/config"
	"jxzy/bs/bs_rag/internal/model"
	efactory "jxzy/bs/bs_rag/internal/provider/embedding/factory"
	etypes "jxzy/bs/bs_rag/internal/provider/embedding/types"
	vfactory "jxzy/bs/bs_rag/internal/provider/vectorstore/factory"
	vtypes "jxzy/bs/bs_rag/internal/provider/vectorstore/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	VectorProvider      vtypes.VectorProvider
	EmbeddingSceneModel model.EmbeddingSceneModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建向量数据库工厂
	factory := &vfactory.VectorProviderFactory{}

	// 根据配置创建向量数据库提供者
	var vectorProvider vtypes.VectorProvider
	var err error

	switch c.VectorDB.Type {
	case "faiss":
		vectorProvider, err = factory.NewVectorProvider(vtypes.VectorProviderTypeFaiss, c.Faiss)
	case "dashvector":
		vectorProvider, err = factory.NewVectorProvider(vtypes.VectorProviderTypeDashVector, c.DashVector)
	case "mock":
		vectorProvider, err = factory.NewVectorProvider(vtypes.VectorProviderTypeMock, nil)
	default:
		// 默认使用 mock 提供者
		vectorProvider, err = factory.NewVectorProvider(vtypes.VectorProviderTypeMock, nil)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to create vector provider: %v", err))
	}

	// 初始化数据库连接和embedding_scene模型
	var embeddingSceneModel model.EmbeddingSceneModel
	if c.MySQL.DataSource != "" {
		conn := sqlx.NewMysql(c.MySQL.DataSource)
		embeddingSceneModel = model.NewEmbeddingSceneModel(conn)
	}

	return &ServiceContext{
		Config:              c,
		VectorProvider:      vectorProvider,
		EmbeddingSceneModel: embeddingSceneModel,
	}
}

// GetEmbeddingProvider 根据scene_code获取embedding provider
func (s *ServiceContext) GetEmbeddingProvider(ctx context.Context, sceneCode string) (etypes.EmbeddingProvider, int64, error) {
	// 查询embedding_scene表
	scene, err := s.EmbeddingSceneModel.FindOneBySceneCode(ctx, sceneCode)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find embedding scene by scene_code %s: %w", sceneCode, err)
	}

	// 从配置中获取provider配置
	providerConfig, ok := s.Config.EmbeddingProviders[scene.ProviderCode]
	if !ok {
		return nil, 0, fmt.Errorf("provider_code %s not found in config", scene.ProviderCode)
	}

	// 创建embedding provider
	eFactory := &efactory.EmbeddingProviderFactory{}
	var embeddingProvider etypes.EmbeddingProvider

	switch scene.ProviderCode {
	case "bailian":
		embeddingProvider, err = eFactory.NewProvider(etypes.EmbeddingProviderTypeBailian, efactory.BailianProviderConfig{
			APIKey:          providerConfig.APIKey,
			ModelCode:       scene.ModelCode,
			VectorDimension: scene.VectorDimension,
		})
	default:
		return nil, 0, fmt.Errorf("unsupported provider_code: %s", scene.ProviderCode)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to create embedding provider: %w", err)
	}

	return embeddingProvider, scene.VectorDimension, nil
}
