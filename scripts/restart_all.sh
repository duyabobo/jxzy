#!/bin/bash

# JXZY 统一服务管理脚本
# 直接执行：停止所有服务，然后按顺序启动

set -e

echo "=========================================="
echo "    JXZY 服务管理脚本"
echo "=========================================="

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

echo "📁 项目根目录: $PROJECT_ROOT"

# 1. 停止所有服务
echo ""
echo "🛑 步骤1: 停止所有服务..."
pkill -f "chat-api|api-knowledge|bll-context|bs-llm|bs_rag" || true
sleep 3

# 2. 启动 bs-llm 服务
echo ""
echo "🚀 步骤2: 启动 bs-llm 服务..."
cd bs/bs_llm
bash scripts/restart.sh
cd "$PROJECT_ROOT"
sleep 2

# 3. 启动 bs-rag 服务
echo ""
echo "🚀 步骤3: 启动 bs-rag 服务..."
cd bs/bs_rag
bash scripts/restart.sh
cd "$PROJECT_ROOT"
sleep 2

# 4. 启动 bll-context 服务
echo ""
echo "🚀 步骤4: 启动 bll-context 服务..."
cd bll/bll_context
bash scripts/restart.sh
cd "$PROJECT_ROOT"
sleep 2

# 5. 启动 chat-api 服务
echo ""
echo "🚀 步骤5: 启动 chat-api 服务..."
cd apis/api_chat
bash scripts/restart.sh
cd "$PROJECT_ROOT"
sleep 2

# 6. 启动 api-knowledge 服务
echo ""
echo "🚀 步骤6: 启动 api-knowledge 服务..."
cd apis/api_knowledge
bash scripts/restart.sh
cd "$PROJECT_ROOT"

echo ""
echo "=========================================="
echo "✅ 所有服务启动完成！"
echo "=========================================="
echo ""
echo "服务状态:"
echo "  bs-llm        (端口: 8081) - 运行中"
echo "  bs-rag        (端口: 8082) - 运行中"
echo "  bll-context   (端口: 8080) - 运行中"
echo "  chat-api      (端口: 8888) - 运行中"
echo "  api-knowledge (端口: 8889) - 运行中"
echo ""
echo "日志文件位置:"
echo "  所有服务日志统一存储在项目根目录的 logs/ 目录下"
echo "  可执行文件统一存储在项目根目录的 bin/ 目录下"
echo ""
echo "API测试:"
echo "  # 聊天API测试:"
echo "  curl -X POST http://localhost:8888/api/v1/chat/stream \\"
echo "    -H \"Content-Type: application/json\" \\"
echo "    -d '{\"user_id\": \"test\", \"message\": \"你好\", \"scene_code\": \"chat_general\"}'"
echo ""
echo "  # 知识库API测试:"
echo "  curl -X POST http://localhost:8889/api/v1/knowledge/add \\"
echo "    -H \"Content-Type: application/json\" \\"
echo "    -d '{\"summary\": \"测试总结\", \"content\": \"测试内容\", \"user_id\": \"test\"}'"
echo ""
echo "  curl -X POST http://localhost:8889/api/v1/knowledge/delete \\"
echo "    -H \"Content-Type: application/json\" \\"
echo "    -d '{\"vector_id\": \"test_id\", \"user_id\": \"test\"}'"
echo ""
