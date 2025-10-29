#!/bin/bash

# JXZY Chat API 重启脚本
# 用于重启 Chat API 服务

set -e

echo "🔄 开始重启 Chat API 服务..."

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

echo "📁 工作目录: $(pwd)"

echo "=========================================="
echo "开始重启 Chat API 服务"
echo "项目目录: $(pwd)"
echo "=========================================="

# 停止服务
echo "🔧 步骤1: 停止服务..."
bash scripts/stop.sh

# 启动服务
echo "🔧 步骤2: 启动服务..."
bash scripts/start.sh

echo "🎉 Chat API 服务重启完成！"
