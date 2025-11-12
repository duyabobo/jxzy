package bailian

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	consts "jxzy/common/const"

	"github.com/zeromicro/go-zero/core/logx"
)

// Provider 阿里云百炼嵌入模型实现
type Provider struct {
	logger          logx.Logger
	apiKey          string
	modelCode       string
	vectorDimension int64
}

// NewBailianEmbeddingProvider 构造函数
func NewBailianEmbeddingProvider(apiKey string, modelCode string, vectorDimension int64) *Provider {
	if modelCode == "" {
		modelCode = "text-embedding-v4" // 默认模型
	}
	if vectorDimension == 0 {
		vectorDimension = consts.DashVectorDefaultDimension // 默认维度
	}
	return &Provider{
		logger:          logx.WithContext(nil),
		apiKey:          apiKey,
		modelCode:       modelCode,
		vectorDimension: vectorDimension,
	}
}

// GenerateEmbedding 生成文本的向量表示
func (p *Provider) GenerateEmbedding(text string) ([]float32, error) {
	return p.callBailianEmbeddingAPI(text)
}

func (p *Provider) callBailianEmbeddingAPI(text string) ([]float32, error) {
	requestBody := map[string]interface{}{
		"model": p.modelCode,
		"input": map[string]interface{}{
			"texts": []string{text},
		},
		"parameters": map[string]interface{}{
			"dimensions": p.vectorDimension,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	apiKey := p.getBailianAPIKey()
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-DashScope-SSE", "disable")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var response struct {
		Output struct {
			Embeddings []struct {
				Embedding []float32 `json:"embedding"`
			} `json:"embeddings"`
		} `json:"output"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	if len(response.Output.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings in response")
	}

	embedding := response.Output.Embeddings[0].Embedding
	p.logger.Debugf("Generated embedding vector for text '%s': length=%d", text, len(embedding))
	return embedding, nil
}

func (p *Provider) getBailianAPIKey() string {
	if p.apiKey != "" {
		return p.apiKey
	}
	if apiKey := os.Getenv("BAILIAN_API_KEY"); apiKey != "" {
		return apiKey
	}
	return "your-bailian-api-key"
}
