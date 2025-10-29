#!/bin/bash

# JXZY BLLS Context 停止脚本
# 用于停止 BLLS Context 服务

set -e

echo "🛑 开始停止 BLLS Context 服务..."

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

echo "📁 工作目录: $(pwd)"

# 配置文件路径
CONFIG_FILE="etc/bllcontext.yaml"
PORT=8080

echo "=========================================="
echo "开始停止 BLLS Context 服务"
echo "项目目录: $(pwd)"
echo "配置文件: $CONFIG_FILE"
echo "端口: $PORT"
echo "=========================================="

echo "🔧 步骤1: 根据监听端口停止现有服务..."
# 根据监听端口8080查找并停止进程
PIDS=$(lsof -ti:$PORT 2>/dev/null || true)
if [ -n "$PIDS" ]; then
    echo "🛑 找到监听端口 $PORT 的进程: $PIDS"
    for PID in $PIDS; do
        echo "   正在停止进程 $PID..."
        kill -TERM "$PID" 2>/dev/null || true
    done
    
    # 等待进程结束
    sleep 3
    
    # 强制杀死仍在运行的进程
    PIDS=$(lsof -ti:$PORT 2>/dev/null || true)
    if [ -n "$PIDS" ]; then
        echo "⚠️  强制停止进程..."
        for PID in $PIDS; do
            kill -KILL "$PID" 2>/dev/null || true
        done
    fi
    echo "✅ 端口 $PORT 上的服务已停止"
else
    echo "ℹ️  端口 $PORT 上没有找到运行中的服务"
fi

echo "等待端口 $PORT 释放..."
sleep 2

# 检查端口是否已释放
if lsof -i:$PORT >/dev/null 2>&1; then
    echo "❌ 端口 $PORT 仍被占用"
    exit 1
else
    echo "✅ 端口 $PORT 已释放"
fi

echo "🎉 BLLS Context 服务停止完成！"
