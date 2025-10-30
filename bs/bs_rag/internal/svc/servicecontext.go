package svc

import (
    "fmt"
    "jxzy/bs/bs_rag/internal/config"
    efactory "jxzy/bs/bs_rag/internal/provider/embedding/factory"
    etypes "jxzy/bs/bs_rag/internal/provider/embedding/types"
    vfactory "jxzy/bs/bs_rag/internal/provider/vectorstore/factory"
    vtypes "jxzy/bs/bs_rag/internal/provider/vectorstore/types"
)

type ServiceContext struct {
	Config           config.Config
    VectorProvider   vtypes.VectorProvider
    EmbeddingService etypes.EmbeddingProvider
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

    // 初始化 Embedding Provider（目前默认使用 Bailian）
    eFactory := &efactory.EmbeddingProviderFactory{}
    embeddingService, err := eFactory.NewProvider(etypes.EmbeddingProviderTypeBailian, c.Bailian)
    if err != nil {
        panic(fmt.Sprintf("Failed to create embedding provider: %v", err))
    }

	return &ServiceContext{
		Config:           c,
		VectorProvider:   vectorProvider,
		EmbeddingService: embeddingService,
	}
}
