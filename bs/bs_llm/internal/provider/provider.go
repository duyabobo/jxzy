package provider

import (
	"context"
	"io"
	"jxzy/bs/bs_llm/bs_llm"
)

// StreamResponse 流式响应接口
type StreamResponse interface {
	// Delta 增量内容
	Delta() string
	// Finished 是否结束
	Finished() bool
	// Usage token使用情况
	Usage() *bs_llm.LLMUsage
	// FinishReason 结束原因
	FinishReason() string
	// Error 错误信息
	Error() error
}

// Provider LLM供应商接口
type Provider interface {
	// Name 供应商名称
	Name() string

	// CallLLM 非流式调用
	CallLLM(ctx context.Context, req *LLMRequest) (*LLMResponse, error)

	// StreamLLM 流式调用
	StreamLLM(ctx context.Context, req *LLMRequest) (StreamReader, error)

	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error
}

// StreamReader 流式读取器接口
type StreamReader interface {
	// Read 读取下一个响应
	Read() (StreamResponse, error)
	// Close 关闭流
	Close() error
}

// LLMRequest 标准化的LLM请求
type LLMRequest struct {
	Messages    []*ChatMessage    `json:"messages"`
	ModelCode   string            `json:"model_code"`
	Temperature float64           `json:"temperature"`
	MaxTokens   int64             `json:"max_tokens"`
	Stream      bool              `json:"stream"`
	ExtraParams map[string]string `json:"extra_params"`
	Config      *ProviderConfig   `json:"config"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMResponse 标准化的LLM响应
type LLMResponse struct {
	Content          string `json:"content"`
	ModelCode        string `json:"model_code"`
	PromptTokens     int64  `json:"prompt_tokens"`
	CompletionTokens int64  `json:"completion_tokens"`
	TotalTokens      int64  `json:"total_tokens"`
	FinishReason     string `json:"finish_reason"`
}

// ProviderConfig 供应商配置
type ProviderConfig struct {
	APIEndpoint   string            `json:"api_endpoint"`
	APIKey        string            `json:"api_key"`
	Headers       map[string]string `json:"headers"`
	DefaultParams map[string]string `json:"default_params"`
	Timeout       int32             `json:"timeout"`
	RetryCount    int32             `json:"retry_count"`
}

// Manager 供应商管理器
type Manager struct {
	providers map[string]Provider
}

// NewManager 创建供应商管理器
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
	}
}

// Register 注册供应商
func (m *Manager) Register(name string, provider Provider) {
	m.providers[name] = provider
}

// GetProvider 获取供应商
func (m *Manager) GetProvider(name string) Provider {
	return m.providers[name]
}

// ListProviders 列出所有供应商
func (m *Manager) ListProviders() []string {
	var names []string
	for name := range m.providers {
		names = append(names, name)
	}
	return names
}

// BaseStreamResponse 基础流式响应实现
type BaseStreamResponse struct {
	delta        string
	finished     bool
	usage        *bs_llm.LLMUsage
	finishReason string
	err          error
}

func NewStreamResponse(delta string, finished bool, usage *bs_llm.LLMUsage, finishReason string, err error) *BaseStreamResponse {
	return &BaseStreamResponse{
		delta:        delta,
		finished:     finished,
		usage:        usage,
		finishReason: finishReason,
		err:          err,
	}
}

func (r *BaseStreamResponse) Delta() string {
	return r.delta
}

func (r *BaseStreamResponse) Finished() bool {
	return r.finished
}

func (r *BaseStreamResponse) Usage() *bs_llm.LLMUsage {
	return r.usage
}

func (r *BaseStreamResponse) FinishReason() string {
	return r.finishReason
}

func (r *BaseStreamResponse) Error() error {
	return r.err
}

// BaseStreamReader 基础流式读取器
type BaseStreamReader struct {
	reader io.ReadCloser
	done   bool
}

func NewBaseStreamReader(reader io.ReadCloser) *BaseStreamReader {
	return &BaseStreamReader{
		reader: reader,
		done:   false,
	}
}

func (r *BaseStreamReader) Close() error {
	r.done = true
	if r.reader != nil {
		return r.reader.Close()
	}
	return nil
}
