package bailian

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
	DefaultAPIEndpoint = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	ProviderName       = "bailian"
)

// BailianProvider 百炼供应商实现
type BailianProvider struct {
	client *http.Client
	logger logx.Logger
}

// NewBailianProvider 创建百炼供应商
func NewBailianProvider() *BailianProvider {
	return &BailianProvider{
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		logger: logx.WithContext(context.Background()),
	}
}

// Name 供应商名称
func (p *BailianProvider) Name() string {
	return ProviderName
}

// createHTTPClient 根据配置创建HTTP客户端
func (p *BailianProvider) createHTTPClient(config *provider.ProviderConfig) *http.Client {
	timeout := 120 * time.Second
	if config != nil && config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	return &http.Client{
		Timeout: timeout,
	}
}

// CallLLM 非流式调用
func (p *BailianProvider) CallLLM(ctx context.Context, req *provider.LLMRequest) (*provider.LLMResponse, error) {
	// 构建百炼API请求
	apiReq := &BailianRequest{
		Model: req.ModelCode,
		Input: &BailianInput{
			Messages: convertMessages(req.Messages),
		},
		Parameters: &BailianParameters{
			Temperature: req.Temperature,
			MaxTokens:   req.MaxTokens,
			TopP:        0.8,
			TopK:        50,
		},
	}

	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	endpoint := req.Config.APIEndpoint
	if endpoint == "" {
		endpoint = DefaultAPIEndpoint
	}

	// 添加调试日志
	p.logger.Infof("Bailian API Request - Endpoint: %s, Model: %s", endpoint, req.ModelCode)
	p.logger.Infof("Bailian API Request Body: %s", string(reqBody))

	// 创建HTTP客户端
	client := p.createHTTPClient(req.Config)

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

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	p.logger.Infof("Bailian API Response Status: %d", resp.StatusCode)
	p.logger.Infof("Bailian API Response Body: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp BailianResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	// 检查输出内容是否为空
	if apiResp.Output.Text == "" {
		return nil, fmt.Errorf("no text content in response")
	}

	p.logger.Infof("Bailian API Output: finish_reason=%s, text_length=%d",
		apiResp.Output.FinishReason, len(apiResp.Output.Text))

	return &provider.LLMResponse{
		Content:          apiResp.Output.Text,
		ModelCode:        req.ModelCode,
		PromptTokens:     apiResp.Usage.InputTokens,
		CompletionTokens: apiResp.Usage.OutputTokens,
		TotalTokens:      apiResp.Usage.TotalTokens,
		FinishReason:     apiResp.Output.FinishReason,
	}, nil
}

// StreamLLM 流式调用
func (p *BailianProvider) StreamLLM(ctx context.Context, req *provider.LLMRequest) (provider.StreamReader, error) {
	// 构建百炼API请求
	apiReq := &BailianRequest{
		Model: req.ModelCode,
		Input: &BailianInput{
			Messages: convertMessages(req.Messages),
		},
		Parameters: &BailianParameters{
			Temperature: req.Temperature,
			MaxTokens:   req.MaxTokens,
			TopP:        0.8,
			TopK:        50,
			Stream:      true,
		},
	}

	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	endpoint := req.Config.APIEndpoint
	if endpoint == "" {
		endpoint = DefaultAPIEndpoint
	}

	// 添加调试日志
	p.logger.Infof("Bailian Stream API Request - Endpoint: %s, Model: %s", endpoint, req.ModelCode)
	p.logger.Infof("Bailian Stream API Request Body: %s", string(reqBody))

	// 创建HTTP客户端
	client := p.createHTTPClient(req.Config)

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

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return NewBailianStreamReader(resp.Body, req.ModelCode, p.logger), nil
}

// HealthCheck 健康检查
func (p *BailianProvider) HealthCheck(ctx context.Context) error {
	// 简单的健康检查，可以调用一个简单的API来验证连接性
	return nil
}

// convertMessages 转换消息格式
func convertMessages(messages []*provider.ChatMessage) []BailianMessage {
	var result []BailianMessage
	for _, msg := range messages {
		result = append(result, BailianMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return result
}

// BailianStreamReader 百炼流式读取器
type BailianStreamReader struct {
	reader       io.ReadCloser
	scanner      *bufio.Scanner
	modelCode    string
	finished     bool
	logger       logx.Logger
	previousText string // 用于计算增量
}

// NewBailianStreamReader 创建百炼流式读取器
func NewBailianStreamReader(reader io.ReadCloser, modelCode string, logger logx.Logger) *BailianStreamReader {
	scanner := bufio.NewScanner(reader)
	// 设置扫描器的缓冲区大小，确保能处理长行
	scanner.Buffer(make([]byte, 64*1024), 64*1024)

	return &BailianStreamReader{
		reader:       reader,
		scanner:      scanner,
		modelCode:    modelCode,
		finished:     false,
		logger:       logger,
		previousText: "",
	}
}

// Read 读取下一个响应
func (r *BailianStreamReader) Read() (provider.StreamResponse, error) {
	if r.finished {
		return nil, io.EOF
	}

	for r.scanner.Scan() {
		line := r.scanner.Text()

		// 跳过空行
		if strings.TrimSpace(line) == "" {
			continue
		}

		// SSE格式: data: {...}
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")

			// 检查是否是结束标志
			if data == "[DONE]" {
				r.finished = true
				r.logger.Infof("Bailian Stream: Received [DONE]")
				return provider.NewStreamResponse("", true, nil, "stop", nil), nil
			}

			// 解析JSON数据
			var chunk BailianStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				r.logger.Errorf("Bailian Stream: Failed to parse JSON: %v, data: %s", err, data)
				continue // 跳过无法解析的行
			}

			// 处理流式响应
			output := chunk.Output

			// 计算增量内容（百炼返回的是完整文本，需要计算增量）
			delta := ""
			if output.Text != "" {
				// 计算增量：当前文本减去之前的文本
				if len(output.Text) > len(r.previousText) {
					delta = output.Text[len(r.previousText):]
				}
				r.previousText = output.Text
			}

			// 检查是否结束
			if output.FinishReason != "" && output.FinishReason != "null" {
				r.finished = true
				var usage *bs_llm.LLMUsage
				if chunk.Usage != nil {
					usage = &bs_llm.LLMUsage{
						PromptTokens:     chunk.Usage.InputTokens,
						CompletionTokens: chunk.Usage.OutputTokens,
						TotalTokens:      chunk.Usage.InputTokens + chunk.Usage.OutputTokens,
					}
				}
				return provider.NewStreamResponse(delta, true, usage, output.FinishReason, nil), nil
			}

			// 只有在有增量内容时才返回响应
			if delta != "" {
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
func (r *BailianStreamReader) Close() error {
	r.finished = true
	return r.reader.Close()
}

// 百炼API请求和响应结构体
type BailianRequest struct {
	Model      string             `json:"model"`
	Input      *BailianInput      `json:"input"`
	Parameters *BailianParameters `json:"parameters"`
}

type BailianInput struct {
	Messages []BailianMessage `json:"messages"`
}

type BailianParameters struct {
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int64   `json:"max_tokens,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	Seed        int64   `json:"seed,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
}

type BailianMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BailianResponse struct {
	RequestID string        `json:"request_id"`
	Output    BailianOutput `json:"output"`
	Usage     BailianUsage  `json:"usage"`
}

type BailianOutput struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

// 保留旧的结构体定义以备兼容性使用（如果需要）
type BailianChoice struct {
	Message      BailianMessage `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

type BailianUsage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	TotalTokens  int64 `json:"total_tokens"`
}

// 流式响应结构体
type BailianStreamChunk struct {
	RequestID string              `json:"request_id"`
	Output    BailianStreamOutput `json:"output"`
	Usage     *BailianUsage       `json:"usage,omitempty"`
}

type BailianStreamOutput struct {
	FinishReason string `json:"finish_reason"`
	Text         string `json:"text"`
}
