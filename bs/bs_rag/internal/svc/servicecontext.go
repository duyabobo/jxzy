package svc

import (
	"fmt"
	"jxzy/bs/bs_rag/internal/config"
	"jxzy/bs/bs_rag/internal/provider/factory"
	"jxzy/bs/bs_rag/internal/provider/types"
)

type ServiceContext struct {
	Config         config.Config
	VectorProvider types.VectorProvider
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建向量数据库工厂
	factory := &factory.VectorProviderFactory{}

	// 根据配置创建向量数据库提供者
	var vectorProvider types.VectorProvider
	var err error

	switch c.VectorDB.Type {
	case "faiss":
		vectorProvider, err = factory.NewVectorProvider(types.VectorProviderTypeFaiss, c.Faiss)
	case "dashvector":
		vectorProvider, err = factory.NewVectorProvider(types.VectorProviderTypeDashVector, c.DashVector)
	case "mock":
		vectorProvider, err = factory.NewVectorProvider(types.VectorProviderTypeMock, nil)
	default:
		// 默认使用 mock 提供者
		vectorProvider, err = factory.NewVectorProvider(types.VectorProviderTypeMock, nil)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to create vector provider: %v", err))
	}

	return &ServiceContext{
		Config:         c,
		VectorProvider: vectorProvider,
	}
}
