# JXZY Chat API 服务

基于 go-zero 框架的 AI 聊天服务 API，提供流式聊天功能。

## 📁 项目结构

```
apis/
├── chat.api                    # API 定义文件
├── chat.go                     # 主程序入口
├── etc/
│   └── chat-api.yaml          # 配置文件
├── internal/
│   ├── config/                # 配置相关
│   ├── handler/               # HTTP 处理器
│   ├── logic/                 # 业务逻辑
│   ├── middleware/            # 中间件
│   ├── svc/                   # 服务上下文
│   └── types/                 # 类型定义
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

2. **Go**: 1.16 或更高版本

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
# 重新生成代码
goctl api go -api chat.api -dir .

# 编译项目
go build -o chat-api .

# 启动服务
./chat-api -f etc/chat-api.yaml
```

## 📋 API 接口

### 流式聊天接口

- **URL**: `POST /api/v1/chat/stream`
- **Content-Type**: `application/json`
- **Response-Type**: `text/event-stream` (SSE)

#### 请求参数
```json
{
  "message": "你好，请介绍一下你自己",
  "session_id": "session_20231201120000",  // 可选，为空时自动创建
  "scene_code": "general"                  // 必填，业务场景编码
}
```

#### 响应格式
```json
{
  "session_id": "session_20231201120000",
  "scene_code": "general",
  "delta": "Hello! I received your message: 你好，请介绍一下你自己",
  "finished": false
}
```

## 🔧 配置说明

配置文件：`etc/chat-api.yaml`

```yaml
Name: chat-api
Host: 0.0.0.0
Port: 8888

BllContextRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: bll_context.rpc
  NonBlock: true
```

## 📝 开发指南

### 修改 API 定义

1. 编辑 `chat.api` 文件
2. 运行 `./scripts/compile.sh` 重新生成代码
3. 运行 `./scripts/restart.sh` 重启服务

### 添加新的处理器

1. 在 `chat.api` 中定义新的接口
2. 运行 `./scripts/compile.sh` 生成代码框架
3. 在 `internal/logic/` 中实现业务逻辑
4. 重启服务

## 🚨 注意事项

1. **服务依赖**: 需要 bll_context RPC 服务
2. **端口配置**: 默认使用 8888 端口
3. **日志文件**: 服务日志保存在 `chat-api.log`
4. **CORS**: 已配置跨域支持

## 🔧 故障排除

### 编译失败
- 检查 goctl 是否正确安装
- 检查 API 文件语法是否正确

### 启动失败
- 检查配置文件是否存在
- 检查端口是否被占用
- 查看日志文件获取详细错误信息

### 连接失败
- 检查 etcd 服务是否运行
- 检查 bll_context RPC 服务是否可用

## 📞 服务信息

- **服务地址**: http://localhost:8888
- **日志文件**: `chat-api.log`
- **配置文件**: `etc/chat-api.yaml`
