package consts

// 向量维度常量定义
const (
	// Faiss 向量数据库默认维度
	FaissDefaultDimension = 1024

	// DashVector 向量数据库默认维度
	DashVectorDefaultDimension = 1024
)

// 向量索引类型常量
const (
	// Faiss 索引类型
	IndexTypeIVFFlat = "IVFFlat"
	IndexTypeFlat    = "Flat"
	IndexTypeHNSW    = "HNSW"
)

// 向量距离度量类型常量
const (
	MetricTypeL2     = "L2"     // 欧几里得距离
	MetricTypeIP     = "IP"     // 内积
	MetricTypeCOSINE = "COSINE" // 余弦相似度
)

// 向量数据库类型常量
const (
	VectorDBTypeFaiss      = "faiss"
	VectorDBTypeDashVector = "dashvector"
	VectorDBTypeMilvus     = "milvus"
	VectorDBTypePinecone   = "pinecone"
	VectorDBTypeWeaviate   = "weaviate"
	VectorDBTypeMock       = "mock"
)

// RAG服务相关常量
const (
	// 默认集合名称
	DefaultCollectionName = "test_collection"

	// 集合配置常量
	MaxCollections            = 100
	MaxDocumentsPerCollection = 1000000
)
