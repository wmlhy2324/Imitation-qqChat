# 🚀 Nacos 配置中心快速启动指南

## 📋 变更说明

已将配置中心从 `sail` 替换为 `Nacos`：
- ✅ **Nacos**: 支持 x86_64 架构，功能强大
- ❌ **sail**: 仅支持 ARM64 架构（已注释）

## ⚡ 快速启动（推荐）

使用自动化脚本一键初始化：

```bash
./init-nacos.sh
```

脚本会自动完成：
1. 检查 Docker 和 Docker Compose
2. 创建必要的目录
3. 启动 MySQL 容器
4. 初始化 Nacos 数据库
5. 启动 Nacos 服务
6. 验证服务状态

## 🔧 手动启动步骤

如果你想手动操作，请按以下步骤：

### 1. 启动 MySQL
```bash
docker compose up -d mysql
```

### 2. 等待 MySQL 就绪（约 20 秒）
```bash
# 查看 MySQL 日志
docker compose logs -f mysql

# 等待看到类似这样的日志：
# [Server] /usr/sbin/mysqld: ready for connections
```

### 3. 初始化 Nacos 数据库
```bash
docker compose exec -i mysql mysql -uroot -peasy-chat < components/nacos/nacos-init.sql
```

### 4. 启动 Nacos
```bash
docker compose up -d nacos
```

### 5. 等待 Nacos 启动（约 30-60 秒）
```bash
docker compose logs -f nacos

# 等待看到类似这样的日志：
# Nacos started successfully in stand alone mode
```

## 🌐 访问 Nacos 控制台

启动成功后，访问：

- **URL**: http://localhost:8848/nacos
- **用户名**: `nacos`
- **密码**: `nacos`

## 📦 启动所有服务

初始化 Nacos 后，启动所有服务：

```bash
docker compose up -d
```

查看服务状态：
```bash
docker compose ps
```

## 🔍 服务列表

| 服务名称 | 端口 | 说明 |
|---------|------|------|
| nacos | 8848, 9848 | 配置中心 + 服务发现 |
| mysql | 23306 | 数据库 |
| redis | 16379 | 缓存 |
| mongo | 47017 | NoSQL 数据库 |
| etcd | 3379, 3380 | 分布式键值存储 |
| kafka | 9092 | 消息队列 |
| zookeeper | 2181 | Kafka 依赖 |
| elasticsearch | 9200, 9300 | 搜索引擎 |
| logstash | 5044, 50000, 9600 | 日志收集 |
| apisix | 9080, 9443 | API 网关 |
| apisix-dashboard | 9000 | API 网关控制台 |
| jeager | 16686 | 链路追踪 |

## 🛠️ 常用命令

```bash
# 启动所有服务
docker compose up -d

# 停止所有服务
docker compose down

# 查看服务状态
docker compose ps

# 查看服务日志
docker compose logs -f [服务名]

# 重启某个服务
docker compose restart [服务名]

# 查看 Nacos 日志
docker compose logs -f nacos

# 进入 MySQL 容器
docker compose exec mysql mysql -uroot -peasy-chat
```

## 🔧 配置 Go 项目使用 Nacos

详细的集成文档请查看：[components/nacos/README.md](components/nacos/README.md)

简单示例：

```go
import (
    "github.com/nacos-group/nacos-sdk-go/v2/clients"
    "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
    "github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// 创建配置客户端
clientConfig := constant.ClientConfig{
    NamespaceId: "",
    TimeoutMs:   5000,
    Username:    "nacos",
    Password:    "nacos",
}

serverConfigs := []constant.ServerConfig{
    {IpAddr: "127.0.0.1", Port: 8848},
}

configClient, _ := clients.NewConfigClient(
    vo.NacosClientParam{
        ClientConfig:  &clientConfig,
        ServerConfigs: serverConfigs,
    },
)

// 获取配置
content, _ := configClient.GetConfig(vo.ConfigParam{
    DataId: "application.yaml",
    Group:  "DEFAULT_GROUP",
})
```

## ❓ 故障排查

### Nacos 启动失败

```bash
# 1. 检查 MySQL 是否正常运行
docker compose ps mysql

# 2. 检查数据库是否已初始化
docker compose exec mysql mysql -uroot -peasy-chat -e "SHOW DATABASES LIKE 'nacos';"

# 3. 查看 Nacos 详细日志
docker compose logs nacos

# 4. 重新初始化
docker compose down nacos
./init-nacos.sh
```

### 无法访问控制台

```bash
# 检查容器状态
docker compose ps nacos

# 检查端口占用
netstat -tunlp | grep 8848

# 等待更长时间（Nacos 启动需要 30-60 秒）
docker compose logs -f nacos
```

### 端口冲突

如果某个端口被占用，可以修改 `docker-compose.yaml` 中的端口映射：

```yaml
ports:
  - "18848:8848"  # 将 8848 改为 18848
```

## 📚 更多信息

- Nacos 详细使用文档：[components/nacos/README.md](components/nacos/README.md)
- Nacos 官方文档：https://nacos.io/zh-cn/docs/
- Go SDK 文档：https://github.com/nacos-group/nacos-sdk-go

## ✅ 验证清单

- [ ] MySQL 容器正常运行
- [ ] Nacos 数据库已创建
- [ ] Nacos 容器正常运行
- [ ] 可以访问 http://localhost:8848/nacos
- [ ] 可以使用 nacos/nacos 登录控制台
- [ ] 其他服务正常启动

---

**享受使用 Nacos 配置中心！** 🎉

