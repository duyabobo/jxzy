# bs_llm RPC 测试文件

本目录包含了 `bs_llm` RPC 服务的测试文件，每个 RPC 接口都有独立的测试文件，每个文件只维护一个测试案例。

## 文件结构

```
test/
├── test_utils.go           # 共享测试工具和配置
├── llm_rpc_test.go         # LLM RPC 接口测试（非流式）
├── stream_llm_rpc_test.go  # StreamLLM RPC 接口测试（流式）
└── README.md              # 本说明文件
```

## 测试文件说明

### 1. test_utils.go
- 包含共享的测试工具和配置
- `MockStreamServer`: 模拟 gRPC 流服务器
- `GetTestConfig()`: 获取测试配置

### 2. llm_rpc_test.go
- 测试 `LLM` RPC 接口（非流式调用）
- 测试场景：使用百炼场景进行非流式对话
- 验证响应格式、内容完整性、token 使用情况等

### 3. stream_llm_rpc_test.go
- 测试 `StreamLLM` RPC 接口（流式调用）
- 测试场景：使用百炼场景进行流式对话
- 验证流式响应的连续性、完成状态、增量内容等

## 运行测试

### 运行所有测试
```bash
go test -v
```

### 运行特定测试
```bash
# 运行 LLM RPC 测试
go test -v -run TestLLM_RPC

# 运行 StreamLLM RPC 测试
go test -v -run TestStreamLLM_RPC

# 运行两个 RPC 测试
go test -v -run "TestLLM_RPC|TestStreamLLM_RPC"
```

## 测试特点

1. **独立性**: 每个 RPC 接口有独立的测试文件
2. **单一职责**: 每个测试文件只维护一个测试案例
3. **错误处理**: 包含完善的错误处理和 API 配置问题处理
4. **详细验证**: 验证响应格式、内容完整性、业务逻辑等
5. **日志记录**: 提供详细的测试执行日志

## 注意事项

- 测试需要正确的配置文件（`bsllm.yaml`）
- 测试会实际调用外部 API，需要网络连接
- 如果 API 配置有问题，测试会跳过而不是失败
- 测试会向数据库写入记录，请确保数据库配置正确
