# JXZY BLL Prompt 服务

业务逻辑层服务，负责提示词工程（模板管理、动态生成、A/B 测试、合规校验），为 `BllsContext-RPC` 提供提示词能力。

## 📁 项目结构

```
bll_prompt/
├── bllprompt.proto            # Proto 定义文件
├── bll_prompt/                # 生成的 proto 代码
├── bllpromptservice/          # 服务端注册（自动生成）
├── etc/
│   └── bllprompt.yaml         # 配置文件（端口 8005）
├── internal/
│   ├── config/                # 配置结构体
│   ├── logic/                 # 业务逻辑
│   ├── server/                # gRPC Server（自动生成）
│   └── svc/                   # 服务上下文
└── scripts/
    ├── compile.sh
    ├── restart.sh
    ├── start.sh
    └── stop.sh
```

## 🚀 快速开始

```bash
cd bll/bll_prompt
./scripts/restart.sh
```

## 🔧 配置说明（etc/bllprompt.yaml）

- ListenOn: 0.0.0.0:8005
- MySQL: 模板与版本管理（预留）

## 📋 服务接口

暂为占位，后续扩展模板管理与动态生成等接口。


