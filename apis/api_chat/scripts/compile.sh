#!/bin/bash

# JXZY Chat API 编译脚本
# 用于重新生成代码并编译项目

set -e

echo "🚀 开始编译 Chat API 项目..."

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

echo "📁 工作目录: $(pwd)"

# 检查goctl是否安装
if ! command -v goctl &> /dev/null; then
    echo "❌ goctl 未安装，请先安装 goctl"
    echo "   安装命令: go install github.com/zeromicro/go-zero/tools/goctl@latest"
    exit 1
fi

# 检查API文件是否存在
if [ ! -f "chat.api" ]; then
    echo "❌ chat.api 文件不存在"
    exit 1
fi

echo "🔧 步骤1: 使用 goctl 生成代码..."
# 生成代码（goctl默认不会覆盖已存在的文件）
goctl api go -api chat.api -dir .

if [ $? -eq 0 ]; then
    echo "✅ 代码生成成功"
else
    echo "❌ 代码生成失败"
    exit 1
fi

echo "�� 步骤2: 编译项目..."
# 创建项目根目录的bin目录
mkdir -p ../../bin

# 编译项目到项目根目录的bin目录
go build -o ../../bin/chat-api .

if [ $? -eq 0 ]; then
    echo "✅ 编译成功"
    echo "📦 可执行文件: ../../bin/chat-api"
else
    echo "❌ 编译失败"
    exit 1
fi

echo "🎉 Chat API 编译完成！"
echo "💡 提示: 使用 ./scripts/restart.sh 重启服务"
