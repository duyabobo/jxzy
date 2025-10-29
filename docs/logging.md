# JXZY 统一日志配置

## 概述

JXZY项目使用Go-Zero框架的统一日志系统，所有服务的日志都会收集到项目根目录的`logs`文件夹中。

## 日志文件结构

```
logs/
├── .gitkeep                    # 确保logs目录被git跟踪
├── jxzy.log                   # 统一的应用日志
├── jxzy.log.2025-08-16        # 按日期归档的日志文件
└── jxzy.log.2025-08-16.gz     # 压缩的归档日志文件
```

## 配置说明

### 1. 服务配置文件

每个服务的配置文件（YAML）中的Log部分：

```yaml
Log:
  ServiceName: chat-api          # 服务名称
  Mode: file                     # 日志模式：file
  Path: ../../logs              # 日志路径：指向项目根目录的logs
  Level: info                   # 日志级别
  Encoding: json                # 日志格式：JSON
  TimeFormat: 2006-01-02T15:04:05.000Z07:00  # 时间格式
  Compress: true                # 是否压缩
  KeepDays: 7                   # 保留天数
  StackCooldownMillis: 100      # 堆栈冷却时间
  MaxSize: 100                  # 单个日志文件最大大小(MB)
  MaxBackups: 10                # 最大备份文件数
  Stat: true                    # 是否启用统计
```

### 2. 统一日志初始化

在服务启动时，会自动初始化统一日志系统：

```go
// 初始化统一日志系统
if err := logger.InitUnifiedLogger("service-name"); err != nil {
    logx.Errorf("Failed to initialize logger: %v", err)
    return
}
```

### 3. 日志使用

使用Go-Zero的logx包记录日志：

```go
import "github.com/zeromicro/go-zero/core/logx"

// 信息日志
logx.Infof("Service started successfully")

// 错误日志
logx.Errorf("Failed to connect to database: %v", err)

// 调试日志
logx.Infof("Processing request: %s", requestID)

// 警告日志
logx.Errorf("Warning: resource usage is high")
```

## 服务配置

### API服务 (chat-api)
- 配置文件：`apis/api_chat/etc/chat-api.yaml`
- 服务名称：`chat-api`
- 端口：8888

### BLL服务 (bll-context)
- 配置文件：`bll/bll_context/etc/bllcontext.yaml`
- 服务名称：`bll-context`
- 端口：8080

### BS服务 (bs-llm)
- 配置文件：`bs/bs_llm/etc/bsllm.yaml`
- 服务名称：`bs-llm`
- 端口：8081

## 日志级别

- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息
- `fatal`: 致命错误

## 日志格式

日志采用JSON格式，包含以下字段：

```json
{
  "@timestamp": "2025-08-16T12:00:59.625+08:00",
  "caller": "svc/servicecontext.go:43",
  "content": "Successfully connected to LLM RPC via direct connection",
  "level": "info"
}
```

## 日志轮转

- 单个日志文件最大大小：100MB
- 最大备份文件数：10个
- 日志保留天数：7天
- 自动压缩归档文件

## 注意事项

1. 所有服务的日志都会统一收集到项目根目录的`logs`文件夹
2. 日志文件会自动按日期归档和压缩
3. 使用相对路径`../../logs`确保从各个服务目录都能正确指向项目根目录
4. 日志初始化在配置加载后、服务启动前进行

## 故障排除

如果日志没有正确收集到统一位置：

1. 检查配置文件中的`Path`设置是否正确
2. 确认项目根目录存在`logs`文件夹
3. 检查服务启动时的日志初始化是否成功
4. 验证文件权限是否正确
