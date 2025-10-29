#!/bin/bash

# JXZY BS RAG 启动脚本
# 用于启动 BS RAG 服务

set -e

echo "🚀 开始启动 BS RAG 服务..."

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

echo "📁 工作目录: $(pwd)"

# 配置文件路径
CONFIG_FILE="etc/bsrag.yaml"
PORT=8082

echo "=========================================="
echo "开始启动 BS RAG 服务"
echo "项目目录: $(pwd)"
echo "配置文件: $CONFIG_FILE"
echo "端口: $PORT"
echo "=========================================="

# 检查配置文件是否存在
if [ ! -f "etc/bsrag.yaml" ]; then
    echo "❌ 配置文件 etc/bsrag.yaml 不存在"
    exit 1
fi

echo "🔧 步骤1: 删除旧的可执行文件..."
# 获取项目根目录
PROJECT_ROOT="$(cd ../../.. && pwd)"
BIN_DIR="$PROJECT_ROOT/bin"

# 删除旧的可执行文件
if [ -f "$BIN_DIR/bs_rag" ]; then
    rm -f "$BIN_DIR/bs_rag"
    echo "🗑️  已删除旧的可执行文件"
fi

echo "🔧 步骤2: 重新编译项目..."
# 运行编译脚本
if [ -f "scripts/compile.sh" ]; then
    bash scripts/compile.sh
    if [ $? -ne 0 ]; then
        echo "❌ 编译失败"
        exit 1
    fi
else
    echo "❌ 编译脚本 scripts/compile.sh 不存在"
    exit 1
fi

echo "🔧 步骤3: 启动服务..."
# 检查编译后的可执行文件是否存在
if [ ! -f "$BIN_DIR/bs_rag" ]; then
    echo "❌ 编译后的可执行文件 $BIN_DIR/bs_rag 不存在"
    exit 1
fi

# 启动服务
echo "启动 BS RAG 服务..."
# 确保logs目录存在
mkdir -p ../../logs
nohup "$BIN_DIR/bs_rag" -f etc/bsrag.yaml >> ../../logs/access.log 2>&1 &
SERVICE_PID=$!

# 等待服务启动
sleep 5

# 检查服务是否成功启动
if kill -0 "$SERVICE_PID" 2>/dev/null; then
    echo "✅ 服务启动成功 (PID: $SERVICE_PID)"
    echo "📝 日志文件: ../../logs/access.log"
    echo "📋 查看日志: tail -f ../../logs/access.log"
    
    # 检查端口是否正常监听
    if lsof -i:$PORT >/dev/null 2>&1; then
        echo "✅ 端口 $PORT 监听正常"
    else
        echo "⚠️  端口 $PORT 监听异常，请检查日志"
    fi
else
    echo "❌ 服务启动失败"
    echo "📋 查看错误日志: tail -f ../../logs/access.log"
    exit 1
fi

echo "🎉 BS RAG 服务启动完成！"
