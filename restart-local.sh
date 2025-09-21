#!/bin/bash

# Easy-Chat 重启脚本

echo "重启 Easy-Chat 服务..."

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 先停止所有服务
echo "停止现有服务..."
"$SCRIPT_DIR/stop-local.sh"

# 等待一会儿确保服务完全停止
sleep 3

# 重新启动服务
echo "重新启动服务..."
"$SCRIPT_DIR/start-local.sh"
