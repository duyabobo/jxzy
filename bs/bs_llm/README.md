# JXZY BS LLM 服务

基础服务层 LLM 服务，负责与各种大语言模型进行交互。

## 📁 项目结构

```
bs_llm/
├── bsllm.proto                 # Proto 定义文件
├── bs_llm/                     # 生成的 proto 代码
├── bsllmservice/               # 服务实现
├── etc/
│   └── bsllm.yaml             # 配置文件
├── internal/
│   ├── common/                # 公共工具
│   ├── config/                # 配置相关
│   ├── logic/                 # 业务逻辑
│   ├── model/                 # 数据模型
│   ├── provider/              # LLM 提供商
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
goctl rpc protoc bsllm.proto --go_out=. --go-grpc_out=. --zrpc_out=.

# 编译项目
go build -o bs_llm .

# 启动服务
./bs_llm -f etc/bsllm.yaml
```

## 📋 服务接口

### StreamLLM RPC 接口

- **服务名**: `BsLlmService`
- **方法**: `StreamLLM`
- **协议**: gRPC

#### 请求参数
```protobuf
message StreamLLMRequest {
    string message = 1;
    string model_id = 2;
    float temperature = 3;
    int64 max_tokens = 4;
    string scene_code = 5;
}
```

#### 响应格式
```protobuf
message StreamLLMResponse {
    string delta = 1;
    bool finished = 2;
    TokenUsage usage = 3;
}
```

## 🔧 配置说明

配置文件：`etc/bsllm.yaml`

```yaml
Name: bs_llm
Host: 0.0.0.0
Port: 8081

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: bs_llm.rpc

DataSource:
  Host: 127.0.0.1
  Port: 3306
  User: root
  Password: password
  Database: jxzy

LLM:
  Doubao:
    ApiKey: your_api_key_here
    BaseURL: https://api.doubao.com
```

## 📝 开发指南

### 修改 Proto 定义

1. 编辑 `bsllm.proto` 文件
2. 运行 `./scripts/compile.sh` 重新生成代码
3. 运行 `./scripts/restart.sh` 重启服务

### 添加新的 LLM 提供商

1. 在 `internal/provider/` 目录下创建新的提供商实现
2. 实现 `Provider` 接口
3. 在配置文件中添加相应的配置
4. 重启服务

### 添加新的 RPC 方法

1. 在 `bsllm.proto` 中定义新的方法
2. 运行 `./scripts/compile.sh` 生成代码框架
3. 在 `internal/logic/` 中实现业务逻辑
4. 重启服务

## 🚨 注意事项

1. **服务依赖**: 需要 etcd 和 MySQL 数据库
2. **端口配置**: 默认使用 8081 端口
3. **日志文件**: 服务日志保存在 `bs_llm.log`
4. **API密钥**: 需要配置相应的 LLM 提供商 API 密钥
5. **数据库**: 需要初始化数据库表结构

## 🔧 故障排除

### 编译失败
- 检查 goctl 和 protoc 是否正确安装
- 检查 proto 文件语法是否正确

### 启动失败
- 检查配置文件是否存在
- 检查端口是否被占用
- 检查数据库连接是否正常
- 检查 API 密钥配置是否正确
- 查看日志文件获取详细错误信息

### 连接失败
- 检查 etcd 服务是否运行
- 检查数据库服务是否可用
- 检查 LLM 提供商 API 是否可访问

## 📞 服务信息

- **服务类型**: gRPC RPC 服务
- **注册地址**: etcd://127.0.0.1:2379/bs_llm.rpc
- **日志文件**: `bs_llm.log`
- **配置文件**: `etc/bsllm.yaml`
