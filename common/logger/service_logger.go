package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// ServiceLogger 带服务名的日志记录器
type ServiceLogger struct {
	serviceName string
	innerLogger logx.Logger
	fields      []logx.LogField
}

// NewServiceLogger 创建带服务名的日志记录器
func NewServiceLogger(serviceName string) *ServiceLogger {
	return &ServiceLogger{
		serviceName: serviceName,
		innerLogger: logx.WithContext(context.Background()),
	}
}

// WithContext 创建带上下文的服务日志记录器
func (s *ServiceLogger) WithContext(ctx context.Context) logx.Logger {
	return &ServiceLogger{
		serviceName: s.serviceName,
		innerLogger: logx.WithContext(ctx),
	}
}

// 实现 logx.Logger 接口的所有方法

func (s *ServiceLogger) Debug(v ...any) {
	s.innerLogger.Debugf("[%s] %s", s.serviceName, fmt.Sprint(v...))
}

func (s *ServiceLogger) Debugf(format string, v ...any) {
	s.innerLogger.Debugf("[%s] "+format, append([]any{s.serviceName}, v...)...)
}

func (s *ServiceLogger) Debugv(v any) {
	s.innerLogger.Infof("[%s] %v", s.serviceName, v)
}

func (s *ServiceLogger) Debugw(msg string, fields ...logx.LogField) {
	s.innerLogger.Infof("[%s] %s", s.serviceName, msg)
}

func (s *ServiceLogger) Error(v ...any) {
	s.innerLogger.Errorf("[%s] %s", s.serviceName, fmt.Sprint(v...))
}

func (s *ServiceLogger) Errorf(format string, v ...any) {
	s.innerLogger.Errorf("[%s] "+format, append([]any{s.serviceName}, v...)...)
}

func (s *ServiceLogger) Errorv(v any) {
	s.innerLogger.Errorf("[%s] %v", s.serviceName, v)
}

func (s *ServiceLogger) Errorw(msg string, fields ...logx.LogField) {
	s.innerLogger.Errorf("[%s] %s", s.serviceName, msg)
}

func (s *ServiceLogger) Info(v ...any) {
	s.innerLogger.Infof("[%s] %s", s.serviceName, fmt.Sprint(v...))
}

func (s *ServiceLogger) Infof(format string, v ...any) {
	s.innerLogger.Infof("[%s] "+format, append([]any{s.serviceName}, v...)...)
}

func (s *ServiceLogger) Infov(v any) {
	s.innerLogger.Infof("[%s] %v", s.serviceName, v)
}

func (s *ServiceLogger) Infow(msg string, fields ...logx.LogField) {
	s.innerLogger.Infof("[%s] %s", s.serviceName, msg)
}

func (s *ServiceLogger) Slow(v ...any) {
	s.innerLogger.Infof("[%s] %s", s.serviceName, fmt.Sprint(v...))
}

func (s *ServiceLogger) Slowf(format string, v ...any) {
	s.innerLogger.Infof("[%s] "+format, append([]any{s.serviceName}, v...)...)
}

func (s *ServiceLogger) Slowv(v any) {
	s.innerLogger.Infof("[%s] %v", s.serviceName, v)
}

func (s *ServiceLogger) Sloww(msg string, fields ...logx.LogField) {
	s.innerLogger.Infof("[%s] %s", s.serviceName, msg)
}

// WithCallerSkip 返回带调用者跳过的日志记录器
func (s *ServiceLogger) WithCallerSkip(skip int) logx.Logger {
	return &ServiceLogger{
		serviceName: s.serviceName,
		innerLogger: logx.WithCallerSkip(skip),
	}
}

// WithDuration 返回带持续时间的日志记录器
func (s *ServiceLogger) WithDuration(d time.Duration) logx.Logger {
	return &ServiceLogger{
		serviceName: s.serviceName,
		innerLogger: logx.WithDuration(d),
	}
}

// WithFields 返回带字段的日志记录器
func (s *ServiceLogger) WithFields(fields ...logx.LogField) logx.Logger {
	s.fields = append(s.fields, fields...)
	return s
}
