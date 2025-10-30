#!/bin/bash

# JXZY BLL Knowledge 编译脚本

set -e

echo "🚀 开始编译 BLL Knowledge 项目..."

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

echo "📁 工作目录: $(pwd)"

if ! command -v protoc &> /dev/null; then
    echo "❌ protoc 未安装"
    exit 1
fi

if ! command -v goctl &> /dev/null; then
    echo "❌ goctl 未安装"
    exit 1
fi

if [ ! -f "bllknowledge.proto" ]; then
    echo "❌ bllknowledge.proto 文件不存在"
    exit 1
fi

echo "🔧 生成 proto 代码..."
rm -rf bll_knowledge/*.pb.go bll_knowledge/*_grpc.pb.go
goctl rpc protoc bllknowledge.proto --go_out=. --go-grpc_out=. --zrpc_out=.

echo "🔧 编译可执行文件..."
mkdir -p ../../bin
go build -o ../../bin/bll-knowledge .

echo "✅ 编译成功: ../../bin/bll-knowledge"

