package common

import (
	"context"
	"fmt"
	"strings"
	"time"

	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/model"
	"jxzy/bs/bs_llm/internal/provider"
	"jxzy/bs/bs_llm/internal/provider/bailian"
	"jxzy/bs/bs_llm/internal/provider/doubao"
	"jxzy/bs/bs_llm/internal/svc"
	"jxzy/common/logger"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

// LLMCommon 公共的LLM逻辑
type LLMCommon struct {
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	providerManager *provider.Manager
	logger          logx.Logger
}

// NewLLMCommon 创建公共LLM逻辑实例
func NewLLMCommon(ctx context.Context, svcCtx interface{}) *LLMCommon {
	// 初始化供应商管理器
	manager := provider.NewManager()
	manager.Register("doubao", doubao.NewDoubaoProvider())
	manager.Register("bailian", bailian.NewBailianProvider())

	// 使用自定义的 ServiceLogger，在日志中显示服务名
	serviceLogger := logger.NewServiceLogger("bs-llm").WithContext(ctx)

	return &LLMCommon{
		ctx:             ctx,
		svcCtx:          svcCtx.(*svc.ServiceContext),
		providerManager: manager,
		logger:          serviceLogger,
	}
}

// GetContext 获取上下文
func (c *LLMCommon) GetContext() context.Context {
	return c.ctx
}

// GetServiceContext 获取服务上下文
func (c *LLMCommon) GetServiceContext() *svc.ServiceContext {
	return c.svcCtx
}

// GetProviderManager 获取供应商管理器
func (c *LLMCommon) GetProviderManager() *provider.Manager {
	return c.providerManager
}

// GetLogger 获取日志器
func (c *LLMCommon) GetLogger() logx.Logger {
	return c.logger
}

// InitializeCompletion 初始化完成记录
func (c *LLMCommon) InitializeCompletion(sceneCode string, messages []*bs_llm.ChatMessage, userId string) (*model.LlmCompletion, string) {
	requestId := uuid.New().String()

	completion := &model.LlmCompletion{
		SceneCode:    sceneCode,
		Prompt:       c.BuildPromptText(messages),
		RequestId:    requestId,
		UserId:       userId,
		InputTokens:  0,
		OutputTokens: 0,
		TotalTokens:  0,
		Status:       0, // 初始状态为失败，成功时会更新
		CreatedAt:    time.Now(),
	}

	c.logger.Infof("Completion record initialized - RequestId: %s, Prompt: %s", requestId, completion.Prompt)

	return completion, requestId
}

// GetSceneConfig 获取场景配置
func (c *LLMCommon) GetSceneConfig(sceneCode string) (*model.LlmScene, error) {
	c.logger.Infof("getSceneConfig called for scene_code: %s", sceneCode)

	if c.svcCtx.LlmSceneModel == nil {
		c.logger.Error("scene model not initialized")
		return nil, fmt.Errorf("scene model not initialized")
	}

	sceneInfo, err := c.svcCtx.LlmSceneModel.FindOneBySceneCode(c.ctx, sceneCode)
	if err != nil {
		if err == sqlc.ErrNotFound {
			c.logger.Errorf("Scene_code %s not found", sceneCode)
			return nil, fmt.Errorf("scene_code %s not found", sceneCode)
		}
		c.logger.Errorf("Failed to get scene info for %s: %v", sceneCode, err)
		return nil, fmt.Errorf("failed to get scene info: %w", err)
	}

	c.logger.Infof("Found scene config - SceneCode: %s, ProviderCode: %s, ModelCode: %s",
		sceneInfo.SceneCode, sceneInfo.ProviderCode, sceneInfo.ModelCode)
	return sceneInfo, nil
}

// GetProviderConfig 获取供应商配置
func (c *LLMCommon) GetProviderConfig(providerCode string) *provider.ProviderConfig {
	// 这里应该从配置文件或数据库中获取供应商配置
	// 为了简化，这里提供一个默认配置
	switch providerCode {
	case "doubao":
		return &provider.ProviderConfig{
			APIEndpoint: "https://ark.cn-beijing.volces.com/api/v3/chat/completions",
			APIKey:      c.svcCtx.Config.DoubaoAPIKey,
			Headers:     make(map[string]string),
			Timeout:     30,
			RetryCount:  3,
		}
	case "bailian":
		return &provider.ProviderConfig{
			APIEndpoint: "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
			APIKey:      c.svcCtx.Config.BailianAPIKey,
			Headers:     make(map[string]string),
			Timeout:     120,
			RetryCount:  3,
		}
	default:
		return &provider.ProviderConfig{
			Headers:    make(map[string]string),
			Timeout:    30,
			RetryCount: 3,
		}
	}
}

// BuildPromptText 构建提示词文本
func (c *LLMCommon) BuildPromptText(messages []*bs_llm.ChatMessage) string {
	var parts []string
	for _, msg := range messages {
		parts = append(parts, fmt.Sprintf("[%s]: %s", msg.Role, msg.Content))
	}
	return strings.Join(parts, "\n")
}

// EstimateTokens 估算文本的 token 数量
// 这是一个简化的估算方法，基于字符数来计算
// 对于英文，大约 4 个字符 = 1 个 token
// 对于中文，大约 1.5 个字符 = 1 个 token
func (c *LLMCommon) EstimateTokens(text string) int64 {
	if text == "" {
		return 0
	}

	// 计算中文字符数量
	chineseChars := 0
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			chineseChars++
		}
	}

	// 计算英文字符数量（包括空格、标点等）
	totalChars := len(text)
	englishChars := totalChars - chineseChars

	// 估算 token 数量
	// 中文：1.5 字符 = 1 token
	// 英文：4 字符 = 1 token
	chineseTokens := int64(float64(chineseChars) / 1.5)
	englishTokens := int64(float64(englishChars) / 4.0)

	totalTokens := chineseTokens + englishTokens

	// 确保至少返回 1 个 token
	if totalTokens < 1 {
		totalTokens = 1
	}

	return totalTokens
}

// SaveCompletion 保存问答记录
func (c *LLMCommon) SaveCompletion(completion *model.LlmCompletion) {
	c.logger.Infof("saveCompletion called for request_id: %s", completion.RequestId)

	// 检查context是否已取消
	select {
	case <-c.ctx.Done():
		c.logger.Infof("Context canceled, skipping completion save")
		return
	default:
	}

	if c.svcCtx.LlmCompletionModel == nil {
		c.logger.Error("completion model not initialized")
		return
	}

	// 使用独立的 context 进行数据库操作，避免 context canceled 错误
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.svcCtx.LlmCompletionModel.Insert(dbCtx, completion)
	if err != nil {
		c.logger.Errorf("Failed to save completion record: %v", err)
	} else {
		c.logger.Infof("Saved completion record: request_id=%s, status=%d, response_time=%.2fs",
			completion.RequestId, completion.Status, completion.ResponseTime.Float64)
	}
}

// ValidateSceneCode 验证场景码
func (c *LLMCommon) ValidateSceneCode(sceneCode string) error {
	if sceneCode == "" {
		c.logger.Error("scene_code is required")
		return fmt.Errorf("scene_code is required")
	}
	return nil
}

// GetOrDefaultUserId 获取用户ID，如果为空则使用默认值
func (c *LLMCommon) GetOrDefaultUserId(userId string) string {
	if userId != "" {
		return userId
	}
	return "anonymous"
}
