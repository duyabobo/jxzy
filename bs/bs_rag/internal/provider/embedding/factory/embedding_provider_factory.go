package factory

import (
    "fmt"
    "jxzy/bs/bs_rag/internal/config"
    bailian "jxzy/bs/bs_rag/internal/provider/embedding/bailian"
    etypes "jxzy/bs/bs_rag/internal/provider/embedding/types"
)

// EmbeddingProviderFactory 工厂：创建嵌入模型提供者
type EmbeddingProviderFactory struct{}

// NewProvider 根据类型与配置创建嵌入模型提供者
func (f *EmbeddingProviderFactory) NewProvider(t etypes.EmbeddingProviderType, cfg interface{}) (etypes.EmbeddingProvider, error) {
    switch t {
    case etypes.EmbeddingProviderTypeBailian:
        var c config.BailianConfig
        if cfg != nil {
            if typed, ok := cfg.(config.BailianConfig); ok {
                c = typed
            } else {
                return nil, fmt.Errorf("invalid config type for bailian embedding provider")
            }
        }
        return bailian.NewBailianEmbeddingProvider(c.APIKey), nil
    default:
        return nil, fmt.Errorf("unsupported embedding provider type: %s", t)
    }
}


