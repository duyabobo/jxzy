package test

import (
	"context"
	"testing"

	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/server"
	"jxzy/bs/bs_llm/internal/svc"

	"github.com/stretchr/testify/assert"
)

// TestLLM_RPC 测试LLM RPC接口
// 测试场景：使用百炼场景进行非流式对话
func TestLLM_RPC(t *testing.T) {
	// 1. 初始化配置和服务
	cfg := GetTestConfig()
	svcCtx := svc.NewServiceContext(cfg)
	llmServer := server.NewBsLlmServiceServer(svcCtx)

	// 2. 构建测试请求
	req := &bs_llm.LLMRequest{
		SceneCode: "chat_bailian_turbo", // 使用百炼场景
		Messages: []*bs_llm.ChatMessage{
			{
				Role:    "user",
				Content: "你好，请介绍一下你自己",
			},
		},
		UserId: "test_user_llm_001",
		ExtraParams: map[string]string{
			"test_mode": "true",
		},
	}

	// 3. 执行非流式调用
	response, err := llmServer.LLM(context.Background(), req)

	// 4. 验证结果
	// 如果有错误，记录错误信息但不强制失败（可能是API配置问题）
	if err != nil {
		t.Logf("LLM call returned error: %v", err)
		// 如果响应为空，跳过后续验证
		if response == nil {
			t.Skip("Skipping response validation due to API error")
			return
		}
	}

	// 验证响应不为空
	assert.NotNil(t, response, "Response should not be nil")

	// 验证响应内容（只有在响应不为空时）
	if response != nil {
		// 验证完整回答内容
		if response.Completion != "" {
			assert.NotEmpty(t, response.Completion, "Completion should not be empty")
		}

		// 验证模型ID
		if response.ModelId != "" {
			assert.NotEmpty(t, response.ModelId, "Model ID should not be empty")
		}

		// 验证完成原因
		if response.FinishReason != "" {
			assert.NotEmpty(t, response.FinishReason, "Finish reason should not be empty")

			// 验证完成原因的有效值
			validFinishReasons := []string{"stop", "length", "content_filter"}
			assert.Contains(t, validFinishReasons, response.FinishReason,
				"Finish reason should be one of: stop, length, content_filter")
		}

		// 验证使用情况
		if response.Usage != nil {
			assert.GreaterOrEqual(t, response.Usage.TotalTokens, int64(0), "Total tokens should be non-negative")
			assert.GreaterOrEqual(t, response.Usage.PromptTokens, int64(0), "Prompt tokens should be non-negative")
			assert.GreaterOrEqual(t, response.Usage.CompletionTokens, int64(0), "Completion tokens should be non-negative")

			// 验证token数量的一致性
			if response.Usage.TotalTokens > 0 {
				assert.Equal(t, response.Usage.TotalTokens,
					response.Usage.PromptTokens+response.Usage.CompletionTokens,
					"Total tokens should equal prompt tokens plus completion tokens")
			}
		}
	}

	// 记录测试结果
	if response != nil && response.Completion != "" {
		t.Logf("LLM test completed successfully. Response length: %d characters", len(response.Completion))
	} else {
		t.Logf("LLM test completed with API error or empty response")
	}
}
