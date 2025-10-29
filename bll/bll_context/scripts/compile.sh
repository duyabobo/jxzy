#!/bin/bash

# JXZY BLLS Context 编译脚本
# 用于重新生成proto代码并编译项目

set -e

echo "🚀 开始编译 BLLS Context 项目..."

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

echo "📁 工作目录: $(pwd)"

# 检查protoc是否安装
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc 未安装，请先安装 protobuf"
    echo "   macOS: brew install protobuf"
    echo "   Ubuntu: sudo apt-get install protobuf-compiler"
    exit 1
fi

# 检查goctl是否安装
if ! command -v goctl &> /dev/null; then
    echo "❌ goctl 未安装，请先安装 goctl"
    echo "   安装命令: go install github.com/zeromicro/go-zero/tools/goctl@latest"
    exit 1
fi

# 检查proto文件是否存在
if [ ! -f "bllcontext.proto" ]; then
    echo "❌ bllcontext.proto 文件不存在"
    exit 1
fi

echo "🔧 步骤1: 使用 goctl 重新生成proto代码..."
# 删除旧的生成文件
rm -rf bll_context/*.pb.go
rm -rf bll_context/*_grpc.pb.go

# 重新生成proto代码
goctl rpc protoc bllcontext.proto --go_out=. --go-grpc_out=. --zrpc_out=.

if [ $? -eq 0 ]; then
    echo "✅ Proto代码生成成功"
else
    echo "❌ Proto代码生成失败"
    exit 1
fi

echo "🔧 步骤2: 编译项目..."
# 创建项目根目录的bin目录
mkdir -p ../../bin

# 编译项目到项目根目录的bin目录
go build -o ../../bin/bll-context .

if [ $? -eq 0 ]; then
    echo "✅ 编译成功"
    echo "📦 可执行文件: ../../bin/bll-context"
else
    echo "❌ 编译失败"
    exit 1
fi

echo "🎉 BLLS Context 编译完成！"
echo "💡 提示: 使用 ./scripts/restart.sh 重启服务"
