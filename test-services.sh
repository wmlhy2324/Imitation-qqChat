#!/bin/bash

# Easy-Chat 服务测试脚本

echo "========================================="
echo "        Easy-Chat 服务测试"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试函数
test_service() {
    local service_name=$1
    local url=$2
    local expected_response=$3
    
    echo -n "测试 $service_name: "
    
    # 使用 curl 测试服务
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)
    
    if [ "$response" = "200" ] || [ "$response" = "404" ] || [ "$response" = "401" ]; then
        echo -e "${GREEN}✓ 运行正常 (HTTP $response)${NC}"
        return 0
    else
        echo -e "${RED}✗ 服务异常 (HTTP $response)${NC}"
        return 1
    fi
}

test_rpc_service() {
    local service_name=$1
    local host=$2
    local port=$3
    
    echo -n "测试 $service_name: "
    
    # 测试端口是否开放
    if nc -z "$host" "$port" 2>/dev/null; then
        echo -e "${GREEN}✓ 端口开放 ($host:$port)${NC}"
        return 0
    else
        echo -e "${RED}✗ 端口未开放 ($host:$port)${NC}"
        return 1
    fi
}

echo -e "${BLUE}测试 RPC 服务...${NC}"
test_rpc_service "User RPC" "127.0.0.1" "10000"
test_rpc_service "Social RPC" "127.0.0.1" "10001"
test_rpc_service "IM RPC" "127.0.0.1" "10002"

echo ""
echo -e "${BLUE}测试 API 服务...${NC}"
test_service "User API" "http://localhost:8888/health" 
test_service "Social API" "http://localhost:8881/health"
test_service "IM API" "http://localhost:8882/health"

echo ""
echo -e "${BLUE}测试依赖服务...${NC}"
test_rpc_service "MySQL" "127.0.0.1" "13306"
test_rpc_service "Redis" "127.0.0.1" "16379"
test_rpc_service "ETCD" "127.0.0.1" "3379"
test_rpc_service "MongoDB" "127.0.0.1" "47017"

echo ""
echo "测试完成！"
