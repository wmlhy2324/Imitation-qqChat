#!/bin/bash

# Nacos 初始化脚本
# 用于自动初始化 Nacos 配置中心

set -e

echo "========================================="
echo "  Nacos 配置中心初始化脚本"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 步骤 1: 启动 MySQL
echo -e "${YELLOW}[1/5]${NC} 启动 MySQL 容器..."
docker compose up -d mysql

# 步骤 2: 等待 MySQL 就绪
echo -e "${YELLOW}[2/5]${NC} 等待 MySQL 启动完成..."
sleep 5

# 检查 MySQL 是否就绪
MAX_TRIES=30
COUNTER=0
until docker compose exec mysql mysqladmin ping -uroot -peasy-chat --silent &> /dev/null ; do
    COUNTER=$((COUNTER+1))
    if [ $COUNTER -gt $MAX_TRIES ]; then
        echo -e "${RED}错误: MySQL 启动超时${NC}"
        exit 1
    fi
    echo -n "."
    sleep 1
done
echo ""
echo -e "${GREEN}✓ MySQL 已就绪${NC}"

# 步骤 3: 检查数据库是否已初始化
echo -e "${YELLOW}[3/5]${NC} 检查 Nacos 数据库..."
DB_EXISTS=$(docker compose exec mysql mysql -uroot -peasy-chat -e "SHOW DATABASES LIKE 'nacos';" 2>/dev/null | grep -c "nacos" || true)

if [ "$DB_EXISTS" -eq "0" ]; then
    echo "Nacos 数据库不存在，开始初始化..."
    # 步骤 4: 初始化 Nacos 数据库
    echo -e "${YELLOW}[4/5]${NC} 初始化 Nacos 数据库..."
    docker compose exec -T mysql mysql -uroot -peasy-chat < components/nacos/nacos-init.sql
    echo -e "${GREEN}✓ Nacos 数据库初始化完成${NC}"
else
    echo -e "${GREEN}✓ Nacos 数据库已存在，跳过初始化${NC}"
fi

# 步骤 5: 启动 Nacos
echo -e "${YELLOW}[5/5]${NC} 启动 Nacos 容器..."
docker compose up -d nacos

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}  Nacos 初始化完成！${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "📌 访问信息："
echo "   控制台地址: http://localhost:8848/nacos"
echo "   用户名: nacos"
echo "   密码: nacos"
echo ""
echo "📌 查看日志："
echo "   docker compose logs -f nacos"
echo ""
echo "⏳ Nacos 启动需要约 30 秒，请稍候..."
echo ""

