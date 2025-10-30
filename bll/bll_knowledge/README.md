# JXZY BLL Knowledge 服务

业务逻辑层服务，负责知识库与记忆模块的管理，底层调用 `BsRAG-RPC` 执行向量化与检索。

## 📁 项目结构

```
bll_knowledge/
├── bllknowledge.proto          # Proto 定义文件
├── bll_knowledge/              # 生成的 proto 代码
├── bllknowledgeservice/        # 服务端注册（自动生成）
├── etc/
│   └── bllknowledge.yaml       # 配置文件（端口 8006）
├── internal/
│   ├── config/                 # 配置结构体
│   ├── logic/                  # 业务逻辑（Add/Delete Vector Knowledge）
│   ├── model/                  # 数据模型（knowledge_* ORM）
│   ├── server/                 # gRPC Server（自动生成）
│   └── svc/                    # 服务上下文（RAG 客户端、DB 等）
└── scripts/                    # 脚本目录
    ├── compile.sh
    ├── restart.sh
    ├── start.sh
    └── stop.sh
```

## 🚀 快速开始

```bash
cd bll/bll_knowledge
./scripts/restart.sh
```

## 🔧 配置说明（etc/bllknowledge.yaml）

- ListenOn: 0.0.0.0:8006
- MySQL: 元数据（知识文件/语义段）
- BsRagRpc: 直连 `127.0.0.1:8082`

## 📋 服务接口

- AddVectorKnowledge
- DeleteVectorKnowledge

详见 `bllknowledge.proto`。


