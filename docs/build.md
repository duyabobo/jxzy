# JXZY 编译说明

## 目录结构

```
JXZY/
├── bin/                    # 统一的可执行文件目录
│   ├── chat-api           # API Gateway 可执行文件
│   ├── bll-context        # BLL Context 服务可执行文件
│   └── bs-llm             # BS LLM 服务可执行文件
├── apis/
│   └── api_chat/          # API Gateway 项目
├── bll/
│   └── bll_context/       # BLL Context 项目
├── bs/
│   └── bs_llm/            # BS LLM 项目
├── logs/                  # 统一日志目录
├── scripts/
│   ├── build_all.sh       # 统一编译脚本
│   └── view_logs.sh       # 日志查看脚本
└── docs/                  # 文档目录
```

## 编译方式

### 1. 统一编译所有服务

```bash
# 在项目根目录执行
./scripts/build_all.sh
```

这个脚本会：
- 自动创建 `bin/` 目录
- 按顺序编译所有服务
- 将所有可执行文件放到 `bin/` 目录下

### 2. 单独编译服务

```bash
# 编译 BS LLM 服务
cd bs/bs_llm && ./scripts/compile.sh

# 编译 BLL Context 服务
cd bll/bll_context && ./scripts/compile.sh

# 编译 API Gateway 服务
cd apis/api_chat && ./scripts/compile.sh
```

## 可执行文件

编译完成后，所有可执行文件都会放在 `bin/` 目录下：

| 服务 | 可执行文件 | 说明 |
|------|------------|------|
| **API Gateway** | `bin/chat-api` | HTTP API 入口，监听端口 8888 |
| **BLL Context** | `bin/bll-context` | 业务逻辑层 RPC 服务，监听端口 8080 |
| **BS LLM** | `bin/bs-llm` | LLM 服务 RPC，监听端口 8081 |

## 启动服务

编译完成后，使用以下命令启动服务：

```bash
# 启动 BS LLM 服务
cd bs/bs_llm && ./scripts/restart.sh

# 启动 BLL Context 服务
cd bll/bll_context && ./scripts/restart.sh

# 启动 API Gateway 服务
cd apis/api_chat && ./scripts/restart.sh
```

## 重启脚本说明

重启脚本会：
1. 根据监听端口停止现有服务
2. 删除旧的可执行文件
3. 重新编译项目
4. 启动新服务
5. 验证服务状态

## 日志管理

所有服务的日志都统一输出到 `logs/` 目录：

```bash
# 查看所有日志
./scripts/view_logs.sh

# 查看特定服务日志
tail -f logs/chat-api.log
tail -f logs/bll-context.log
tail -f logs/bs-llm.log
```

## 开发流程

### 1. 首次开发

```bash
# 1. 统一编译所有服务
./scripts/build_all.sh

# 2. 启动所有服务
cd bs/bs_llm && ./scripts/restart.sh
cd bll/bll_context && ./scripts/restart.sh
cd apis/api_chat && ./scripts/restart.sh
```

### 2. 日常开发

```bash
# 修改代码后，重启对应服务即可
cd apis/api_chat && ./scripts/restart.sh
```

### 3. 查看日志

```bash
# 使用统一日志查看工具
./scripts/view_logs.sh
```

## 注意事项

1. **可执行文件位置**: 所有可执行文件都统一放在 `bin/` 目录下
2. **日志位置**: 所有日志都统一放在 `logs/` 目录下
3. **配置文件**: 各服务的配置文件仍在各自的 `etc/` 目录下
4. **端口冲突**: 确保端口 8080、8081、8888 没有被其他程序占用

## 故障排除

### 编译失败

1. 检查 Go 环境是否正确安装
2. 检查 goctl 工具是否正确安装
3. 检查 protoc 工具是否正确安装（RPC 服务需要）

### 启动失败

1. 检查端口是否被占用
2. 检查配置文件是否正确
3. 查看日志文件排查错误

### 服务无法连接

1. 检查各服务的端口监听状态
2. 检查 RPC 服务配置
3. 查看日志文件排查网络问题
