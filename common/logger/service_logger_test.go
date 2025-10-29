package logger

import (
	"context"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
)

func TestServiceLogger(t *testing.T) {
	// 初始化日志系统
	InitLogger()

	// 创建带服务名的日志记录器
	serviceLogger := NewServiceLogger("test-service").WithContext(context.Background())

	// 测试各种日志级别
	t.Run("Test Info Logging", func(t *testing.T) {
		serviceLogger.Info("This is an info message")
		serviceLogger.Infof("This is a formatted info message: %s", "test")
	})

	t.Run("Test Error Logging", func(t *testing.T) {
		serviceLogger.Error("This is an error message")
		serviceLogger.Errorf("This is a formatted error message: %s", "test")
	})

	t.Run("Test Debug Logging", func(t *testing.T) {
		serviceLogger.Debug("This is a debug message")
		serviceLogger.Infof("This is a formatted debug message: %s", "test")
	})

	t.Run("Test WithFields", func(t *testing.T) {
		loggerWithFields := serviceLogger.WithFields(
			logx.Field("user_id", "123"),
			logx.Field("action", "test"),
		)
		loggerWithFields.Info("Message with fields")
	})

	t.Run("Test WithContext", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "request_id", "test-123")
		loggerWithContext := serviceLogger.WithContext(ctx)
		loggerWithContext.Info("Message with context")
	})
}

func TestServiceLoggerDifferentServices(t *testing.T) {
	// 初始化日志系统
	InitLogger()

	// 测试不同服务的日志记录器
	services := []string{"bll-context", "bs-llm", "api-chat"}

	for _, serviceName := range services {
		t.Run("Test "+serviceName, func(t *testing.T) {
			serviceLogger := NewServiceLogger(serviceName).WithContext(context.Background())
			serviceLogger.Infof("Service %s is running", serviceName)
			serviceLogger.Error("An error occurred in " + serviceName)
		})
	}
}

func TestServiceLoggerOutputFormat(t *testing.T) {
	// 初始化日志系统
	InitLogger()

	// 创建服务日志记录器
	serviceLogger := NewServiceLogger("test-service").WithContext(context.Background())

	// 测试日志输出格式
	t.Run("Test Service Name Prefix", func(t *testing.T) {
		// 这里我们无法直接捕获输出，但可以验证方法调用不会出错
		serviceLogger.Info("Test message")
		serviceLogger.Error("Test error")
		serviceLogger.Debug("Test debug")

		// 验证服务名是否正确设置
		if serviceLogger.(*ServiceLogger).serviceName != "test-service" {
			t.Errorf("Expected service name 'test-service', got '%s'", serviceLogger.(*ServiceLogger).serviceName)
		}
	})
}

func TestServiceLoggerWithFields(t *testing.T) {
	// 初始化日志系统
	InitLogger()

	serviceLogger := NewServiceLogger("test-service").WithContext(context.Background())

	t.Run("Test WithFields Method", func(t *testing.T) {
		loggerWithFields := serviceLogger.WithFields(
			logx.Field("key1", "value1"),
			logx.Field("key2", "value2"),
		)

		// 验证返回的是同一个 logger 实例
		if loggerWithFields != serviceLogger {
			t.Error("WithFields should return the same logger instance")
		}

		// 验证字段是否正确添加
		fields := loggerWithFields.(*ServiceLogger).fields
		if len(fields) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(fields))
		}
	})
}
