# 🚀 Easy-Chat 完整启动指南

## ⚠️ 重要提示

**代码已迁移到 Nacos 配置中心**，不能直接 `docker compose up -d` 启动所有服务！

必须按照以下步骤操作，否则应用服务会因为找不到配置而启动失败。

## 📋 启动前检查

### 1. 系统要求
- Docker 和 Docker Compose 已安装
- Go 1.21+ （如果要编译运行 Go 服务）
- 至少 4GB 可用内存
- 至少 10GB 可用磁盘空间

### 2. 架构确认
```bash
uname -m
# 输出应该是: x86_64 (AMD64)
```

## 🎯 分步启动流程

### 第一阶段：启动基础设施服务（必需）

#### 步骤 1: 启动 MySQL
```bash
docker compose up -d mysql
```

等待 MySQL 启动（约 20 秒）：
```bash
# 检查状态
docker compose ps mysql

# 查看日志，等待看到 "ready for connections"
docker compose logs -f mysql
```

#### 步骤 2: 初始化 Nacos 数据库
```bash
# 方法 1: 使用一键脚本（推荐）
./init-nacos.sh

# 方法 2: 手动执行
docker compose exec -i mysql mysql -uroot -peasy-chat < components/nacos/nacos-init.sql
```

#### 步骤 3: 启动 Nacos
```bash
docker compose up -d nacos
```

等待 Nacos 启动（约 30-60 秒）：
```bash
# 查看启动日志
docker compose logs -f nacos

# 等待看到类似信息：
# Nacos started successfully in stand alone mode
```

#### 步骤 4: 访问 Nacos 控制台并配置

**访问地址**: http://localhost:8848/nacos

**登录信息**:
- 用户名: `nacos`
- 密码: `nacos`

**配置 Nacos**（重要！）:

1. **创建命名空间**
   - 点击左侧 "命名空间"
   - 点击右上角 "新建命名空间"
   - 创建以下三个命名空间：
     * 命名空间ID: `user`，命名空间名: `user`
     * 命名空间ID: `im`，命名空间名: `im`
     * 命名空间ID: `social`，命名空间名: `social`

2. **上传配置文件**
   - 点击左侧 "配置管理" → "配置列表"
   - 选择对应的命名空间
   - 点击 "+" 创建配置
   - 为每个服务创建对应的配置文件：

   | 命名空间 | Group | Data ID | 配置来源 |
   |---------|-------|---------|----------|
   | user | DEFAULT_GROUP | user-api.yaml | etc/dev/user.yaml |
   | user | DEFAULT_GROUP | user-rpc.yaml | etc/dev/user.yaml (修改端口) |
   | im | DEFAULT_GROUP | im-api.yaml | etc/dev/im.yaml |
   | im | DEFAULT_GROUP | im-rpc.yaml | etc/dev/im.yaml (修改端口) |
   | im | DEFAULT_GROUP | im-ws.yaml | etc/dev/im.yaml (修改端口) |
   | social | DEFAULT_GROUP | social-api.yaml | etc/dev/social.yaml |
   | social | DEFAULT_GROUP | social-rpc.yaml | etc/dev/social.yaml (修改端口) |

   **配置示例**（user-api.yaml）:
   ```yaml
   Name: user-api
   Host: 0.0.0.0
   Port: 10001
   Mode: dev

   JwtAuth:
     AccessSecret: your-secret-key
     AccessExpire: 86400

   Mysql:
     DataSource: root:easy-chat@tcp(127.0.0.1:23306)/easy-chat?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

   CacheRedis:
     - Host: 127.0.0.1:16379
       Pass: easy-chat
       Type: node

   Etcd:
     Hosts:
       - 127.0.0.1:3379
     Key: user.rpc

   Log:
     ServiceName: user-api
     Mode: console
     Level: info
   ```

   **注意**: 
   - 将 `etc/dev/` 目录下的配置文件内容复制到 Nacos
   - 确保数据库地址、Redis地址等配置正确
   - 每个服务的端口要不同

#### 步骤 5: 启动其他基础服务
```bash
# 启动 Redis
docker compose up -d redis

# 启动 etcd（用于服务发现）
docker compose up -d etcd

# 启动 MongoDB
docker compose up -d mongo

# 启动 Kafka 相关
docker compose up -d zookeeper
sleep 10
docker compose up -d kafka

# 启动 Elasticsearch（可选）
docker compose up -d elasticsearch

# 启动 Logstash（可选）
docker compose up -d logstash

# 启动 API 网关
docker compose up -d apisix-dashboard apisix

# 启动链路追踪（可选）
docker compose up -d jeager
```

### 第二阶段：验证基础服务

```bash
# 查看所有服务状态
docker compose ps

# 应该看到以下服务都是 running 状态：
# - mysql
# - nacos
# - redis
# - etcd
# - mongo
# - zookeeper
# - kafka
# - elasticsearch (可选)
# - logstash (可选)
# - apisix
# - apisix-dashboard
# - jeager (可选)
```

### 第三阶段：启动 Go 应用服务

**确认 Nacos 中已配置好所有配置文件后**，启动 Go 服务：

```bash
# 1. 启动 user-rpc
cd apps/user/rpc
go run user.go -f etc/dev/user.yaml

# 2. 启动 user-api（新终端）
cd apps/user/api
go run user.go -f etc/dev/user.yaml

# 3. 启动 im-rpc（新终端）
cd apps/im/rpc
go run im.go -f etc/dev/im.yaml

# 4. 启动 im-api（新终端）
cd apps/im/api
go run im.go -f etc/dev/im.yaml

# 5. 启动 im-ws（新终端）
cd apps/im/ws
go run im.go -f etc/dev/im.yaml

# 6. 启动 social-rpc（新终端）
cd apps/social/rpc
go run social.go -f etc/dev/social.yaml

# 7. 启动 social-api（新终端）
cd apps/social/api
go run social.go -f etc/dev/social.yaml

# 8. 启动 task-mq（新终端，可选）
cd apps/task/mq
go run task.go -f etc/dev/task.yaml
```

## 🔍 服务访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| Nacos 控制台 | http://localhost:8848/nacos | 配置中心 |
| MySQL | localhost:23306 | 数据库 |
| Redis | localhost:16379 | 缓存 |
| MongoDB | localhost:47017 | NoSQL |
| Kafka | localhost:9092 | 消息队列 |
| Elasticsearch | http://localhost:9200 | 搜索引擎 |
| APISIX Dashboard | http://localhost:9000 | API网关控制台 |
| APISIX Gateway | http://localhost:9080 | API网关 |
| Jaeger UI | http://localhost:16686 | 链路追踪 |

## ✅ 验证清单

- [ ] MySQL 正常运行 (`docker compose ps mysql`)
- [ ] Nacos 正常运行并可访问控制台
- [ ] Nacos 中创建了 3 个命名空间 (user, im, social)
- [ ] Nacos 中上传了 7 个配置文件
- [ ] Redis、etcd、MongoDB 正常运行
- [ ] Kafka 和 Zookeeper 正常运行
- [ ] 所有 Go 服务能正常启动并连接到 Nacos
- [ ] 查看 Go 服务日志没有错误

## 🛠️ 快速命令

### 查看所有容器状态
```bash
docker compose ps
```

### 查看某个服务日志
```bash
docker compose logs -f nacos
docker compose logs -f mysql
docker compose logs -f redis
```

### 停止所有服务
```bash
# 停止 Docker 服务
docker compose down

# 停止 Go 服务（Ctrl+C 每个终端）
```

### 重启某个服务
```bash
docker compose restart nacos
docker compose restart mysql
```

### 完全清理并重新开始
```bash
# 警告：会删除所有数据！
docker compose down -v
rm -rf components/*/data
rm -rf components/*/logs
```

## ❓ 常见问题

### Q1: Nacos 启动失败
```bash
# 检查 MySQL 是否正常
docker compose ps mysql

# 检查 nacos 数据库是否存在
docker compose exec mysql mysql -uroot -peasy-chat -e "SHOW DATABASES;"

# 重新初始化
docker compose exec -i mysql mysql -uroot -peasy-chat < components/nacos/nacos-init.sql
docker compose restart nacos
```

### Q2: Go 服务连接 Nacos 失败
**检查项**：
- Nacos 是否启动成功？`docker compose ps nacos`
- Nacos 控制台能否访问？http://localhost:8848/nacos
- 配置文件是否已在 Nacos 中创建？
- 命名空间、Group、DataId 是否匹配？

### Q3: 配置文件找不到
确保：
1. 命名空间已创建（如 `user`）
2. 在正确的命名空间下创建配置
3. Data ID 完全匹配（如 `user-api.yaml`）
4. Group 为 `DEFAULT_GROUP`

### Q4: 端口冲突
如果某个端口被占用：
```bash
# 查看端口占用
netstat -tunlp | grep 8848

# 修改 docker-compose.yaml 中的端口映射
# 例如: "18848:8848"
```

## 📚 相关文档

- [NACOS-QUICKSTART.md](./NACOS-QUICKSTART.md) - Nacos 快速启动
- [NACOS-MIGRATION-GUIDE.md](./NACOS-MIGRATION-GUIDE.md) - 从 Sail 迁移到 Nacos 指南
- [components/nacos/README.md](./components/nacos/README.md) - Nacos 使用文档

## 🎉 启动成功标志

当你看到以下日志，说明启动成功：

```
Starting server at 0.0.0.0:10001...  # user-api
Starting rpc server at 0.0.0.0:10011...  # user-rpc
Starting server at 0.0.0.0:10002...  # im-api
Starting rpc server at 0.0.0.0:10012...  # im-rpc
start websocket server at 0.0.0.0:10003 .....  # im-ws
Starting server at 0.0.0.0:10004...  # social-api
Starting rpc server at 0.0.0.0:10014...  # social-rpc
```

祝使用愉快！ 🚀

