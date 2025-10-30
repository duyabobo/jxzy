#!/bin/bash

set -e

bash "$(dirname "$0")/stop.sh"
bash "$(dirname "$0")/start.sh"

echo "🎉 bll-prompt 重启完成"

