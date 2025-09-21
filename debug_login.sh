#!/bin/bash

# 登录流程调试脚本
# 作者: AI Assistant
# 用途: 调试用户登录流程

set -e

echo "=== 登录流程调试脚本 ==="
echo "时间: $(date)"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        log_error "Go未安装"
        exit 1
    fi
    log_success "Go已安装: $(go version)"
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装"
        exit 1
    fi
    log_success "Docker已安装: $(docker --version)"
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装"
        exit 1
    fi
    log_success "Docker Compose已安装: $(docker-compose --version)"
    
    echo ""
}

# 启动依赖服务
start_dependencies() {
    log_info "启动依赖服务..."
    
    cd /testlinux/code/easy-chat
    
    # 检查Docker服务是否运行
    if ! docker info &> /dev/null; then
        log_error "Docker服务未运行"
        exit 1
    fi
    
    # 启动Docker Compose服务
    log_info "启动MySQL、Redis、ETCD等服务..."
    docker-compose up -d
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 10
    
    # 检查服务状态
    log_info "检查服务状态..."
    docker-compose ps
    
    log_success "依赖服务启动完成"
    echo ""
}

# 检查端口是否可用
check_port() {
    local port=$1
    local service=$2
    
    if netstat -tuln | grep -q ":$port "; then
        log_warning "$service 端口 $port 已被占用"
        return 1
    else
        log_success "$service 端口 $port 可用"
        return 0
    fi
}

# 启动user-rpc服务
start_user_rpc() {
    log_info "启动user-rpc服务..."
    
    # 检查端口
    check_port 10000 "user-rpc"
    
    cd /testlinux/code/easy-chat/apps/user/rpc
    
    # 后台启动user-rpc
    nohup go run user.go -f etc/dev/user.yaml > /tmp/user-rpc.log 2>&1 &
    USER_RPC_PID=$!
    
    # 等待服务启动
    sleep 5
    
    # 检查服务是否启动成功
    if ps -p $USER_RPC_PID > /dev/null; then
        log_success "user-rpc服务启动成功 (PID: $USER_RPC_PID)"
        echo $USER_RPC_PID > /tmp/user-rpc.pid
    else
        log_error "user-rpc服务启动失败"
        cat /tmp/user-rpc.log
        exit 1
    fi
    
    echo ""
}

# 启动user-api服务
start_user_api() {
    log_info "启动user-api服务..."
    
    # 检查端口
    check_port 8888 "user-api"
    
    cd /testlinux/code/easy-chat/apps/user/api
    
    # 后台启动user-api
    nohup go run user.go -f etc/dev/user.yaml > /tmp/user-api.log 2>&1 &
    USER_API_PID=$!
    
    # 等待服务启动
    sleep 5
    
    # 检查服务是否启动成功
    if ps -p $USER_API_PID > /dev/null; then
        log_success "user-api服务启动成功 (PID: $USER_API_PID)"
        echo $USER_API_PID > /tmp/user-api.pid
    else
        log_error "user-api服务启动失败"
        cat /tmp/user-api.log
        exit 1
    fi
    
    echo ""
}

# 测试登录接口
test_login() {
    log_info "测试登录接口..."
    
    # 等待服务完全启动
    sleep 3
    
    # 测试登录接口
    log_info "发送登录请求..."
    
    # 创建测试用户数据
    TEST_PHONE="13800138000"
    TEST_PASSWORD="123456"
    
    # 发送登录请求
    RESPONSE=$(curl -s -X POST http://localhost:8888/v1/user/login \
        -H "Content-Type: application/json" \
        -d "{
            \"phone\": \"$TEST_PHONE\",
            \"password\": \"$TEST_PASSWORD\"
        }")
    
    log_info "登录响应: $RESPONSE"
    
    # 检查响应
    if echo "$RESPONSE" | grep -q "token"; then
        log_success "登录测试成功"
        
        # 提取token
        TOKEN=$(echo "$RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        log_info "获取到Token: ${TOKEN:0:10}****"
        
        # 测试获取用户信息
        test_user_info "$TOKEN"
    else
        log_error "登录测试失败"
        log_info "可能的原因:"
        log_info "1. 用户不存在，需要先注册"
        log_info "2. 密码错误"
        log_info "3. 服务未正常启动"
    fi
    
    echo ""
}

# 测试获取用户信息
test_user_info() {
    local token=$1
    
    log_info "测试获取用户信息..."
    
    RESPONSE=$(curl -s -X GET http://localhost:8888/v1/user/user \
        -H "Authorization: Bearer $token")
    
    log_info "用户信息响应: $RESPONSE"
    
    if echo "$RESPONSE" | grep -q "id"; then
        log_success "获取用户信息成功"
    else
        log_warning "获取用户信息失败"
    fi
    
    echo ""
}

# 注册测试用户
register_test_user() {
    log_info "注册测试用户..."
    
    TEST_PHONE="13800138000"
    TEST_PASSWORD="123456"
    TEST_NICKNAME="测试用户"
    
    RESPONSE=$(curl -s -X POST http://localhost:8888/v1/user/register \
        -H "Content-Type: application/json" \
        -d "{
            \"phone\": \"$TEST_PHONE\",
            \"password\": \"$TEST_PASSWORD\",
            \"nickname\": \"$TEST_NICKNAME\",
            \"sex\": 1,
            \"avatar\": \"https://example.com/avatar.jpg\"
        }")
    
    log_info "注册响应: $RESPONSE"
    
    if echo "$RESPONSE" | grep -q "token"; then
        log_success "用户注册成功"
    else
        log_warning "用户注册失败或用户已存在"
    fi
    
    echo ""
}

# 显示日志
show_logs() {
    log_info "显示服务日志..."
    
    echo "=== user-api 日志 ==="
    tail -n 20 /tmp/user-api.log 2>/dev/null || echo "日志文件不存在"
    echo ""
    
    echo "=== user-rpc 日志 ==="
    tail -n 20 /tmp/user-rpc.log 2>/dev/null || echo "日志文件不存在"
    echo ""
}

# 清理服务
cleanup() {
    log_info "清理服务..."
    
    # 停止user-api
    if [ -f /tmp/user-api.pid ]; then
        USER_API_PID=$(cat /tmp/user-api.pid)
        if ps -p $USER_API_PID > /dev/null; then
            kill $USER_API_PID
            log_success "user-api服务已停止"
        fi
        rm -f /tmp/user-api.pid
    fi
    
    # 停止user-rpc
    if [ -f /tmp/user-rpc.pid ]; then
        USER_RPC_PID=$(cat /tmp/user-rpc.pid)
        if ps -p $USER_RPC_PID > /dev/null; then
            kill $USER_RPC_PID
            log_success "user-rpc服务已停止"
        fi
        rm -f /tmp/user-rpc.pid
    fi
    
    # 停止Docker服务
    cd /testlinux/code/easy-chat
    docker-compose down
    
    log_success "清理完成"
}

# 主函数
main() {
    case "${1:-help}" in
        "start")
            check_dependencies
            start_dependencies
            start_user_rpc
            start_user_api
            log_success "所有服务启动完成"
            ;;
        "test")
            register_test_user
            test_login
            ;;
        "logs")
            show_logs
            ;;
        "cleanup")
            cleanup
            ;;
        "full")
            check_dependencies
            start_dependencies
            start_user_rpc
            start_user_api
            register_test_user
            test_login
            show_logs
            ;;
        *)
            echo "用法: $0 {start|test|logs|cleanup|full}"
            echo ""
            echo "命令说明:"
            echo "  start   - 启动所有服务"
            echo "  test    - 测试登录功能"
            echo "  logs    - 显示服务日志"
            echo "  cleanup - 清理所有服务"
            echo "  full    - 完整流程：启动服务 + 测试 + 显示日志"
            echo ""
            echo "调试步骤:"
            echo "1. 运行 '$0 full' 进行完整测试"
            echo "2. 运行 '$0 logs' 查看详细日志"
            echo "3. 运行 '$0 cleanup' 清理环境"
            ;;
    esac
}

# 捕获中断信号
trap cleanup EXIT

# 执行主函数
main "$@"




