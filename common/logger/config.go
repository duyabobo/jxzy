package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
)

// LogConfig 日志配置结构
type LogConfig struct {
	ServiceName         string `json:"serviceName" yaml:"serviceName"`
	Mode                string `json:"mode" yaml:"mode"`
	Path                string `json:"path" yaml:"path"`
	Level               string `json:"level" yaml:"level"`
	Encoding            string `json:"encoding" yaml:"encoding"`
	TimeFormat          string `json:"timeFormat" yaml:"timeFormat"`
	Compress            bool   `json:"compress" yaml:"compress"`
	KeepDays            int    `json:"keepDays" yaml:"keepDays"`
	StackCooldownMillis int    `json:"stackCooldownMillis" yaml:"stackCooldownMillis"`
	MaxSize             int    `json:"maxSize" yaml:"maxSize"`
	MaxBackups          int    `json:"maxBackups" yaml:"maxBackups"`
	Stat                bool   `json:"stat" yaml:"stat"`
}

// InitUnifiedLogger 初始化统一日志系统
func InitUnifiedLogger(serviceName string) error {
	if initialized {
		return nil
	}

	// 获取项目根目录的logs路径
	logsDir, err := getProjectLogsDir()
	if err != nil {
		return fmt.Errorf("failed to get project logs directory: %w", err)
	}

	// 创建logs目录
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// 配置日志
	config := logx.LogConf{
		ServiceName:         serviceName,
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
	}

	// 设置自定义日志格式，包含服务名称
	logx.SetUp(config)

	// 记录服务启动日志，包含服务名称
	logx.Infof("[%s] Service logger initialized", serviceName)

	initialized = true
	return nil
}

// getProjectLogsDir 获取项目根目录的logs路径
func getProjectLogsDir() (string, error) {
	// 从当前工作目录开始，向上查找go.mod文件来确定项目根目录
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上查找go.mod文件
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			// 找到go.mod文件，这就是项目根目录
			return filepath.Join(currentDir, "logs"), nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// 已经到达根目录，没有找到go.mod
			break
		}
		currentDir = parentDir
	}

	// 如果没有找到go.mod，使用当前目录下的logs
	return "logs", nil
}

// InitLoggerWithConfig 使用自定义配置初始化日志
func InitLoggerWithConfig(config LogConfig) error {
	if initialized {
		return nil
	}

	// 获取项目根目录的logs路径
	logsDir, err := getProjectLogsDir()
	if err != nil {
		return fmt.Errorf("failed to get project logs directory: %w", err)
	}

	// 创建logs目录
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// 使用配置中的路径，但确保指向项目根目录
	if config.Path != "" {
		// 如果配置的是相对路径，则相对于项目根目录
		if !filepath.IsAbs(config.Path) {
			config.Path = filepath.Join(logsDir, config.Path)
		}
	} else {
		config.Path = logsDir
	}

	// 设置默认值
	if config.ServiceName == "" {
		config.ServiceName = "jxzy"
	}
	if config.Mode == "" {
		config.Mode = "file"
	}
	if config.Level == "" {
		config.Level = "info"
	}
	if config.Encoding == "" {
		config.Encoding = "json"
	}
	if config.TimeFormat == "" {
		config.TimeFormat = "2006-01-02T15:04:05.000Z07:00"
	}
	if config.KeepDays == 0 {
		config.KeepDays = 7
	}
	if config.StackCooldownMillis == 0 {
		config.StackCooldownMillis = 100
	}
	if config.MaxSize == 0 {
		config.MaxSize = 100
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 10
	}

	// 配置日志
	logxConfig := logx.LogConf{
		ServiceName:         config.ServiceName,
		Mode:                config.Mode,
		Path:                config.Path,
		Level:               config.Level,
		Encoding:            config.Encoding,
		TimeFormat:          config.TimeFormat,
		Compress:            config.Compress,
		KeepDays:            config.KeepDays,
		StackCooldownMillis: config.StackCooldownMillis,
		MaxSize:             config.MaxSize,
		MaxBackups:          config.MaxBackups,
		Stat:                config.Stat,
	}

	logx.SetUp(logxConfig)

	// 记录服务启动日志，包含服务名称
	logx.Infof("[%s] Service logger initialized with custom config", config.ServiceName)

	initialized = true
	return nil
}
