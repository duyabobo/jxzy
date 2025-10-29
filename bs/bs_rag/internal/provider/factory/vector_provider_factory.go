package factory

import (
	"jxzy/bs/bs_rag/internal/config"
	"jxzy/bs/bs_rag/internal/provider/dashvector"
	"jxzy/bs/bs_rag/internal/provider/faiss"
	"jxzy/bs/bs_rag/internal/provider/types"
)

// VectorProviderFactory 向量数据库工厂
type VectorProviderFactory struct{}

// NewVectorProvider 创建向量数据库提供者
func (f *VectorProviderFactory) NewVectorProvider(providerType types.VectorProviderType, cfg interface{}) (types.VectorProvider, error) {
	switch providerType {
	case types.VectorProviderTypeFaiss:
		if faissConfig, ok := cfg.(config.FaissConfig); ok {
			return faiss.NewFaissProvider(faissConfig), nil
		}
		return nil, types.ErrInvalidConfig
	case types.VectorProviderTypeDashVector:
		if dashVectorConfig, ok := cfg.(config.DashVectorConfig); ok {
			dashConfig := dashvector.DashVectorConfig{
				Endpoint: dashVectorConfig.Endpoint,
				APIKey:   dashVectorConfig.APIKey,
				Region:   dashVectorConfig.Region,
				Timeout:  dashVectorConfig.Timeout,
				Headers:  dashVectorConfig.Headers,
			}
			return dashvector.NewDashVectorProvider(dashConfig), nil
		}
		return nil, types.ErrInvalidConfig
	default:
		return nil, types.ErrUnsupportedProvider
	}
}
