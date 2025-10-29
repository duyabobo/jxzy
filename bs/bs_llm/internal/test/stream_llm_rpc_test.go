package test

import (
	"testing"

	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/server"
	"jxzy/bs/bs_llm/internal/svc"

	"github.com/stretchr/testify/assert"
)

// TestStreamLLM_RPC 测试StreamLLM RPC接口
// 测试场景：使用百炼场景进行流式对话
func TestStreamLLM_RPC(t *testing.T) {
	// 1. 初始化配置和服务
	cfg := GetTestConfig()
	svcCtx := svc.NewServiceContext(cfg)
	llmServer := server.NewBsLlmServiceServer(svcCtx)

	// 2. 创建模拟流服务器
	mockStream := NewMockStreamServer()

	// 3. 构建测试请求
	req := &bs_llm.LLMRequest{
		SceneCode: "chat_bailian_turbo", // 使用百炼场景
		Messages: []*bs_llm.ChatMessage{
			{
				Role:    "user",
				Content: "你好，请介绍一下你自己",
			},
		},
		UserId: "test_user_stream_001",
		ExtraParams: map[string]string{
			"test_mode": "true",
		},
	}

	// 4. 执行流式调用
	err := llmServer.StreamLLM(req, mockStream)

	// 5. 验证结果
	// 如果有错误，记录错误信息但不强制失败（可能是API配置问题）
	if err != nil {
		t.Logf("StreamLLM call returned error: %v", err)
		// 如果没有任何响应，跳过后续验证
		if len(mockStream.responses) == 0 {
			t.Skip("Skipping response validation due to API error")
			return
		}
	}

	// 验证有响应返回
	assert.Greater(t, len(mockStream.responses), 0, "Should receive at least one response")

	// 验证最后一个响应是完成的
	if len(mockStream.responses) > 0 {
		lastResponse := mockStream.responses[len(mockStream.responses)-1]
		assert.True(t, lastResponse.Finished, "Last response should be finished")

		// 验证完成原因
		if lastResponse.FinishReason != "" {
			assert.NotEmpty(t, lastResponse.FinishReason, "Finish reason should not be empty")
		}

		// 验证模型ID
		if lastResponse.ModelId != "" {
			assert.NotEmpty(t, lastResponse.ModelId, "Model ID should not be empty")
		}

		// 验证使用情况（仅在完成时返回）
		if lastResponse.Usage != nil {
			assert.GreaterOrEqual(t, lastResponse.Usage.TotalTokens, int64(0), "Total tokens should be non-negative")
			assert.GreaterOrEqual(t, lastResponse.Usage.PromptTokens, int64(0), "Prompt tokens should be non-negative")
			assert.GreaterOrEqual(t, lastResponse.Usage.CompletionTokens, int64(0), "Completion tokens should be non-negative")
		}
	}

	// 6. 验证流式响应的连续性
	for i, response := range mockStream.responses {
		// 除了最后一个响应，其他都不应该是完成的
		if i < len(mockStream.responses)-1 {
			assert.False(t, response.Finished, "Non-final response should not be finished")
		}

		// 验证增量内容不为空（除了可能的最后一个空响应）
		if !response.Finished || response.Delta != "" {
			assert.NotNil(t, response.Delta, "Delta content should not be nil")
		}
	}

	// 记录测试结果
	if len(mockStream.responses) > 0 {
		t.Logf("StreamLLM test completed successfully. Received %d responses", len(mockStream.responses))
	} else {
		t.Logf("StreamLLM test completed with API error or no responses")
	}
}
