package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MySQL       MySQLConfig              `json:"MySQL"`
	Faiss       FaissConfig              `json:"Faiss"`
	DashVector  DashVectorConfig         `json:"DashVector"`
	Collections CollectionsConfig        `json:"Collections"`
	VectorDB    VectorDBConfig           `json:"VectorDB"`
	EmbeddingProviders map[string]EmbeddingProviderConfig `json:"EmbeddingProviders"` // provider_code -> provider config
}

type VectorDBConfig struct {
	Type   string                 `json:"Type"`   // 向量数据库类型: faiss, milvus, pinecone, weaviate, mock
	Config map[string]interface{} `json:"Config"` // 具体配置
}

type MySQLConfig struct {
	DataSource string `json:"DataSource"`
}

type FaissConfig struct {
	IndexPath      string `json:"IndexPath"`
	IndexType      string `json:"IndexType"`
	Nlist          int    `json:"Nlist"`
	Nprobe         int    `json:"Nprobe"`
	M              int    `json:"M"`
	EfConstruction int    `json:"EfConstruction"`
	EfSearch       int    `json:"EfSearch"`
	MetricType     string `json:"MetricType"`
}

type DashVectorConfig struct {
	Endpoint string            `json:"Endpoint"` // DashVector 服务端点
	APIKey   string            `json:"APIKey"`   // API 密钥
	Region   string            `json:"Region"`   // 地域
	Timeout  int               `json:"Timeout"`  // 请求超时时间（秒）
	Headers  map[string]string `json:"Headers"`  // 自定义请求头
}

type CollectionsConfig struct {
	MaxCollections            int `json:"MaxCollections"`
	MaxDocumentsPerCollection int `json:"MaxDocumentsPerCollection"`
}

// EmbeddingProviderConfig 嵌入模型提供商配置
type EmbeddingProviderConfig struct {
	APIKey string `json:"APIKey"` // API密钥
	// 可以添加其他提供商特定的配置字段
}

// BailianConfig 百炼配置（用于向后兼容factory）
type BailianConfig struct {
	APIKey string `json:"APIKey"`
}
