package doubao

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/provider"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	DefaultAPIEndpoint = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
	ProviderName       = "doubao"
)

// DoubaoProvider 豆包供应商实现
type DoubaoProvider struct {
	client *http.Client
	logger logx.Logger
}

// NewDoubaoProvider 创建豆包供应商
func NewDoubaoProvider() *DoubaoProvider {
	return &DoubaoProvider{
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		logger: logx.WithContext(context.Background()),
	}
}

// Name 供应商名称
func (p *DoubaoProvider) Name() string {
	return ProviderName
}

// CallLLM 非流式调用
func (p *DoubaoProvider) CallLLM(ctx context.Context, req *provider.LLMRequest) (*provider.LLMResponse, error) {
	apiReq := &ChatCompletionRequest{
		Model:       req.ModelCode,
		Messages:    convertMessages(req.Messages),
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      false,
	}

	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	endpoint := req.Config.APIEndpoint
	if endpoint == "" {
		endpoint = DefaultAPIEndpoint
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+req.Config.APIKey)
	for k, v := range req.Config.Headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := apiResp.Choices[0]
	return &provider.LLMResponse{
		Content:          choice.Message.Content,
		ModelCode:        req.ModelCode,
		PromptTokens:     apiResp.Usage.PromptTokens,
		CompletionTokens: apiResp.Usage.CompletionTokens,
		TotalTokens:      apiResp.Usage.TotalTokens,
		FinishReason:     choice.FinishReason,
	}, nil
}

// StreamLLM 流式调用
func (p *DoubaoProvider) StreamLLM(ctx context.Context, req *provider.LLMRequest) (provider.StreamReader, error) {
	apiReq := &ChatCompletionRequest{
		Model:       req.ModelCode,
		Messages:    convertMessages(req.Messages),
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      true,
	}

	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	endpoint := req.Config.APIEndpoint
	if endpoint == "" {
		endpoint = DefaultAPIEndpoint
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+req.Config.APIKey)
	httpReq.Header.Set("Accept", "text/event-stream")
	for k, v := range req.Config.Headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return NewDoubaoStreamReader(resp.Body, req.ModelCode), nil
}

// HealthCheck 健康检查
func (p *DoubaoProvider) HealthCheck(ctx context.Context) error {
	// 这里可以实现简单的健康检查逻辑
	// 比如调用一个简单的API来验证连接性
	return nil
}

// convertMessages 转换消息格式
func convertMessages(messages []*provider.ChatMessage) []ChatMessage {
	var result []ChatMessage
	for _, msg := range messages {
		result = append(result, ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return result
}

// DoubaoStreamReader 豆包流式读取器
type DoubaoStreamReader struct {
	reader    io.ReadCloser
	scanner   *bufio.Scanner
	modelCode string
	finished  bool
}

// NewDoubaoStreamReader 创建豆包流式读取器
func NewDoubaoStreamReader(reader io.ReadCloser, modelCode string) *DoubaoStreamReader {
	return &DoubaoStreamReader{
		reader:    reader,
		scanner:   bufio.NewScanner(reader),
		modelCode: modelCode,
		finished:  false,
	}
}

// Read 读取下一个响应
func (r *DoubaoStreamReader) Read() (provider.StreamResponse, error) {
	if r.finished {
		return nil, io.EOF
	}

	for r.scanner.Scan() {
		line := r.scanner.Text()

		// SSE格式: data: {...}
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// 检查是否是结束标志
			if data == "[DONE]" {
				r.finished = true
				return provider.NewStreamResponse("", true, nil, "stop", nil), nil
			}

			// 解析JSON数据
			var chunk ChatCompletionChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue // 跳过无法解析的行
			}

			// 处理流式响应
			if len(chunk.Choices) > 0 {
				choice := chunk.Choices[0]
				delta := ""
				if choice.Delta.Content != "" {
					delta = choice.Delta.Content
				}

				// 检查是否结束
				if choice.FinishReason != "" {
					r.finished = true
					var usage *bs_llm.LLMUsage
					if chunk.Usage != nil {
						usage = &bs_llm.LLMUsage{
							PromptTokens:     chunk.Usage.PromptTokens,
							CompletionTokens: chunk.Usage.CompletionTokens,
							TotalTokens:      chunk.Usage.TotalTokens,
						}
					}
					return provider.NewStreamResponse(delta, true, usage, choice.FinishReason, nil), nil
				}

				return provider.NewStreamResponse(delta, false, nil, "", nil), nil
			}
		}
	}

	// 检查扫描器错误
	if err := r.scanner.Err(); err != nil {
		return nil, err
	}

	r.finished = true
	return provider.NewStreamResponse("", true, nil, "stop", nil), nil
}

// Close 关闭流
func (r *DoubaoStreamReader) Close() error {
	r.finished = true
	return r.reader.Close()
}

// API请求和响应结构体
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int64         `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

// 流式响应结构体
type ChatCompletionChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChoiceChunk `json:"choices"`
	Usage   *Usage        `json:"usage,omitempty"`
}

type ChoiceChunk struct {
	Index        int         `json:"index"`
	Delta        ChatMessage `json:"delta"`
	FinishReason string      `json:"finish_reason"`
}
