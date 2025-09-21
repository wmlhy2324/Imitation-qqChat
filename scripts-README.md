# Easy-Chat 启动脚本使用说明

## 🚀 快速开始

```bash
# 1. 给脚本添加执行权限
chmod +x *.sh

# 2. 启动所有服务
./start-local.sh

# 3. 测试服务状态
./test-services.sh

# 4. 停止所有服务
./stop-local.sh
```

## 📁 脚本文件说明

| 脚本文件 | 功能描述 |
|---------|---------|
| `start-local.sh` | 一键启动所有服务（RPC + API） |
| `stop-local.sh` | 一键停止所有服务 |
| `restart-local.sh` | 重启所有服务 |
| `test-services.sh` | 测试所有服务状态 |

## 🔧 启动顺序

脚本会按照以下顺序启动服务：

1. **RPC 服务**（后端服务）
   - User RPC (端口: 10000)
   - Social RPC (端口: 10001)  
   - IM RPC (端口: 10002)

2. **API 服务**（HTTP 接口）
   - User API (端口: 8888)
   - Social API (端口: 8881)
   - IM API (端口: 8882)

## 📋 依赖检查

启动前会自动检查以下依赖服务：

- ✅ MySQL (127.0.0.1:13306)
- ✅ Redis (127.0.0.1:16379)  
- ✅ ETCD (127.0.0.1:3379)
- ⚠️ MongoDB (127.0.0.1:47017) - IM 服务可选

## 📊 日志管理

- 所有服务日志保存在 `logs/` 目录
- 日志文件命名格式：`{service}-{type}.log`

```bash
# 实时查看日志
tail -f logs/user-rpc.log
tail -f logs/social-api.log

# 查看所有日志
ls -la logs/
```

## 🛠️ 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 查看端口占用
   lsof -i :8888
   
   # 强制停止
   ./stop-local.sh
   ```

2. **依赖服务未启动**
   ```bash
   # 检查依赖
   ./test-services.sh
   ```

3. **启动失败**
   ```bash
   # 查看具体错误
   tail -f logs/服务名.log
   ```

### 手动清理

```bash
# 杀死所有相关进程
pkill -f "go run"

# 清理日志文件
rm -rf logs/*
```

## 💡 使用技巧

1. **开发调试**: 使用 `tail -f logs/*.log` 同时监控多个日志
2. **快速重启**: 修改代码后使用 `./restart-local.sh`
3. **状态检查**: 定期运行 `./test-services.sh` 检查服务健康状态
4. **优雅停止**: 在启动脚本中按 `Ctrl+C` 会优雅停止所有服务

## ⚙️ 自定义配置

如需修改端口或其他配置，请编辑：
- 脚本中的端口定义
- 各服务的配置文件 (`etc/dev/*.yaml`)

---

**注意**: 推荐使用这些脚本而不是手动启动，可以避免依赖顺序和端口冲突问题。
