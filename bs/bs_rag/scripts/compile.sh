#!/bin/bash

# 编译 bs_rag 服务
echo "Compiling bs_rag service..."

# 设置工作目录
cd "$(dirname "$0")/.."

# 获取项目根目录
PROJECT_ROOT="$(cd ../../.. && pwd)"
BIN_DIR="$PROJECT_ROOT/bin"

# 确保 bin 目录存在
mkdir -p "$BIN_DIR"

# 编译
go build -o "$BIN_DIR/bs_rag" bsrag.go

if [ $? -eq 0 ]; then
    echo "✅ bs_rag service compiled successfully!"
    echo "Binary location: $BIN_DIR/bs_rag"
else
    echo "❌ Compilation failed!"
    exit 1
fi
