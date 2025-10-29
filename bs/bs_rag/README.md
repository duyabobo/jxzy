# BS RAG Service

基于 Faiss 向量数据库的 RAG (Retrieval-Augmented Generation) 服务，提供向量搜索、插入、删除等功能。

## 功能特性

- **向量相似度搜索**: 支持基于向量相似度的文档检索
- **向量文档管理**: 支持向量文档的插入、删除操作
- **集合管理**: 支持多个向量集合的管理和信息查询
- **高性能**: 基于 Faiss 向量数据库，支持高效的向量索引和搜索
- **可配置**: 支持多种索引类型和参数配置

## 服务接口

### 1. 向量搜索 (VectorSearch)
- **功能**: 根据查询向量搜索最相似的文档
- **参数**: 
  - `query_vector`: 查询向量
  - `top_k`: 返回结果数量
  - `min_score`: 最小相似度阈值
  - `collection_name`: 集合名称
  - `filters`: 过滤条件

### 2. 向量插入 (VectorInsert)
- **功能**: 向指定集合插入向量文档
- **参数**:
  - `collection_name`: 集合名称
  - `documents`: 要插入的文档列表
  - `user_id`: 用户ID

### 3. 向量删除 (VectorDelete)
- **功能**: 从指定集合删除向量文档
- **参数**:
  - `collection_name`: 集合名称
  - `document_ids`: 要删除的文档ID列表
  - `user_id`: 用户ID

### 4. 获取集合信息 (GetCollectionInfo)
- **功能**: 获取指定集合的详细信息
- **参数**:
  - `collection_name`: 集合名称
  - `user_id`: 用户ID

## 配置说明

### Faiss 配置
```yaml
Faiss:
  IndexPath: ./data/faiss_indexes    # Faiss索引文件存储路径
  DefaultDimension: 1536             # 默认向量维度 (DashVectorDefaultDimension)
  IndexType: "IVFFlat"               # 索引类型: IVFFlat, Flat, HNSW
  Nlist: 100                         # IVF索引的聚类中心数量
  Nprobe: 10                         # IVF搜索时的聚类中心数量
  M: 16                              # HNSW索引的层数
  EfConstruction: 200                # HNSW构建时的搜索深度
  EfSearch: 50                       # HNSW搜索时的搜索深度
  MetricType: "L2"                   # 距离度量类型: L2, IP(内积), COSINE
```

### 集合配置
```yaml
Collections:
  DefaultCollection: "default"        # 默认集合名称
  MaxCollections: 100                # 最大集合数量
  MaxDocumentsPerCollection: 1000000 # 每个集合最大文档数
```

## 快速开始

### 1. 编译服务
```bash
cd bs/bs_rag
chmod +x scripts/compile.sh
./scripts/compile.sh
```

### 2. 启动服务
```bash
chmod +x scripts/restart.sh
./scripts/restart.sh
```

### 3. 测试服务
```bash
# 使用 grpcurl 测试 (需要先安装 grpcurl)
grpcurl -plaintext -d '{
  "query_vector": [0.1, 0.2, 0.3],
  "top_k": 5,
  "collection_name": "test"
}' localhost:8082 bs_rag.BsRagService/VectorSearch
```

## 开发说明

### 项目结构
```
bs/bs_rag/
├── bsrag.proto              # Protocol Buffers 定义
├── bsrag.go                 # 服务入口文件
├── etc/
│   └── bsrag.yaml          # 配置文件
├── internal/
│   ├── config/             # 配置结构
│   ├── logic/              # 业务逻辑
│   ├── provider/           # 向量数据库提供者
│   │   ├── vector_provider.go  # 向量数据库接口定义
│   │   ├── faiss.go           # Faiss 实现
│   │   └── mock_provider.go   # Mock 实现
│   ├── server/             # RPC 服务器
│   ├── svc/                # 服务上下文
│   └── test/               # 单元测试
├── scripts/                # 脚本文件
└── README.md               # 说明文档
```

### 向量数据库接口设计

项目采用接口设计模式，支持多种向量数据库：

#### 1. VectorProvider 接口
```go
type VectorProvider interface {
    Search(ctx context.Context, collectionName string, queryVector []float32, topK int, minScore float32) ([]SearchResult, error)
    Insert(ctx context.Context, collectionName string, documents []Document) error
    Delete(ctx context.Context, collectionName string, documentIDs []string) error
    GetCollectionInfo(ctx context.Context, collectionName string) (*CollectionInfo, error)
    CreateCollection(ctx context.Context, collectionName string, dimension int, indexType string) error
    DeleteCollection(ctx context.Context, collectionName string) error
    ListCollections(ctx context.Context) ([]string, error)
    Close() error
}
```

#### 2. 支持的向量数据库类型
- **Faiss**: Facebook 开源的向量搜索库
- **DashVector**: 阿里云向量检索服务
- **Milvus**: 云原生向量数据库
- **Pinecone**: 托管的向量数据库服务
- **Weaviate**: 向量搜索引擎
- **Mock**: 用于测试的模拟实现

#### 3. 配置方式
```yaml
VectorDB:
  Type: "dashvector"         # 向量数据库类型
  Config: {}                 # 具体配置

# DashVector 具体配置
DashVector:
  Endpoint: "https://dashvector.cn-hangzhou.aliyuncs.com"  # 服务端点
  APIKey: "your-api-key-here"                              # API 密钥
  Region: "cn-hangzhou"                                    # 地域
  DefaultDimension: 1536                                   # 默认向量维度 (DashVectorDefaultDimension)
  Timeout: 30                                              # 请求超时时间（秒）
```

### 扩展新的向量数据库

要添加新的向量数据库支持：

1. 在 `provider/` 目录下创建新的实现文件（如 `milvus.go`）
2. 实现 `VectorProvider` 接口的所有方法
3. 在 `vector_provider.go` 中添加新的类型常量
4. 在 `VectorProviderFactory` 中添加创建逻辑
5. 在配置文件中添加相应的配置结构

### 当前实现状态

- ✅ **Mock Provider**: 完整的模拟实现，用于测试
- ✅ **DashVector Provider**: 阿里云向量检索服务，完整实现
- 🔄 **Faiss Provider**: 基础框架，需要集成真实的 Faiss 库
- ⏳ **其他 Provider**: 待实现

### 集成真实 Faiss

要集成真实的 Faiss，需要：

1. 添加 Faiss Go 绑定依赖
2. 在 `provider/faiss.go` 中实现真实的 Faiss 操作
3. 添加索引持久化和加载功能
4. 实现向量数据的存储和管理

### 添加新的 RPC 方法

1. 在 `bsrag.proto` 中定义新的消息和服务方法
2. 使用 `goctl` 重新生成代码
3. 在 `internal/logic/` 中实现业务逻辑
4. 在 `internal/provider/` 中添加相应的底层操作

## 注意事项

- 当前实现为模拟版本，实际使用时需要集成真实的 Faiss 库
- 向量维度需要在配置中正确设置
- 大量数据插入时建议使用批量操作
- 生产环境建议配置适当的日志级别和监控
