#!/bin/bash

# api_knowledge 编译脚本
# 用于编译 api_knowledge 服务

set -e

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "=========================================="
echo "开始编译 api_knowledge 服务"
echo "项目目录: $PROJECT_DIR"
echo "=========================================="

# 切换到项目目录
cd "$PROJECT_DIR"

# 检查go.mod文件是否存在
if [ ! -f "go.mod" ]; then
    echo "错误: 在 $PROJECT_DIR 中未找到 go.mod 文件"
    exit 1
fi

# 清理之前的构建
echo "清理之前的构建文件..."
rm -f ../../bin/api-knowledge-server

# 下载依赖
echo "下载依赖..."
go mod tidy

# 编译项目
echo "编译 api_knowledge 服务..."
go build -o ../../bin/api-knowledge-server .

# 检查编译结果
if [ -f "../../bin/api-knowledge-server" ]; then
    echo "=========================================="
    echo "✅ api_knowledge 服务编译成功!"
    echo "可执行文件: ../../bin/api-knowledge-server"
    echo "=========================================="
    
    # 显示文件信息
    ls -lh ../../bin/api-knowledge-server
else
    echo "=========================================="
    echo "❌ api_knowledge 服务编译失败!"
    echo "=========================================="
    exit 1
fi
