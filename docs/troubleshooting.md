# JXZY 问题诊断和解决方案

## 常见问题

### 1. 聊天流启动失败

**错误信息**：
```
failed to start chat stream
event: error
data: http: Server closed
```

**可能原因和解决方案**：

#### 1.1 服务未启动
**症状**：`dial tcp 127.0.0.1:8080: connect: connection refused`

**解决方案**：
```bash
# 启动所有服务
cd apis/api_chat && ./scripts/restart.sh
cd ../../bll/bll_context && ./scripts/restart.sh
cd ../../bs/bs_llm && ./scripts/restart.sh
```

#### 1.2 数据库连接失败
**症状**：`dial tcp [::1]:13309: connect: connection refused`

**解决方案**：
```bash
# 方案1：启动MySQL服务
brew install mysql
brew services start mysql

# 方案2：使用Docker启动MySQL
docker run --name mysql-jxzy \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=user \
  -e MYSQL_USER=dever \
  -e MYSQL_PASSWORD=dever \
  -p 13309:3306 \
  -d mysql:8.0

# 方案3：修改配置使用其他数据库
# 编辑 bll/bll_context/etc/bllcontext.yaml
```

#### 1.3 数据库模型未初始化
**症状**：`chat session model not initialized`

**解决方案**：
确保配置文件包含MySQL配置：
```yaml
MySQL:
  DataSource: dever:dever@tcp(localhost:13309)/user?charset=utf8mb4&parseTime=true&loc=Local
```

### 2. 日志配置问题

**症状**：日志分散在各个服务目录

**解决方案**：
```bash
# 运行日志配置测试
./scripts/test_logging.sh

# 查看统一日志
./scripts/view_logs.sh
```

### 3. 服务端口冲突

**症状**：`bind: address already in use`

**解决方案**：
```bash
# 查找占用端口的进程
lsof -i :8080 -i :8081 -i :8888

# 停止冲突的进程
kill -9 <PID>
```

## 服务启动顺序

正确的服务启动顺序：

1. **MySQL数据库** (端口 13309)
2. **BS服务** (端口 8081) - bs-llm
3. **BLL服务** (端口 8080) - bll-context
4. **API服务** (端口 8888) - chat-api

## 快速诊断脚本

```bash
# 检查所有服务状态
./scripts/test_chat.sh

# 检查日志配置
./scripts/test_logging.sh

# 查看实时日志
./scripts/view_logs.sh
```

## 数据库初始化

如果使用MySQL，需要先创建数据库和表：

```sql
-- 创建数据库
CREATE DATABASE user CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'dever'@'%' IDENTIFIED BY 'dever';
GRANT ALL PRIVILEGES ON user.* TO 'dever'@'%';
FLUSH PRIVILEGES;

-- 执行初始化脚本
-- 参考: bll/bll_context/internal/model/sql/session.sql
```

## 配置文件检查清单

确保以下配置文件正确：

1. `apis/api_chat/etc/chat-api.yaml` - API服务配置
2. `bll/bll_context/etc/bllcontext.yaml` - BLL服务配置
3. `bs/bs_llm/etc/bsllm.yaml` - BS服务配置

关键配置项：
- 日志路径：`Path: ../../logs`
- 数据库连接：`DataSource`
- RPC目标地址：`Target`

## 联系支持

如果问题仍然存在，请：

1. 收集完整的错误日志
2. 运行诊断脚本
3. 提供系统环境信息
4. 创建Issue报告
