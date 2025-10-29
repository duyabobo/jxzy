package main

import (
	"context"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	// 初始化日志系统
	logger.InitLogger()

	// 创建不同服务的日志记录器
	bllContextLogger := logger.NewServiceLogger("bll-context").WithContext(context.Background())
	bsLlmLogger := logger.NewServiceLogger("bs-llm").WithContext(context.Background())
	apiChatLogger := logger.NewServiceLogger("api-chat").WithContext(context.Background())

	// 模拟不同服务的日志输出
	bllContextLogger.Info("BLL Context service started")
	bllContextLogger.Error("Failed to find session for update: context canceled")

	bsLlmLogger.Info("BS LLM service started")
	bsLlmLogger.Error("Failed to save completion record: context canceled")

	apiChatLogger.Info("API Chat service started")
	apiChatLogger.Error("Stream chat failed: http: Server closed")

	// 测试带字段的日志
	loggerWithFields := bllContextLogger.WithFields(
		logx.Field("user_id", "test_user_001"),
		logx.Field("session_id", "123"),
	)
	loggerWithFields.Info("Processing user request")

	// 测试带上下文的日志
	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	loggerWithContext := bsLlmLogger.WithContext(ctx)
	loggerWithContext.Info("Processing LLM request")
}
