package test

import (
	"testing"
	"time"

	"jxzy/common/errorx"
	"jxzy/common/utils"

	"github.com/stretchr/testify/assert"
)

// TestUUIDGeneration 测试UUID生成功能
func TestUUIDGeneration(t *testing.T) {
	// 测试标准UUID生成
	uuid1 := utils.GenerateUUID()
	uuid2 := utils.GenerateUUID()

	// 验证UUID格式和唯一性
	assert.NotEmpty(t, uuid1)
	assert.NotEmpty(t, uuid2)
	assert.NotEqual(t, uuid1, uuid2)
	assert.True(t, utils.IsValidUUID(uuid1))
	assert.True(t, utils.IsValidUUID(uuid2))

	// 测试短UUID生成
	shortUuid1 := utils.GenerateShortUUID()
	shortUuid2 := utils.GenerateShortUUID()

	// 验证短UUID格式和唯一性
	assert.NotEmpty(t, shortUuid1)
	assert.NotEmpty(t, shortUuid2)
	assert.NotEqual(t, shortUuid1, shortUuid2)
	assert.NotContains(t, shortUuid1, "-") // 短UUID不包含横线
	assert.NotContains(t, shortUuid2, "-")

	// 测试无效UUID验证
	assert.False(t, utils.IsValidUUID("invalid-uuid"))
	assert.False(t, utils.IsValidUUID(""))
	assert.False(t, utils.IsValidUUID("12345"))
}

// TestTimeUtils 测试时间工具函数
func TestTimeUtils(t *testing.T) {
	// 获取当前时间戳
	timestamp := utils.GetCurrentTimestamp()
	milliTimestamp := utils.GetCurrentMillisTimestamp()

	// 验证时间戳格式
	assert.Greater(t, timestamp, int64(0))
	assert.Greater(t, milliTimestamp, int64(0))
	assert.Greater(t, milliTimestamp, timestamp*1000-1000) // 毫秒时间戳应该更大

	// 测试时间格式化
	now := time.Now()
	formatted := utils.FormatTime(now)
	assert.NotEmpty(t, formatted)
	assert.Contains(t, formatted, "-") // 应该包含日期分隔符
	assert.Contains(t, formatted, ":") // 应该包含时间分隔符

	// 测试时间解析
	parsed, err := utils.ParseTime(formatted)
	assert.NoError(t, err)
	assert.Equal(t, now.Year(), parsed.Year())
	assert.Equal(t, now.Month(), parsed.Month())
	assert.Equal(t, now.Day(), parsed.Day())

	// 测试时间戳转换
	convertedTime := utils.TimestampToTime(timestamp)
	assert.WithinDuration(t, time.Now(), convertedTime, time.Minute)
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 测试预定义错误
	assert.Equal(t, 200, errorx.ErrSuccess.GetCode())
	assert.Equal(t, 400, errorx.ErrParamError.GetCode())
	assert.Equal(t, 401, errorx.ErrUnauthorized.GetCode())
	assert.Equal(t, 404, errorx.ErrNotFound.GetCode())
	assert.Equal(t, 500, errorx.ErrSystemError.GetCode())

	// 测试业务错误
	assert.Equal(t, 10001, errorx.ErrUserNotFound.GetCode())
	assert.Equal(t, 11001, errorx.ErrContextNotFound.GetCode())
	assert.Equal(t, 12001, errorx.ErrPromptNotFound.GetCode())
	assert.Equal(t, 13001, errorx.ErrLLMNotAvailable.GetCode())
	assert.Equal(t, 14001, errorx.ErrDocumentNotFound.GetCode())

	// 测试自定义错误创建
	customErr := errorx.NewCodeError(99999, "自定义错误")
	assert.Equal(t, 99999, customErr.GetCode())
	assert.Equal(t, "自定义错误", customErr.GetMsg())

	// 测试格式化错误创建
	formatErr := errorx.NewCodeErrorf(99998, "格式化错误: %s", "测试参数")
	assert.Equal(t, 99998, formatErr.GetCode())
	assert.Contains(t, formatErr.GetMsg(), "测试参数")

	// 测试错误字符串表示
	errStr := customErr.Error()
	assert.Contains(t, errStr, "99999")
	assert.Contains(t, errStr, "自定义错误")
}

// TestSessionIdGeneration 测试会话ID生成
func TestSessionIdGeneration(t *testing.T) {
	// 模拟会话ID生成逻辑
	generateSessionId := func(userId string) string {
		return "sess_" + userId + "_" + utils.GenerateShortUUID()
	}

	// 生成测试会话ID
	userId := "user123"
	sessionId1 := generateSessionId(userId)
	sessionId2 := generateSessionId(userId)

	// 验证会话ID格式
	assert.NotEmpty(t, sessionId1)
	assert.NotEmpty(t, sessionId2)
	assert.NotEqual(t, sessionId1, sessionId2)
	assert.Contains(t, sessionId1, "sess_")
	assert.Contains(t, sessionId1, userId)
	assert.Contains(t, sessionId2, "sess_")
	assert.Contains(t, sessionId2, userId)
}

// TestConfigValidation 测试配置验证
func TestConfigValidation(t *testing.T) {
	// 模拟配置验证函数
	validateConfig := func(config map[string]interface{}) error {
		if config["model"] == "" {
			return errorx.NewCodeError(400, "模型配置不能为空")
		}
		if temp, ok := config["temperature"].(float64); ok && (temp < 0 || temp > 2) {
			return errorx.NewCodeError(400, "温度参数必须在0-2之间")
		}
		return nil
	}

	// 测试有效配置
	validConfig := map[string]interface{}{
		"model":       "doubao",
		"temperature": 0.7,
		"max_tokens":  1000,
	}
	err := validateConfig(validConfig)
	assert.NoError(t, err)

	// 测试无效配置 - 缺少模型
	invalidConfig1 := map[string]interface{}{
		"model":       "",
		"temperature": 0.7,
	}
	err = validateConfig(invalidConfig1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "模型配置不能为空")

	// 测试无效配置 - 温度参数超范围
	invalidConfig2 := map[string]interface{}{
		"model":       "doubao",
		"temperature": 3.0,
	}
	err = validateConfig(invalidConfig2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "温度参数必须在0-2之间")
}
