#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

CONFIG_FILE="etc/bllknowledge.yaml"
PORT=8006

echo "🔧 启动 bll-knowledge (port: $PORT)"

if [ ! -f "$CONFIG_FILE" ]; then
  echo "❌ 配置文件 $CONFIG_FILE 不存在"
  exit 1
fi

rm -f ../../bin/bll-knowledge || true
bash scripts/compile.sh

mkdir -p ../../logs
nohup ../../bin/bll-knowledge -f "$CONFIG_FILE" >> ../../logs/access.log 2>&1 &
PID=$!
sleep 3

if kill -0 "$PID" 2>/dev/null; then
  echo "✅ bll-knowledge 已启动 (PID: $PID)"
else
  echo "❌ bll-knowledge 启动失败"
  exit 1
fi

