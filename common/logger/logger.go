package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	initialized bool
)

// InitLogger 初始化统一日志系统
func InitLogger() error {
	if initialized {
		return nil
	}

	// 创建logs目录
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// 配置日志 - 使用Go-Zero的完整配置
	logx.SetUp(logx.LogConf{
		ServiceName:         "jxzy",
		Mode:                "file",
		Path:                logsDir,
		Level:               "info",
		Encoding:            "json",
		TimeFormat:          "2006-01-02T15:04:05.000Z07:00",
		Compress:            true,
		KeepDays:            7,
		StackCooldownMillis: 100,
		MaxSize:             100,
		MaxBackups:          10,
		Stat:                true,
	})

	initialized = true
	return nil
}

// Access 记录访问日志
func Access(format string, args ...interface{}) {
	if !initialized {
		InitLogger()
	}
	logx.Infof("[ACCESS] "+format, args...)
}

// Error 记录错误日志
func Error(format string, args ...interface{}) {
	if !initialized {
		InitLogger()
	}
	logx.Errorf("[ERROR] "+format, args...)
}

// Info 记录信息日志
func Info(format string, args ...interface{}) {
	if !initialized {
		InitLogger()
	}
	logx.Infof("[INFO] "+format, args...)
}

// Warn 记录警告日志
func Warn(format string, args ...interface{}) {
	if !initialized {
		InitLogger()
	}
	logx.Errorf("[WARN] "+format, args...)
}

// Debug 记录调试日志
func Debug(format string, args ...interface{}) {
	if !initialized {
		InitLogger()
	}
	logx.Infof("[DEBUG] "+format, args...)
}

// Fatal 记录致命错误日志
func Fatal(format string, args ...interface{}) {
	if !initialized {
		InitLogger()
	}
	logx.Errorf("[FATAL] "+format, args...)
	os.Exit(1)
}

// WithContext 创建带上下文的日志记录器
func WithContext(ctx interface{}) *ContextLogger {
	if !initialized {
		InitLogger()
	}
	return &ContextLogger{
		ctx: ctx,
	}
}

// ContextLogger 带上下文的日志记录器
type ContextLogger struct {
	ctx interface{}
}

func (l *ContextLogger) Access(format string, args ...interface{}) {
	Access("[%v] "+format, append([]interface{}{l.ctx}, args...)...)
}

func (l *ContextLogger) Error(format string, args ...interface{}) {
	Error("[%v] "+format, append([]interface{}{l.ctx}, args...)...)
}

func (l *ContextLogger) Info(format string, args ...interface{}) {
	Info("[%v] "+format, append([]interface{}{l.ctx}, args...)...)
}

func (l *ContextLogger) Warn(format string, args ...interface{}) {
	Warn("[%v] "+format, append([]interface{}{l.ctx}, args...)...)
}

func (l *ContextLogger) Debug(format string, args ...interface{}) {
	Debug("[%v] "+format, append([]interface{}{l.ctx}, args...)...)
}

// RequestLogger 请求日志记录器
type RequestLogger struct {
	RequestID string
	UserID    string
	Method    string
	Path      string
	StartTime time.Time
}

// NewRequestLogger 创建请求日志记录器
func NewRequestLogger(requestID, userID, method, path string) *RequestLogger {
	return &RequestLogger{
		RequestID: requestID,
		UserID:    userID,
		Method:    method,
		Path:      path,
		StartTime: time.Now(),
	}
}

// LogRequest 记录请求开始
func (r *RequestLogger) LogRequest() {
	Access("ID:%s User:%s %s %s", r.RequestID, r.UserID, r.Method, r.Path)
}

// LogResponse 记录请求结束
func (r *RequestLogger) LogResponse(statusCode int, duration time.Duration) {
	Access("ID:%s User:%s %s %s Status:%d Duration:%v",
		r.RequestID, r.UserID, r.Method, r.Path, statusCode, duration)
}

// LogError 记录请求错误
func (r *RequestLogger) LogError(err error) {
	Error("ID:%s User:%s %s %s Error:%v",
		r.RequestID, r.UserID, r.Method, r.Path, err)
}
