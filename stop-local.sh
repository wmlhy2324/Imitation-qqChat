#!/bin/bash

# Easy-Chat 停止脚本

echo "正在停止 Easy-Chat 服务..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 停止函数
stop_service() {
    local port=$1
    local service_name=$2
    
    echo -e "${YELLOW}停止 $service_name (端口: $port)...${NC}"
    
    # 查找并杀死占用端口的进程
    local pids=$(lsof -ti :$port 2>/dev/null)
    
    if [ -n "$pids" ]; then
        echo "$pids" | xargs kill 2>/dev/null || true
        sleep 2
        
        # 检查是否还有进程
        local remaining_pids=$(lsof -ti :$port 2>/dev/null)
        if [ -n "$remaining_pids" ]; then
            echo "$remaining_pids" | xargs kill -9 2>/dev/null || true
        fi
        
        echo -e "${GREEN}✓ $service_name 已停止${NC}"
    else
        echo -e "${YELLOW}$service_name 未运行${NC}"
    fi
}

# 按相反顺序停止服务
echo "按服务启动的相反顺序停止..."

stop_service "10091" "Task MQ"
stop_service "10090" "IM WebSocket"
stop_service "8882" "IM API"
stop_service "8881" "Social API"  
stop_service "8888" "User API"
stop_service "10002" "IM RPC"
stop_service "10001" "Social RPC"
stop_service "10010" "User RPC"

# 额外清理：杀死所有相关的 go run 进程
echo -e "${YELLOW}清理相关进程...${NC}"
pkill -f "go run.*\.go" 2>/dev/null || true

echo -e "${GREEN}所有服务已停止！${NC}"
