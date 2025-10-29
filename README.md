# JXZY AI应用服务

基于Go-Zero微服务框架构建的AI应用服务，采用三层微服务架构。

## 📋 项目概述

JXZY是一个现代化的AI应用微服务平台，致力于提供高性能、可扩展的AI服务能力。项目采用三层微服务架构，提供智能对话、会话管理等功能。

## 🏗️ 技术架构

### 核心设计原则
- **分层解耦**: 三层微服务架构，职责清晰
- **高性能**: 支持高并发和低延迟响应
- **可扩展**: 模块化设计，易于水平扩展

### 技术栈
- **框架**: Go-Zero微服务框架
- **语言**: Go 1.19+
- **数据库**: MySQL 8.0+
- **服务发现**: etcd

## 🚀 快速开始

### 环境要求
- Go 1.19+
- MySQL 8.0+
- etcd

### 启动服务
```bash
# 克隆项目
git clone https://github.com/your-org/jxzy.git
cd jxzy

# 测试日志配置
./scripts/test_logging.sh

# 启动各个服务
cd apis/api_chat && ./scripts/compile.sh && ./scripts/restart.sh
cd ../../bll/bll_context && ./scripts/compile.sh && ./scripts/restart.sh
cd ../../bs/bs_llm && ./scripts/compile.sh && ./scripts/restart.sh

# 查看统一日志
./scripts/view_logs.sh
```

## 📊 项目结构

```
jxzy/
├── apis/api_chat/           # Gateway API层
│   ├── scripts/            # 编译和重启脚本
│   └── README.md           # API服务说明
├── bll/bll_context/        # Business Logic Layer
│   ├── scripts/            # 编译和重启脚本
│   └── README.md           # BLLS服务说明
├── bs/bs_llm/              # Basic Service Layer
│   ├── scripts/            # 编译和重启脚本
│   └── README.md           # BS服务说明
├── common/                 # 公共组件
│   └── logger/             # 统一日志系统
├── docs/                   # 文档目录
│   └── logging.md          # 日志配置说明
├── logs/                   # 统一日志目录
├── scripts/                # 项目脚本
│   ├── test_logging.sh     # 日志配置测试
│   └── view_logs.sh        # 日志查看工具
└── README.md               # 本文件
```

## 📝 开发指南

### 日志管理
```bash
# 测试日志配置
./scripts/test_logging.sh

# 查看统一日志
./scripts/view_logs.sh

# 查看日志配置说明
cat docs/logging.md
```

### 修改API定义
```bash
cd apis/api_chat
# 编辑 chat.api 文件
./scripts/compile.sh
./scripts/restart.sh
```

### 修改Proto定义
```bash
cd bll/bll_context
# 编辑 bllcontext.proto 文件
./scripts/compile.sh
./scripts/restart.sh
```

### 修改BS服务
```bash
cd bs/bs_llm
# 编辑 bsllm.proto 文件
./scripts/compile.sh
./scripts/restart.sh
```

## 🤝 贡献指南

我们欢迎社区贡献！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系我们

- 项目首页: [GitHub Repository](https://github.com/your-org/jxzy)
- 问题反馈: [Issues](https://github.com/your-org/jxzy/issues)
- 讨论交流: [Discussions](https://github.com/your-org/jxzy/discussions)

