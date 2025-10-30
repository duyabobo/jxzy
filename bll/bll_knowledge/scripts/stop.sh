#!/bin/bash

set -e

PORT=8006
echo "🛑 停止 bll-knowledge (port: $PORT)"

PIDS=$(lsof -ti:$PORT 2>/dev/null || true)
if [ -n "$PIDS" ]; then
  for PID in $PIDS; do
    kill -TERM "$PID" 2>/dev/null || true
  done
  sleep 2
  PIDS=$(lsof -ti:$PORT 2>/dev/null || true)
  if [ -n "$PIDS" ]; then
    for PID in $PIDS; do
      kill -KILL "$PID" 2>/dev/null || true
    done
  fi
  echo "✅ 已停止"
else
  echo "ℹ️ 未发现运行中的进程"
fi

