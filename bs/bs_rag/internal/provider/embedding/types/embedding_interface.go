package types

// EmbeddingProvider 向量化提供者接口
type EmbeddingProvider interface {
    // GenerateEmbedding 根据文本生成向量表示
    GenerateEmbedding(text string) ([]float32, error)
}

// EmbeddingProviderType 枚举可用的嵌入模型提供者
type EmbeddingProviderType string

const (
    // EmbeddingProviderTypeBailian 阿里云百炼嵌入模型
    EmbeddingProviderTypeBailian EmbeddingProviderType = "bailian"
)


