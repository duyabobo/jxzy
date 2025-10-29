package bailian

import (
	"encoding/json"
	"testing"
)

func TestBailianStreamChunkParsing(t *testing.T) {
	// 测试百炼流式响应的JSON解析
	testData := `{"output":{"finish_reason":"null","text":"Hello"},"usage":{"total_tokens":14,"output_tokens":1,"input_tokens":13},"request_id":"test-123"}`

	var chunk BailianStreamChunk
	err := json.Unmarshal([]byte(testData), &chunk)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if chunk.RequestID != "test-123" {
		t.Errorf("Expected request_id 'test-123', got '%s'", chunk.RequestID)
	}

	if chunk.Output.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", chunk.Output.Text)
	}

	if chunk.Output.FinishReason != "null" {
		t.Errorf("Expected finish_reason 'null', got '%s'", chunk.Output.FinishReason)
	}

	if chunk.Usage == nil {
		t.Error("Expected usage to be not nil")
	} else {
		if chunk.Usage.InputTokens != 13 {
			t.Errorf("Expected input_tokens 13, got %d", chunk.Usage.InputTokens)
		}
		if chunk.Usage.OutputTokens != 1 {
			t.Errorf("Expected output_tokens 1, got %d", chunk.Usage.OutputTokens)
		}
		if chunk.Usage.InputTokens+chunk.Usage.OutputTokens != 14 {
			t.Errorf("Expected total_tokens 14, got %d", chunk.Usage.InputTokens+chunk.Usage.OutputTokens)
		}
	}

	t.Logf("Parsed chunk: %+v", chunk)
}

func TestBailianStreamChunkParsingWithStop(t *testing.T) {
	// 测试结束状态的响应
	testData := `{"output":{"finish_reason":"stop","text":"Hello! How can I assist you today? 😊"},"usage":{"total_tokens":24,"output_tokens":11,"input_tokens":13},"request_id":"test-456"}`

	var chunk BailianStreamChunk
	err := json.Unmarshal([]byte(testData), &chunk)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if chunk.Output.FinishReason != "stop" {
		t.Errorf("Expected finish_reason 'stop', got '%s'", chunk.Output.FinishReason)
	}

	if chunk.Output.Text != "Hello! How can I assist you today? 😊" {
		t.Errorf("Expected text 'Hello! How can I assist you today? 😊', got '%s'", chunk.Output.Text)
	}

	t.Logf("Parsed chunk with stop: %+v", chunk)
}
