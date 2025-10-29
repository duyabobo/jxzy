# JXZY 日志统一配置改进总结

## 问题描述

之前各个服务（api/bll/bs）的日志分散在各个服务目录下，没有统一收集到项目根目录的logs文件夹，导致日志管理困难。

## 解决方案

### 1. 配置文件修改

修改了所有服务的YAML配置文件，将日志路径统一指向项目根目录：

**修改前：**
```yaml
Log:
  Path: logs  # 相对路径，指向服务目录下的logs
```

**修改后：**
```yaml
Log:
  Path: ../../logs  # 相对路径，指向项目根目录的logs
```

### 2. 统一日志系统

创建了统一的日志管理系统：

- `common/logger/logger.go` - 基础日志功能
- `common/logger/config.go` - 高级日志配置

### 3. 服务启动代码更新

在每个服务的启动代码中添加了统一日志初始化：

```go
// 初始化统一日志系统
if err := logger.InitUnifiedLogger("service-name"); err != nil {
    logx.Errorf("Failed to initialize logger: %v", err)
    return
}
```

### 4. 工具脚本

创建了便捷的日志管理工具：

- `scripts/test_logging.sh` - 测试日志配置
- `scripts/view_logs.sh` - 查看和管理日志

## 修改的文件列表

### 配置文件
- `apis/api_chat/etc/chat-api.yaml`
- `bll/bll_context/etc/bllcontext.yaml`
- `bs/bs_llm/etc/bsllm.yaml`

### 启动代码
- `apis/api_chat/chat.go`
- `bll/bll_context/bllcontext.go`
- `bs/bs_llm/bsllm.go`

### 新增文件
- `common/logger/config.go`
- `scripts/test_logging.sh`
- `scripts/view_logs.sh`
- `logs/.gitkeep`
- `docs/logging.md`

### 更新文件
- `common/logger/logger.go`
- `README.md`

## 配置特性

### 1. 统一路径
所有服务的日志都收集到项目根目录的`logs/`文件夹

### 2. 自动归档
- 按日期自动归档日志文件
- 自动压缩旧日志文件
- 保留7天的日志记录

### 3. 文件轮转
- 单个日志文件最大100MB
- 最多保留10个备份文件
- 自动分割大文件

### 4. JSON格式
采用结构化JSON格式，便于日志分析和处理

### 5. 服务标识
每个日志条目都包含服务名称，便于区分不同服务的日志

## 使用方法

### 1. 测试配置
```bash
./scripts/test_logging.sh
```

### 2. 查看日志
```bash
./scripts/view_logs.sh
```

### 3. 启动服务
服务启动后会自动创建统一的日志文件：
```bash
# 启动各个服务
cd apis/api_chat && ./scripts/compile.sh && ./scripts/restart.sh
cd ../../bll/bll_context && ./scripts/compile.sh && ./scripts/restart.sh
cd ../../bs/bs_llm && ./scripts/compile.sh && ./scripts/restart.sh
```

## 预期效果

### 1. 统一日志收集
所有服务的日志都会收集到`logs/jxzy.log`文件中

### 2. 便于问题排查
- 可以在一个地方查看所有服务的日志
- 支持按服务名称过滤日志
- 支持实时查看错误日志

### 3. 自动化管理
- 自动归档和压缩
- 自动清理过期日志
- 自动文件轮转

### 4. 开发友好
- 提供便捷的日志查看工具
- 详细的配置文档
- 测试脚本验证配置

## 验证步骤

1. 运行测试脚本验证配置
2. 启动各个服务
3. 检查logs目录下是否生成了统一的日志文件
4. 使用日志查看工具验证功能
5. 确认所有服务的日志都正确收集

## 注意事项

1. 确保logs目录有正确的写入权限
2. 服务启动时会自动创建logs目录
3. 日志文件会自动按日期归档
4. 使用相对路径确保从任何目录启动服务都能正确工作
