# Easy Chat - 仿QQ聊天系统

一个基于Go语言和微服务架构的即时通讯系统，模仿QQ聊天功能的实现。

## 📋 项目概述

Easy Chat 是一个完整的即时通讯解决方案，采用微服务架构设计，支持实时聊天、好友管理、群组管理等核心功能。项目使用Go-Zero框架构建，结合WebSocket实现实时通信，使用MongoDB存储聊天记录，MySQL存储用户数据。

## 🚀 主要特性

- **实时通讯**: 基于WebSocket的实时消息推送
- **微服务架构**: 模块化设计，易于扩展和维护  
- **用户系统**: 用户注册、登录、资料管理
- **社交功能**: 好友添加、好友管理、群组管理
- **消息系统**: 私聊、群聊、消息存储与检索
- **高可用性**: 支持集群部署，负载均衡
- **容器化部署**: 完整的Docker部署方案
- **监控日志**: 集成ELK日志收集与分析

## 🛠 技术栈

### 后端框架
- **Go 1.19**: 主要开发语言
- **Go-Zero**: 微服务框架
- **gRPC**: 服务间通信
- **WebSocket**: 实时通信协议

### 数据存储
- **MySQL 5.7**: 用户数据、社交关系存储
- **MongoDB 4.0**: 聊天记录存储  
- **Redis**: 缓存、会话管理
- **Kafka**: 消息队列

### 基础设施
- **Etcd**: 服务发现与配置管理
- **APISIX**: API网关
- **Docker**: 容器化部署
- **ElasticSearch + Logstash + Kibana**: 日志分析
- **Jaeger**: 分布式链路追踪

### 第三方库
- **gorilla/websocket**: WebSocket支持
- **golang-jwt**: JWT认证
- **go-redis**: Redis客户端
- **mongo-driver**: MongoDB驱动

## 📁 项目结构

```
easy-chat/
├── apps/                           # 应用服务目录
│   ├── user/                       # 用户服务
│   │   ├── api/                    # HTTP API服务
│   │   ├── rpc/                    # gRPC服务
│   │   └── models/                 # 数据模型
│   ├── social/                     # 社交服务
│   │   ├── api/                    # HTTP API服务
│   │   ├── rpc/                    # gRPC服务
│   │   └── socialmodels/           # 数据模型
│   ├── im/                         # 即时通讯服务
│   │   ├── api/                    # HTTP API服务
│   │   ├── rpc/                    # gRPC服务
│   │   ├── ws/                     # WebSocket服务
│   │   └── immodels/               # 数据模型
│   └── task/                       # 任务处理服务
│       └── mq/                     # 消息队列处理
├── components/                     # 基础组件配置
│   ├── mysql/                      # MySQL配置
│   ├── redis/                      # Redis配置
│   ├── mongo/                      # MongoDB配置
│   ├── etcd/                       # Etcd配置
│   ├── apisix/                     # API网关配置
│   ├── elasticsearch/              # ES配置
│   ├── logstash/                   # Logstash配置
│   └── kibana/                     # Kibana配置
├── deploy/                         # 部署相关文件
│   ├── dockerfile/                 # Docker构建文件
│   ├── mk/                         # Makefile模块
│   ├── script/                     # 部署脚本
│   └── sql/                        # 数据库初始化脚本
├── pkg/                            # 公共包
│   ├── constants/                  # 常量定义
│   ├── middleware/                 # 中间件
│   ├── interceptor/                # 拦截器
│   ├── resultx/                    # 响应处理
│   └── xerr/                       # 错误处理
└── test/                           # 测试相关文件
```

## 🔧 环境要求

- Go 1.19+
- Docker & Docker Compose
- Make

## 📦 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/wmlhy2324/Imitation-qqChat.git
cd easy-chat
```

### 2. 启动基础服务

```bash
# 启动Docker基础服务(MySQL, Redis, MongoDB, Kafka等)
make install-docker
```

### 3. 编译并启动所有微服务

```bash
# 编译所有服务
make release-test

# 或者单独启动各个服务
make user-api-dev      # 用户API服务
make user-rpc-dev      # 用户RPC服务  
make social-api-dev    # 社交API服务
make social-rpc-dev    # 社交RPC服务
make im-api-dev        # 即时通讯API服务
make im-rpc-dev        # 即时通讯RPC服务
make im-ws-dev         # WebSocket服务
make task-mq-dev       # 任务队列服务
```

### 4. 安装服务到Docker

```bash
# 安装所有服务到Docker容器
make install-server
```

## 🌐 服务端口

| 服务 | 端口 | 描述 |
|-----|------|------|
| MySQL | 13306 | 用户数据存储 |
| Redis | 16379 | 缓存服务 |
| MongoDB | 47017 | 聊天记录存储 |
| Etcd | 3379 | 服务发现 |
| Kafka | 9092 | 消息队列 |
| APISIX | 9080 | API网关 |
| APISIX Dashboard | 9000 | 网关管理界面 |
| ElasticSearch | 9200 | 日志存储 |
| Kibana | 5601 | 日志分析界面 |

## 🔌 API文档

### 用户服务 (User Service)
- `POST /v1/user/register` - 用户注册
- `POST /v1/user/login` - 用户登录  
- `GET /v1/user/user` - 获取用户信息 (需要JWT认证)

### 社交服务 (Social Service)
- 好友管理 (添加、删除、查询好友)
- 好友请求处理
- 群组管理 (创建、加入、退出群组)

### 即时通讯服务 (IM Service)
- WebSocket连接管理
- 实时消息推送
- 聊天记录查询
- 会话管理

## 🗄️ 数据库设计

### 用户表 (users)
```sql
CREATE TABLE `users` (
  `id` varchar(24) NOT NULL,
  `avatar` varchar(191) NOT NULL DEFAULT '',
  `nickname` varchar(24) NOT NULL,
  `phone` varchar(20) NOT NULL,
  `password` varchar(191) DEFAULT NULL,
  `status` tinyint DEFAULT NULL,
  `sex` tinyint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
);
```

### 消息记录存储
- 使用MongoDB存储聊天记录
- 支持文本、图片、文件等多种消息类型
- 按会话ID分片存储，提高查询效率

## 🔍 监控与日志

项目集成了完整的监控和日志系统：

- **日志收集**: Filebeat收集应用日志
- **日志处理**: Logstash处理和转换日志
- **日志存储**: ElasticSearch存储日志数据
- **日志分析**: Kibana提供可视化分析界面
- **链路追踪**: Jaeger提供分布式链路追踪
- **指标监控**: 支持Prometheus指标收集

访问 `http://localhost:5601` 查看Kibana日志分析界面。

## 🚀 部署指南

### 开发环境部署

1. 确保Docker和Docker Compose已安装
2. 执行 `make install-docker` 启动基础服务
3. 执行 `make install-server` 部署应用服务

### 生产环境部署

1. 修改各服务的配置文件 (位于 `apps/*/api/etc/` 和 `apps/*/rpc/etc/`)
2. 构建生产镜像
3. 使用Kubernetes或Docker Swarm进行集群部署
4. 配置负载均衡和服务发现

## 🔧 配置说明

### 环境变量
- `MYSQL_ROOT_PASSWORD`: MySQL root密码 (默认: easy-chat)
- `MONGO_INITDB_ROOT_USERNAME`: MongoDB用户名 (默认: root)  
- `MONGO_INITDB_ROOT_PASSWORD`: MongoDB密码 (默认: easy-chat)

### Kafka主题配置
- `ws2ms_chat`: WebSocket到微服务的消息
- `ms2ps_chat`: 微服务到推送服务的消息  
- `msg_to_mongo`: 消息持久化到MongoDB

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 📄 许可证

本项目基于MIT许可证开源。查看 [LICENSE](LICENSE) 文件了解更多细节。

## 📞 联系方式

如有问题或建议，欢迎通过以下方式联系：

- GitHub Issues: [https://github.com/wmlhy2324/Imitation-qqChat/issues](https://github.com/wmlhy2324/Imitation-qqChat/issues)

---

⭐ 如果这个项目对你有帮助，请给我们一个Star！