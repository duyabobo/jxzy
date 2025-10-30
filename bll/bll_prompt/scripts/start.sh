#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

CONFIG_FILE="etc/bllprompt.yaml"
PORT=8005

echo "🔧 启动 bll-prompt (port: $PORT)"

if [ ! -f "$CONFIG_FILE" ]; then
  echo "❌ 配置文件 $CONFIG_FILE 不存在"
  exit 1
fi

rm -f ../../bin/bll-prompt || true
bash scripts/compile.sh

mkdir -p ../../logs
nohup ../../bin/bll-prompt -f "$CONFIG_FILE" >> ../../logs/access.log 2>&1 &
PID=$!
sleep 3

if kill -0 "$PID" 2>/dev/null; then
  echo "✅ bll-prompt 已启动 (PID: $PID)"
else
  echo "❌ bll-prompt 启动失败"
  exit 1
fi

