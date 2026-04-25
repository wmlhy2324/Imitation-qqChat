# 📚 从 Sail 迁移到 Nacos 完整指南

## 🎯 迁移概述

本项目已从 **Sail 配置中心**迁移到 **Nacos 配置中心**，主要原因是 Sail 只有 ARM64 镜像，不兼容 x86_64 架构。

### 变更内容

| 项目 | 修改前 (Sail) | 修改后 (Nacos) |
|------|--------------|----------------|
| 配置中心 | Sail (仅 ARM64) | Nacos (支持 x86_64/ARM64) |
| 配置存储 | etcd | MySQL + Nacos |
| Web 管理界面 | http://localhost:8108 | http://localhost:8848/nacos |
| 适配器文件 | pkg/configserver/sail.go | pkg/configserver/nacos.go |
| 服务启动参数 | NewSail() | NewNacos() |

## ✅ 已完成的代码修改

### 1. 创建 Nacos 适配器
- 文件：`pkg/configserver/nacos.go`
- 实现了与 Sail 相同的 `ConfigServer` 接口
- 支持配置动态更新

### 2. 修改所有服务启动文件

已修改以下 7 个服务：

#### User 服务
- ✅ `apps/user/api/user.go` - user-api.yaml
- ✅ `apps/user/rpc/user.go` - user-rpc.yaml

#### IM 服务
- ✅ `apps/im/api/im.go` - im-api.yaml
- ✅ `apps/im/rpc/im.go` - im-rpc.yaml
- ✅ `apps/im/ws/im.go` - im-ws.yaml

#### Social 服务
- ✅ `apps/social/api/social.go` - social-api.yaml
- ✅ `apps/social/rpc/social.go` - social-rpc.yaml

### 3. 添加依赖
```bash
go get github.com/nacos-group/nacos-sdk-go/v2@latest
```

## 🚀 Nacos 初始化步骤

### 步骤 1: 启动 Nacos 服务

使用一键脚本：
```bash
./init-nacos.sh
```

或手动启动：
```bash
# 1. 启动 MySQL
docker compose up -d mysql

# 2. 等待 20 秒，初始化数据库
sleep 20
docker compose exec -i mysql mysql -uroot -peasy-chat < components/nacos/nacos-init.sql

# 3. 启动 Nacos
docker compose up -d nacos

# 4. 等待 Nacos 启动完成（约 30-60 秒）
docker compose logs -f nacos
```

### 步骤 2: 访问 Nacos 控制台

- URL: http://localhost:8848/nacos
- 用户名: `nacos`
- 密码: `nacos`

## 📝 配置文件迁移

### 在 Nacos 中创建配置

你需要在 Nacos 中创建以下 7 个配置文件：

| 命名空间 | Group | Data ID | 说明 |
|---------|-------|---------|------|
| user | DEFAULT_GROUP | user-api.yaml | User API 配置 |
| user | DEFAULT_GROUP | user-rpc.yaml | User RPC 配置 |
| im | DEFAULT_GROUP | im-api.yaml | IM API 配置 |
| im | DEFAULT_GROUP | im-rpc.yaml | IM RPC 配置 |
| im | DEFAULT_GROUP | im-ws.yaml | IM WebSocket 配置 |
| social | DEFAULT_GROUP | social-api.yaml | Social API 配置 |
| social | DEFAULT_GROUP | social-rpc.yaml | Social RPC 配置 |

### 操作步骤

#### 1. 创建命名空间

在 Nacos 控制台：

1. 点击左侧菜单 **"命名空间"**
2. 点击右上角 **"新建命名空间"**
3. 创建三个命名空间：
   - 命名空间名：`user`，命名空间ID：`user`
   - 命名空间名：`im`，命名空间ID：`im`
   - 命名空间名：`social`，命名空间ID：`social`

#### 2. 创建配置文件

在 Nacos 控制台：

1. 点击左侧菜单 **"配置管理" → "配置列表"**
2. 在顶部选择对应的命名空间（如 `user`）
3. 点击右上角 **"+"** 按钮
4. 填写配置信息：
   - **Data ID**: 如 `user-api.yaml`
   - **Group**: `DEFAULT_GROUP`
   - **配置格式**: `YAML`
   - **配置内容**: 粘贴你的 YAML 配置

### 配置文件示例

#### user-api.yaml 示例
```yaml
Name: user-api
Host: 0.0.0.0
Port: 10001
Mode: dev

JwtAuth:
  AccessSecret: your-access-secret-key-here
  AccessExpire: 86400

Mysql:
  DataSource: root:easy-chat@tcp(127.0.0.1:23306)/easy-chat?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

CacheRedis:
  - Host: 127.0.0.1:16379
    Pass: easy-chat
    Type: node

Log:
  ServiceName: user-api
  Mode: console
  Level: info
  Encoding: plain
  Path: logs
  KeepDays: 7
```

#### user-rpc.yaml 示例
```yaml
Name: user-rpc
ListenOn: 0.0.0.0:10011
Mode: dev

Etcd:
  Hosts:
    - 127.0.0.1:3379
  Key: user.rpc

Mysql:
  DataSource: root:easy-chat@tcp(127.0.0.1:23306)/easy-chat?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

CacheRedis:
  - Host: 127.0.0.1:16379
    Pass: easy-chat
    Type: node

Log:
  ServiceName: user-rpc
  Mode: console
  Level: info
  Encoding: plain
  Path: logs
  KeepDays: 7
```

**注意**：其他配置文件（im-api.yaml, im-rpc.yaml, im-ws.yaml, social-api.yaml, social-rpc.yaml）格式类似，只需修改对应的服务名称和端口即可。

### 从本地文件导入配置

如果你已有本地配置文件（如 `etc/dev/user-api.yaml`），可以：

1. 打开 Nacos 控制台
2. 创建配置时，选择 **"从文件导入"**
3. 选择你的本地配置文件
4. 点击发布

## 🔧 配置参数说明

### 代码中的 Nacos 配置

```go
configserver.NewNacos(&configserver.NacosConfig{
    Addr:      "127.0.0.1",        // Nacos 服务器地址
    Namespace: "user",              // 命名空间 ID
    Group:     "DEFAULT_GROUP",     // 配置分组
    DataId:    "user-api.yaml",     // 配置文件 Data ID
    Username:  "nacos",             // Nacos 用户名
    Password:  "nacos",             // Nacos 密码
    LogLevel:  "warn",              // 日志级别
})
```

### 参数修改

如果你的 Nacos 部署在其他地址，修改 `Addr` 参数：

```go
Addr: "192.168.1.100",  // Nacos 服务器 IP
```

如果修改了 Nacos 用户名密码，相应修改 `Username` 和 `Password`。

## 🧪 验证迁移

### 1. 检查 Nacos 配置

```bash
# 查看 Nacos 容器状态
docker compose ps nacos

# 查看 Nacos 日志
docker compose logs nacos

# 访问 Nacos 控制台
http://localhost:8848/nacos
```

### 2. 测试服务启动

启动一个服务测试：

```bash
# 启动 user-api 服务
cd apps/user/api
go run user.go -f etc/dev/user.yaml
```

如果看到类似日志，说明成功连接 Nacos：
```
配置发生变化 - Namespace: user, Group: DEFAULT_GROUP, DataId: user-api.yaml
Starting server at 0.0.0.0:10001...
```

### 3. 测试配置热更新

1. 在 Nacos 控制台修改配置
2. 点击 **"发布"**
3. 观察服务日志，应该会看到配置更新提示
4. 服务会自动重启应用新配置

## 🆚 Sail vs Nacos 对比

| 特性 | Sail | Nacos |
|-----|------|-------|
| 架构支持 | 仅 ARM64 ❌ | x86_64 + ARM64 ✅ |
| 配置存储 | etcd | MySQL |
| Web 界面 | 简单 | 功能丰富 ✅ |
| 社区支持 | 小众 | 阿里巴巴，活跃 ✅ |
| 文档 | 较少 | 完善 ✅ |
| 服务发现 | 依赖 etcd | 内置 ✅ |
| 配置版本管理 | ❌ | ✅ |
| 配置回滚 | ❌ | ✅ |
| 灰度发布 | ❌ | ✅ |

## 🔄 回滚到 Sail（可选）

如果需要回滚到 Sail：

### 1. 修改代码

将所有 `NewNacos()` 改回 `NewSail()`：

```go
// 回滚示例
configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
    ETCDEndpoints:  "127.0.0.1:3379",
    ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
    Namespace:      "user",
    Configs:        "user-api.yaml",
    ConfigFilePath: "./etc/conf",
    LogLevel:       "DEBUG",
}))
```

### 2. 取消 Sail 注释

编辑 `docker-compose.yaml`，取消 sail 服务的注释：

```yaml
sail:
  image: ccr.ccs.tencentyun.com/hyy-yu/sail:latest
  container_name: sail
  platform: linux/arm64  # 需要 QEMU 模拟器
  ports:
    - "8108:8108"
  volumes:
    - "./components/sail/compose-cfg.toml:/app/cfg.toml"
  restart: always
  networks:
    easy-chat:
```

### 3. 重启服务

```bash
docker compose up -d sail
```

## ❓ 常见问题

### Q1: Nacos 启动失败？
**A**: 确保 MySQL 已启动并初始化了 nacos 数据库：
```bash
docker compose exec mysql mysql -uroot -peasy-chat -e "SHOW DATABASES LIKE 'nacos';"
```

### Q2: 服务连接 Nacos 失败？
**A**: 检查 Nacos 地址和端口是否正确：
```go
Addr: "127.0.0.1",  // 确保能访问
```

### Q3: 找不到配置文件？
**A**: 
- 检查命名空间是否创建
- 检查 Data ID 和 Group 是否匹配
- 查看 Nacos 日志：`docker compose logs nacos`

### Q4: 配置热更新不生效？
**A**: 
- 确保在创建 NewNacos 时提供了 onChange 回调
- 检查配置格式是否正确（必须是 YAML）

### Q5: 依赖下载失败？
**A**: 
```bash
# 使用国内代理
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```

## 📚 参考资料

- [Nacos 官方文档](https://nacos.io/zh-cn/docs/what-is-nacos.html)
- [Nacos Go SDK](https://github.com/nacos-group/nacos-sdk-go)
- [NACOS-QUICKSTART.md](./NACOS-QUICKSTART.md) - Nacos 快速启动指南
- [components/nacos/README.md](./components/nacos/README.md) - Nacos 详细使用文档

## ✅ 迁移检查清单

- [ ] Nacos 服务已启动
- [ ] Nacos 数据库已初始化
- [ ] 可以访问 Nacos 控制台（http://localhost:8848/nacos）
- [ ] 已创建 3 个命名空间（user, im, social）
- [ ] 已创建 7 个配置文件
- [ ] 代码已修改为使用 NewNacos()
- [ ] 已安装 nacos-sdk-go 依赖
- [ ] 已测试至少一个服务能正常启动
- [ ] 已测试配置热更新功能

完成以上所有步骤后，迁移就完成了！🎉

---

**如有问题，请查看日志或参考文档。祝使用愉快！**

