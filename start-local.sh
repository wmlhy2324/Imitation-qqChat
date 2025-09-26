#!/bin/bash

# Easy-Chat 本地启动脚本
# 按顺序启动：User RPC -> Social RPC -> IM RPC -> User API -> Social API -> IM API

set -e

echo "========================================="
echo "        Easy-Chat 本地启动脚本"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 日志目录
LOG_DIR="$SCRIPT_DIR/logs"
mkdir -p "$LOG_DIR"

# 清理函数
cleanup() {
    echo -e "\n${RED}正在停止所有服务...${NC}"
    
    # 获取所有子进程的PID
    jobs -p | xargs -r kill 2>/dev/null || true
    
    # 等待进程结束
    sleep 2
    
    # 强制杀死残留进程
    jobs -p | xargs -r kill -9 2>/dev/null || true
    
    echo -e "${GREEN}所有服务已停止${NC}"
    exit 0
}

# 注册信号处理
trap cleanup SIGINT SIGTERM

# 检查端口是否被占用
check_port() {
    local port=$1
    local service=$2
    
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        echo -e "${YELLOW}警告: 端口 $port ($service) 已被占用${NC}"
        return 1
    fi
    return 0
}

# 启动服务函数
start_service() {
    local service_name=$1
    local service_dir=$2
    local service_file=$3
    local log_file=$4
    local port=$5
    local wait_time=${6:-5}
    
    echo -e "${BLUE}启动 $service_name...${NC}"
    
    # 检查端口
    if ! check_port "$port" "$service_name"; then
        echo -e "${RED}端口 $port 被占用，请先停止占用该端口的进程${NC}"
        exit 1
    fi
    
    # 进入服务目录
    cd "$SCRIPT_DIR/$service_dir"
    
    # 启动服务
    go run "$service_file" > "$LOG_DIR/$log_file" 2>&1 &
    local pid=$!
    
    echo -e "${YELLOW}$service_name 启动中... (PID: $pid)${NC}"
    
    # 等待服务启动
    local count=0
    while [ $count -lt $wait_time ]; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            echo -e "${GREEN}✓ $service_name 启动成功 (端口: $port)${NC}"
            return 0
        fi
        sleep 1
        count=$((count + 1))
        echo -n "."
    done
    
    echo -e "\n${RED}✗ $service_name 启动失败或超时${NC}"
    echo -e "${YELLOW}请检查日志文件: $LOG_DIR/$log_file${NC}"
    return 1
}

# 检查依赖服务
echo -e "${BLUE}检查依赖服务...${NC}"

# 检查 MySQL
if ! nc -z 127.0.0.1 13306 2>/dev/null; then
    echo -e "${RED}错误: MySQL 服务 (127.0.0.1:13306) 未启动${NC}"
    exit 1
fi
echo -e "${GREEN}✓ MySQL 服务正常${NC}"

# 检查 Redis
if ! nc -z 127.0.0.1 16379 2>/dev/null; then
    echo -e "${RED}错误: Redis 服务 (127.0.0.1:16379) 未启动${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Redis 服务正常${NC}"

# 检查 ETCD
if ! nc -z 127.0.0.1 3379 2>/dev/null; then
    echo -e "${RED}错误: ETCD 服务 (127.0.0.1:3379) 未启动${NC}"
    exit 1
fi
echo -e "${GREEN}✓ ETCD 服务正常${NC}"

# 检查 MongoDB (IM 服务需要)
if ! nc -z 127.0.0.1 47017 2>/dev/null; then
    echo -e "${YELLOW}警告: MongoDB 服务 (127.0.0.1:47017) 未启动，IM 服务可能无法正常工作${NC}"
fi

# 检查 Kafka (消息队列服务需要)
if ! nc -z 127.0.0.1 9092 2>/dev/null; then
    echo -e "${YELLOW}警告: Kafka 服务 (127.0.0.1:9092) 未启动，消息队列可能无法正常工作${NC}"
fi

echo ""
echo -e "${GREEN}========================================="
echo -e "         开始启动 RPC 服务"
echo -e "=========================================${NC}"

# 1. 启动 User RPC 服务
start_service "User RPC" "apps/user/rpc" "user.go" "user-rpc.log" "10010" 8

# 2. 启动 Social RPC 服务
start_service "Social RPC" "apps/social/rpc" "social.go" "social-rpc.log" "10001" 8

# 3. 启动 IM RPC 服务
start_service "IM RPC" "apps/im/rpc" "im.go" "im-rpc.log" "10002" 8

echo ""
echo -e "${GREEN}========================================="
echo -e "         开始启动 API 服务"
echo -e "=========================================${NC}"

# 4. 启动 User API 服务
start_service "User API" "apps/user/api" "user.go" "user-api.log" "8888" 6

# 5. 启动 Social API 服务
start_service "Social API" "apps/social/api" "social.go" "social-api.log" "8881" 6

# 6. 启动 IM API 服务
start_service "IM API" "apps/im/api" "im.go" "im-api.log" "8882" 6

echo ""
echo -e "${GREEN}========================================="
echo -e "         开始启动 WebSocket 服务"
echo -e "=========================================${NC}"

# 7. 启动 IM WebSocket 服务
start_service "IM WebSocket" "apps/im/ws" "im.go" "im-ws.log" "10090" 8

echo ""
echo -e "${GREEN}========================================="
echo -e "         开始启动 消息队列 服务"
echo -e "=========================================${NC}"

# 8. 启动 Task MQ 服务
start_service "Task MQ" "apps/task/mq" "mq.go" "task-mq.log" "9001" 6

echo ""
echo -e "${GREEN}========================================="
echo -e "          所有服务启动完成！"
echo -e "=========================================${NC}"

echo -e "${BLUE}服务状态:${NC}"
echo -e "  • User RPC:     http://localhost:10010"
echo -e "  • Social RPC:   http://localhost:10001" 
echo -e "  • IM RPC:       http://localhost:10002"
echo -e "  • User API:     http://localhost:8888"
echo -e "  • Social API:   http://localhost:8881"
echo -e "  • IM API:       http://localhost:8882"
echo -e "  • IM WebSocket: ws://localhost:10090/ws"
echo -e "  • Task MQ:      http://localhost:9001"

echo ""
echo -e "${BLUE}日志文件位置: $LOG_DIR/${NC}"
echo -e "${YELLOW}按 Ctrl+C 停止所有服务${NC}"

# 保持脚本运行，等待用户中断
while true; do
    sleep 1
done
