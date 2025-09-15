#!/bin/bash

# 使用Delve调试器调试登录流程
# 作者: AI Assistant

set -e

echo "=== Delve调试器 - 登录流程调试 ==="
echo "时间: $(date)"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# 检查Delve是否安装
check_dlv() {
    if ! command -v dlv &> /dev/null; then
        log_error "Delve未安装，请运行: go install github.com/go-delve/delve/cmd/dlv@latest"
        exit 1
    fi
    log_success "Delve已安装: $(dlv version | head -n1)"
}

# 启动依赖服务
start_dependencies() {
    log_info "启动依赖服务..."
    docker-compose up -d
    sleep 5
    log_success "依赖服务启动完成"
}

# 调试user-rpc服务
debug_user_rpc() {
    log_info "启动user-rpc调试模式..."
    
    cd apps/user/rpc
    
    # 使用Delve启动调试
    log_info "使用Delve启动user-rpc服务..."
    log_info "调试端口: 2345"
    log_info "服务端口: 10000"
    log_info ""
    log_info "调试命令:"
    log_info "  break apps/user/rpc/internal/logic/loginlogic.go:39  # 在Login函数开始处设置断点"
    log_info "  break apps/user/rpc/internal/logic/loginlogic.go:45  # 在用户查询处设置断点"
    log_info "  break apps/user/rpc/internal/logic/loginlogic.go:55  # 在密码验证处设置断点"
    log_info "  break apps/user/rpc/internal/logic/loginlogic.go:65  # 在Token生成处设置断点"
    log_info "  continue  # 继续执行"
    log_info "  next      # 下一步"
    log_info "  step      # 步入函数"
    log_info "  print req # 打印变量"
    log_info "  vars      # 查看所有变量"
    log_info ""
    
    # 启动调试器
    dlv debug user.go -- -f etc/dev/user.yaml
}

# 调试user-api服务
debug_user_api() {
    log_info "启动user-api调试模式..."
    
    cd apps/user/api
    
    # 使用Delve启动调试
    log_info "使用Delve启动user-api服务..."
    log_info "调试端口: 2346"
    log_info "服务端口: 8888"
    log_info ""
    log_info "调试命令:"
    log_info "  break apps/user/api/internal/logic/user/loginlogic.go:32  # 在Login函数开始处设置断点"
    log_info "  break apps/user/api/internal/logic/user/loginlogic.go:40  # 在RPC调用处设置断点"
    log_info "  break apps/user/api/internal/logic/user/loginlogic.go:50  # 在数据转换处设置断点"
    log_info "  break apps/user/api/internal/logic/user/loginlogic.go:60  # 在Redis操作处设置断点"
    log_info "  continue  # 继续执行"
    log_info "  next      # 下一步"
    log_info "  step      # 步入函数"
    log_info "  print req # 打印变量"
    log_info "  vars      # 查看所有变量"
    log_info ""
    
    # 启动调试器
    dlv debug user.go -- -f etc/dev/user.yaml
}

# 测试登录请求
test_login_request() {
    log_info "发送测试登录请求..."
    
    # 等待服务启动
    sleep 3
    
    # 发送登录请求
    curl -X POST http://localhost:8888/v1/user/login \
        -H "Content-Type: application/json" \
        -d '{
            "phone": "13800138000",
            "password": "123456"
        }' &
    
    log_success "登录请求已发送，请查看调试器输出"
}

# 显示调试帮助
show_debug_help() {
    echo "=== Delve调试器使用指南 ==="
    echo ""
    echo "1. 启动调试器后，会进入Delve命令行界面"
    echo ""
    echo "2. 常用调试命令："
    echo "   break <file:line>  - 设置断点"
    echo "   break <function>   - 在函数开始处设置断点"
    echo "   continue (c)       - 继续执行到下一个断点"
    echo "   next (n)           - 执行下一行代码"
    echo "   step (s)           - 步入函数"
    echo "   stepout            - 步出当前函数"
    echo "   print <variable>   - 打印变量值"
    echo "   vars               - 查看所有局部变量"
    echo "   args               - 查看函数参数"
    echo "   locals             - 查看局部变量"
    echo "   stack              - 查看调用栈"
    echo "   goroutines         - 查看所有goroutine"
    echo "   help               - 显示帮助"
    echo "   quit (q)           - 退出调试器"
    echo ""
    echo "3. 登录流程断点建议："
    echo "   API层断点："
    echo "     break apps/user/api/internal/logic/user/loginlogic.go:32"
    echo "     break apps/user/api/internal/logic/user/loginlogic.go:40"
    echo "     break apps/user/api/internal/logic/user/loginlogic.go:50"
    echo ""
    echo "   RPC层断点："
    echo "     break apps/user/rpc/internal/logic/loginlogic.go:39"
    echo "     break apps/user/rpc/internal/logic/loginlogic.go:45"
    echo "     break apps/user/rpc/internal/logic/loginlogic.go:55"
    echo "     break apps/user/rpc/internal/logic/loginlogic.go:65"
    echo ""
    echo "4. 调试步骤："
    echo "   1) 启动调试器"
    echo "   2) 设置断点"
    echo "   3) 运行 continue"
    echo "   4) 发送HTTP请求"
    echo "   5) 观察断点触发"
    echo "   6) 使用 print 查看变量"
    echo "   7) 使用 next/step 逐步调试"
    echo ""
}

# 主函数
main() {
    case "${1:-help}" in
        "rpc")
            check_dlv
            start_dependencies
            debug_user_rpc
            ;;
        "api")
            check_dlv
            start_dependencies
            debug_user_api
            ;;
        "test")
            test_login_request
            ;;
        "help")
            show_debug_help
            ;;
        *)
            echo "用法: $0 {rpc|api|test|help}"
            echo ""
            echo "命令说明:"
            echo "  rpc   - 调试user-rpc服务"
            echo "  api   - 调试user-api服务"
            echo "  test  - 发送测试登录请求"
            echo "  help  - 显示调试帮助"
            echo ""
            echo "调试流程:"
            echo "1. 运行 '$0 rpc' 调试RPC服务"
            echo "2. 运行 '$0 api' 调试API服务"
            echo "3. 运行 '$0 test' 发送测试请求"
            echo ""
            echo "建议先调试RPC服务，再调试API服务"
            ;;
    esac
}

main "$@"



