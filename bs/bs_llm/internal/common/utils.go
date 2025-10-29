package common

import (
	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/provider"
)

// ConvertToProviderMessages 转换消息格式
func ConvertToProviderMessages(messages []*bs_llm.ChatMessage) []*provider.ChatMessage {
	var result []*provider.ChatMessage
	for _, msg := range messages {
		result = append(result, &provider.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return result
}
