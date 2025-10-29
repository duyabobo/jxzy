package logic

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/common"
	"jxzy/bs/bs_llm/internal/provider"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type StreamLLMLogic struct {
	common *common.LLMCommon
	logx.Logger
}

func NewStreamLLMLogic(ctx context.Context, svcCtx interface{}) *StreamLLMLogic {
	commonLogic := common.NewLLMCommon(ctx, svcCtx)

	return &StreamLLMLogic{
		common: commonLogic,
		Logger: commonLogic.GetLogger(),
	}
}

// 流式LLM调用
func (l *StreamLLMLogic) StreamLLM(in *bs_llm.LLMRequest, stream bs_llm.BsLlmService_StreamLLMServer) error {
	startTime := time.Now()
	requestId := uuid.New().String() // 生成请求ID

	l.Logger.Infof("StreamLLM called with scene_code: %s, request_id: %s, messages_count: %d",
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
		return err
	}

	// 2. 获取场景配置
	l.Logger.Debug("Getting scene configuration")
	sceneConfig, err := l.common.GetSceneConfig(in.SceneCode)
	if err != nil {
		completion.ErrorMsg = sql.NullString{String: err.Error(), Valid: true}
		return err
	}

	// 更新completion记录中的模型和供应商信息
	completion.ModelCode = sceneConfig.ModelCode
	completion.ProviderCode = sceneConfig.ProviderCode

	l.Logger.Infof("Resolved provider: %s, model: %s for scene: %s", sceneConfig.ProviderCode, sceneConfig.ModelCode, sceneConfig.SceneCode)
	l.Logger.Infof("Scene config - Temperature: %f, MaxTokens: %d, EnableStream: %d",
		sceneConfig.Temperature, sceneConfig.MaxTokens, sceneConfig.EnableStream)

	// 3. 检查是否启用流式输出
	if sceneConfig.EnableStream == 0 {
		err := fmt.Errorf("scene %s does not support streaming", in.SceneCode)
		completion.ErrorMsg = sql.NullString{String: err.Error(), Valid: true}
		return err
	}

	// 4. 获取供应商
	l.Logger.Infof("Getting provider: %s", sceneConfig.ProviderCode)
	llmProvider := l.common.GetProviderManager().GetProvider(sceneConfig.ProviderCode)
	if llmProvider == nil {
		err := fmt.Errorf("provider %s not supported", sceneConfig.ProviderCode)
		completion.ErrorMsg = sql.NullString{String: err.Error(), Valid: true}
		return err
	}

	// 5. 构建请求
	l.Logger.Debug("Building LLM request")
	req := &provider.LLMRequest{
		Messages:    common.ConvertToProviderMessages(in.Messages),
		ModelCode:   sceneConfig.ModelCode,
		Temperature: sceneConfig.Temperature,
		MaxTokens:   sceneConfig.MaxTokens,
		Stream:      true,
		ExtraParams: in.ExtraParams,
		Config:      l.common.GetProviderConfig(sceneConfig.ProviderCode),
	}

	l.Logger.Infof("LLM request built - Model: %s, Temperature: %f, MaxTokens: %d, Messages: %d",
		req.ModelCode, req.Temperature, req.MaxTokens, len(req.Messages))

	// 6. 调用流式LLM
	l.Logger.Debug("Calling stream LLM")
	streamReader, err := llmProvider.StreamLLM(l.common.GetContext(), req)
	if err != nil {
		l.Logger.Errorf("Failed to call stream LLM: %v", err)
		completion.ErrorMsg = sql.NullString{String: fmt.Sprintf("failed to call stream LLM: %v", err), Valid: true}
		return fmt.Errorf("failed to call stream LLM: %w", err)
	}
	defer streamReader.Close()

	l.Logger.Info("LLM stream established, starting to process responses")

	responseCount := 0
	// 7. 处理流式响应
	var completionText strings.Builder
	var finalUsage *bs_llm.LLMUsage

	for {
		response, err := streamReader.Read()
		if err != nil {
			if err == io.EOF {
				l.Logger.Infof("LLM stream ended normally after %d responses", responseCount)
				break
			}
			l.Logger.Errorf("Failed to read stream response: %v", err)
			completion.ErrorMsg = sql.NullString{String: fmt.Sprintf("failed to read stream response: %v", err), Valid: true}
			return fmt.Errorf("failed to read stream response: %w", err)
		}

		responseCount++
		l.Logger.Infof("Received LLM response %d - Delta: %s, Finished: %v", responseCount, response.Delta(), response.Finished())

		// 累积完整的回答内容
		if response.Delta() != "" {
			completionText.WriteString(response.Delta())
		}

		// 构建gRPC流式响应
		streamResp := &bs_llm.StreamLLMResponse{
			Delta:        response.Delta(),
			ModelId:      sceneConfig.ModelCode,
			Finished:     response.Finished(),
			FinishReason: response.FinishReason(),
		}

		// 如果已完成，保存usage信息
		if response.Finished() && response.Usage() != nil {
			streamResp.Usage = response.Usage()
			finalUsage = response.Usage()
			l.Logger.Infof("Final usage - Prompt: %d, Completion: %d, Total: %d",
				response.Usage().PromptTokens, response.Usage().CompletionTokens, response.Usage().TotalTokens)
		}

		// 发送响应
		if err := stream.Send(streamResp); err != nil {
			l.Logger.Errorf("Failed to send stream response: %v", err)
			completion.ErrorMsg = sql.NullString{String: fmt.Sprintf("failed to send stream response: %v", err), Valid: true}
			return fmt.Errorf("failed to send stream response: %w", err)
		}

		// 如果已完成，退出循环
		if response.Finished() {
			l.Logger.Infof("Stream completed for scene: %s, model: %s, total responses: %d",
				sceneConfig.SceneCode, sceneConfig.ModelCode, responseCount)
			break
		}
	}

	// 8. 更新completion记录为成功状态
	completion.Completion = sql.NullString{String: completionText.String(), Valid: true}
	completion.Status = 1 // 成功

	// 设置 token 使用情况
	if finalUsage != nil {
		// 使用 LLM 返回的实际 usage 信息
		completion.InputTokens = finalUsage.PromptTokens
		completion.OutputTokens = finalUsage.CompletionTokens
		completion.TotalTokens = finalUsage.TotalTokens
		l.Logger.Infof("Using LLM provided usage - Input: %d, Output: %d, Total: %d",
			finalUsage.PromptTokens, finalUsage.CompletionTokens, finalUsage.TotalTokens)
	} else {
		// 如果没有 usage 信息，自己计算 token 数量
		inputTokens := l.common.EstimateTokens(completion.Prompt)
		outputTokens := l.common.EstimateTokens(completionText.String())
		totalTokens := inputTokens + outputTokens

		completion.InputTokens = inputTokens
		completion.OutputTokens = outputTokens
		completion.TotalTokens = totalTokens
		l.Logger.Infof("Estimated token usage - Input: %d, Output: %d, Total: %d",
			inputTokens, outputTokens, totalTokens)
	}

	l.Logger.Info("StreamLLM logic completed successfully")
	return nil
}
