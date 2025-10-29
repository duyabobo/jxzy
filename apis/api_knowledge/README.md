# API Knowledge 服务

知识库管理API服务，提供知识库的添加和删除功能。

## 功能特性

- **添加知识库**: 将知识库内容向量化并存储到向量数据库
- **删除知识库**: 根据向量ID删除指定的知识库

## 服务配置

- **端口**: 8889
- **配置文件**: `etc/apiknowledge.yaml`
- **日志文件**: `logs/api_knowledge.log`

## API接口

### 1. 添加知识库

**接口**: `POST /api/v1/knowledge/add`

**请求参数**:
```json
{
    "summary": "知识库总结",
    "content": "知识库内容",
    "user_id": "用户ID"
}
```

**响应**:
```json
{
    "vector_id": "向量ID",
    "success": true,
    "message": "知识库添加成功"
}
```

### 2. 删除知识库

**接口**: `POST /api/v1/knowledge/delete`

**请求参数**:
```json
{
    "vector_id": "向量ID",
    "user_id": "用户ID"
}
```

**响应**:
```json
{
    "success": true,
    "message": "知识库删除成功"
}
```

## 快速开始

### 编译服务
```bash
bash scripts/compile.sh
```

### 启动服务
```bash
bash scripts/restart.sh
```

### 测试API
```bash
# 添加知识库
curl -X POST http://localhost:8889/api/v1/knowledge/add \
  -H "Content-Type: application/json" \
  -d '{"summary": "测试总结", "content": "测试内容", "user_id": "test"}'

# 删除知识库
curl -X POST http://localhost:8889/api/v1/knowledge/delete \
  -H "Content-Type: application/json" \
  -d '{"vector_id": "vector_id_here", "user_id": "test"}'
```

## 依赖服务

- **bll-context**: RPC服务，提供向量化功能
- **bs-rag**: RPC服务，提供向量数据库操作

## 技术栈

- **框架**: Go-Zero
- **RPC**: gRPC
- **配置**: YAML
- **日志**: Go-Zero Logger
