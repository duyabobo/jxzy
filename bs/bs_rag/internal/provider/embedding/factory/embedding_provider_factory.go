package factory

import (
	"fmt"
	"jxzy/bs/bs_rag/internal/config"
	bailian "jxzy/bs/bs_rag/internal/provider/embedding/bailian"
	etypes "jxzy/bs/bs_rag/internal/provider/embedding/types"
)

// EmbeddingProviderFactory 工厂：创建嵌入模型提供者
type EmbeddingProviderFactory struct{}

// BailianProviderConfig 百炼provider配置（包含model_code和vector_dimension）
type BailianProviderConfig struct {
	APIKey          string
	ModelCode       string
	VectorDimension int64
}

// NewProvider 根据类型与配置创建嵌入模型提供者
func (f *EmbeddingProviderFactory) NewProvider(t etypes.EmbeddingProviderType, cfg interface{}) (etypes.EmbeddingProvider, error) {
	switch t {
	case etypes.EmbeddingProviderTypeBailian:
		var c BailianProviderConfig
		if cfg != nil {
			if typed, ok := cfg.(BailianProviderConfig); ok {
				c = typed
			} else if typed, ok := cfg.(config.BailianConfig); ok {
				// 向后兼容：如果传入的是BailianConfig，转换为BailianProviderConfig
				c = BailianProviderConfig{
					APIKey:          typed.APIKey,
					ModelCode:       "",
					VectorDimension: 0,
				}
			} else {
				return nil, fmt.Errorf("invalid config type for bailian embedding provider")
			}
		}
		return bailian.NewBailianEmbeddingProvider(c.APIKey, c.ModelCode, c.VectorDimension), nil
	default:
		return nil, fmt.Errorf("unsupported embedding provider type: %s", t)
	}
}
