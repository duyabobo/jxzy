# JXZY BLL Context 服务

业务逻辑层服务，负责聊天会话管理和上下文处理。

## 📁 项目结构

```
bll_context/
├── bllcontext.proto            # Proto 定义文件
├── bll_context/                # 生成的 proto 代码
├── bllcontextservice/          # 服务实现
├── etc/
│   └── bllcontext.yaml        # 配置文件
├── internal/
│   ├── common/                # 公共模块（向量化服务等）
│   ├── config/                # 配置相关
│   ├── logic/                 # 业务逻辑
│   ├── model/                 # 数据模型
│   ├── server/                # 服务端实现
│   └── svc/                   # 服务上下文
├── scripts/                   # 脚本目录
│   ├── compile.sh             # 编译脚本
│   └── restart.sh             # 重启脚本
└── README.md                  # 本文件
```

## 🚀 快速开始

### 前置要求

1. **goctl**: go-zero 代码生成工具
   ```bash
   go install github.com/zeromicro/go-zero/tools/goctl@latest
   ```

2. **protoc**: Protocol Buffers 编译器
   ```bash
   # macOS
   brew install protobuf
   
   # Ubuntu/Debian
   sudo apt-get install protobuf-compiler
   ```

3. **Go**: 1.16 或更高版本

### 编译和运行

#### 方法一：使用脚本（推荐）
```bash
# 编译项目
./scripts/compile.sh

# 重启服务
./scripts/restart.sh
```

#### 方法二：手动操作
```bash
# 重新生成 proto 代码
goctl rpc protoc bllcontext.proto --go_out=. --go-grpc_out=. --zrpc_out=.

# 编译项目
go build -o bll_context .

# 启动服务
./bll_context -f etc/bllcontext.yaml
```

## 📋 服务接口

### StreamChat RPC 接口

- **服务名**: `BllContextService`
- **方法**: `StreamChat`
- **协议**: gRPC

#### 请求参数
```protobuf
message StreamChatRequest {
    string message = 1;
    string session_id = 2;
    string scene_code = 3;
}
```

#### 响应格式
```protobuf
message StreamChatResponse {
    string session_id = 1;
    string scene_code = 2;
    string delta = 3;
    bool finished = 4;
    TokenUsage usage = 5;
}
```

## 🔧 配置说明

配置文件：`etc/bllcontext.yaml`

## 🧠 向量化功能

### EmbeddingService

项目集成了阿里云百炼的 Embedding API，用于文本向量化处理。

#### 功能特性

- **高质量向量**: 使用 `text-embedding-v4` 模型生成 1024 维向量
- **多语言支持**: 支持中文、英语、西班牙语、法语等 100+ 主流语种
- **可复用设计**: 封装在 `internal/common/embedding.go` 中，可在多个场景下复用
- **错误处理**: 完善的错误处理和日志记录

#### 使用方式

```go
// 创建 EmbeddingService 实例
embeddingService := common.NewEmbeddingService()

// 生成文本向量
vector, err := embeddingService.GenerateEmbedding("要向量化的文本")
if err != nil {
    // 处理错误
    return err
}

// 使用生成的向量进行 RAG 检索或其他操作
```

#### 配置要求

1. **环境变量配置**:
   ```bash
   export BAILIAN_API_KEY="your-bailian-api-key"
   ```

2. **API 限制**:
   - 单行最大处理 Token 数：8,192
   - 支持批量处理：最多 10 行
   - 向量维度：支持 2,048、1,536、1,024（默认）、768、512、256、128、64

#### 应用场景

- **RAG 检索**: 为关键句生成向量，用于相似度搜索
- **向量插入**: 将文档内容向量化后存储到向量数据库
- **语义搜索**: 基于向量相似度进行语义匹配

## 🚀 并发 RAG 搜索

### 功能特性

- **并发处理**: 使用 `errgroup` 实现协程并发搜索
- **性能优化**: 多个关键句同时进行 RAG 检索，显著提升响应速度
- **错误隔离**: 单个搜索失败不会影响其他搜索
- **线程安全**: 使用互斥锁保护共享资源
- **上下文控制**: 支持请求取消和超时控制

### 实现原理

```go
// 使用 errgroup 并发处理多个关键句
g, ctx := errgroup.WithContext(l.ctx)

for _, sentence := range keySentences {
    g.Go(func() error {
        // 并发执行 RAG 搜索
        ragResp, err := l.svcCtx.RagRpc.VectorSearch(ctx, ragReq)
        // 处理结果...
        return nil
    })
}

// 等待所有协程完成
if err := g.Wait(); err != nil {
    // 处理错误
}
```

### 性能优势

- **并行搜索**: 多个关键句同时搜索，而不是串行处理
- **资源利用**: 充分利用网络 I/O 和 CPU 资源
- **响应时间**: 大幅减少总体响应时间
- **可扩展性**: 易于调整并发数量和处理策略

```yaml
Name: bll_context
Host: 0.0.0.0
Port: 8080

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: bll_context.rpc

DataSource:
  Host: 127.0.0.1
  Port: 3306
  User: root
  Password: password
  Database: jxzy
```

## 📝 开发指南

### 修改 Proto 定义

1. 编辑 `bllcontext.proto` 文件
2. 运行 `./scripts/compile.sh` 重新生成代码
3. 运行 `./scripts/restart.sh` 重启服务

### 添加新的 RPC 方法

1. 在 `bllcontext.proto` 中定义新的方法
2. 运行 `./scripts/compile.sh` 生成代码框架
3. 在 `internal/logic/` 中实现业务逻辑
4. 重启服务

## 🚨 注意事项

1. **服务依赖**: 需要 etcd 和 MySQL 数据库
2. **端口配置**: 默认使用 8080 端口
3. **日志文件**: 服务日志保存在 `bll_context.log`
4. **数据库**: 需要初始化数据库表结构

## 🔧 故障排除

### 编译失败
- 检查 goctl 和 protoc 是否正确安装
- 检查 proto 文件语法是否正确

### 启动失败
- 检查配置文件是否存在
- 检查端口是否被占用
- 检查数据库连接是否正常
- 查看日志文件获取详细错误信息

### 连接失败
- 检查 etcd 服务是否运行
- 检查数据库服务是否可用

## 📞 服务信息

- **服务类型**: gRPC RPC 服务
- **注册地址**: 127.0.0.1:8080
- **日志文件**: `bll_context.log`
- **配置文件**: `etc/bllcontext.yaml`
