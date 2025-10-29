package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/common"
	"jxzy/bs/bs_llm/internal/provider"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type LLMLogic struct {
	common *common.LLMCommon
	logx.Logger
}

func NewLLMLogic(ctx context.Context, svcCtx interface{}) *LLMLogic {
	commonLogic := common.NewLLMCommon(ctx, svcCtx)

	return &LLMLogic{
		common: commonLogic,
		Logger: commonLogic.GetLogger(),
	}
}

// 非流式LLM调用
func (l *LLMLogic) LLM(in *bs_llm.LLMRequest) (*bs_llm.LLMResponse, error) {
	startTime := time.Now()
	requestId := uuid.New().String() // 生成请求ID

	l.Logger.Infof("LLM called with scene_code: %s, request_id: %s, messages_count: %d",
		in.SceneCode, requestId, len(in.Messages))

	// 获取 user_id
	userId := l.common.GetOrDefaultUserId(in.UserId)

	// 初始化完成记录
	completion, _ := l.common.InitializeCompletion(in.SceneCode, in.Messages, userId)
	completion.RequestId = requestId // 使用生成的请求ID

	// 延迟执行：保存问答记录
	defer func() {
		responseTime := time.Since(startTime).Seconds()
		completion.ResponseTime = sql.NullFloat64{Float64: responseTime, Valid: true}
		l.Logger.Infof("Saving completion record - RequestId: %s, ResponseTime: %.2fs", requestId, responseTime)
		l.common.SaveCompletion(completion)
	}()

	// 1. 验证场景码
	if err := l.common.ValidateSceneCode(in.SceneCode); err != nil {
		completion.ErrorMsg = sql.NullString{String: err.Error(), Valid: true}
		return nil, err
	}

	// 2. 获取场景配置
	l.Logger.Debug("Getting scene configuration")
	sceneConfig, err := l.common.GetSceneConfig(in.SceneCode)
	if err != nil {
		completion.ErrorMsg = sql.NullString{String: err.Error(), Valid: true}
		return nil, err
	}

	// 更新completion记录中的模型和供应商信息
	completion.ModelCode = sceneConfig.ModelCode
	completion.ProviderCode = sceneConfig.ProviderCode

	l.Logger.Infof("Resolved provider: %s, model: %s for scene: %s", sceneConfig.ProviderCode, sceneConfig.ModelCode, sceneConfig.SceneCode)
	l.Logger.Infof("Scene config - Temperature: %f, MaxTokens: %d, EnableStream: %d",
		sceneConfig.Temperature, sceneConfig.MaxTokens, sceneConfig.EnableStream)

	// 3. 获取供应商
	l.Logger.Infof("Getting provider: %s", sceneConfig.ProviderCode)
	llmProvider := l.common.GetProviderManager().GetProvider(sceneConfig.ProviderCode)
	if llmProvider == nil {
		err := fmt.Errorf("provider %s not supported", sceneConfig.ProviderCode)
		completion.ErrorMsg = sql.NullString{String: err.Error(), Valid: true}
		return nil, err
	}

	// 4. 构建请求
	l.Logger.Debug("Building LLM request")
	req := &provider.LLMRequest{
		Messages:    common.ConvertToProviderMessages(in.Messages),
		ModelCode:   sceneConfig.ModelCode,
		Temperature: sceneConfig.Temperature,
		MaxTokens:   sceneConfig.MaxTokens,
		Stream:      false, // 非流式调用
		ExtraParams: in.ExtraParams,
		Config:      l.common.GetProviderConfig(sceneConfig.ProviderCode),
	}

	l.Logger.Infof("LLM request built - Model: %s, Temperature: %f, MaxTokens: %d, Messages: %d",
		req.ModelCode, req.Temperature, req.MaxTokens, len(req.Messages))

	// 5. 调用非流式LLM
	l.Logger.Debug("Calling non-stream LLM")
	providerResp, err := llmProvider.CallLLM(l.common.GetContext(), req)
	if err != nil {
		l.Logger.Errorf("Failed to call non-stream LLM: %v", err)
		completion.ErrorMsg = sql.NullString{String: fmt.Sprintf("failed to call non-stream LLM: %v", err), Valid: true}
		return nil, fmt.Errorf("failed to call non-stream LLM: %w", err)
	}

	l.Logger.Infof("LLM response received - Completion: %s, FinishReason: %s", providerResp.Content, providerResp.FinishReason)

	// 6. 更新completion记录为成功状态
	completion.Completion = sql.NullString{String: providerResp.Content, Valid: true}
	completion.Status = 1 // 成功

	// 设置 token 使用情况
	if providerResp.PromptTokens > 0 || providerResp.CompletionTokens > 0 || providerResp.TotalTokens > 0 {
		completion.InputTokens = providerResp.PromptTokens
		completion.OutputTokens = providerResp.CompletionTokens
		completion.TotalTokens = providerResp.TotalTokens
		l.Logger.Infof("LLM provided usage - Input: %d, Output: %d, Total: %d",
			providerResp.PromptTokens, providerResp.CompletionTokens, providerResp.TotalTokens)
	} else {
		// 如果没有 usage 信息，自己计算 token 数量
		inputTokens := l.common.EstimateTokens(completion.Prompt)
		outputTokens := l.common.EstimateTokens(providerResp.Content)
		totalTokens := inputTokens + outputTokens

		completion.InputTokens = inputTokens
		completion.OutputTokens = outputTokens
		completion.TotalTokens = totalTokens
		l.Logger.Infof("Estimated token usage - Input: %d, Output: %d, Total: %d",
			inputTokens, outputTokens, totalTokens)
	}

	// 7. 构建gRPC响应
	llmResp := &bs_llm.LLMResponse{
		Completion:   providerResp.Content,
		ModelId:      sceneConfig.ModelCode,
		FinishReason: providerResp.FinishReason,
	}

	// 添加usage信息
	if providerResp.PromptTokens > 0 || providerResp.CompletionTokens > 0 || providerResp.TotalTokens > 0 {
		llmResp.Usage = &bs_llm.LLMUsage{
			PromptTokens:     providerResp.PromptTokens,
			CompletionTokens: providerResp.CompletionTokens,
			TotalTokens:      providerResp.TotalTokens,
		}
	} else {
		// 估算token使用情况
		inputTokens := l.common.EstimateTokens(completion.Prompt)
		outputTokens := l.common.EstimateTokens(providerResp.Content)
		totalTokens := inputTokens + outputTokens

		llmResp.Usage = &bs_llm.LLMUsage{
			PromptTokens:     inputTokens,
			CompletionTokens: outputTokens,
			TotalTokens:      totalTokens,
		}
	}

	l.Logger.Info("LLM logic completed successfully")
	return llmResp, nil
}
