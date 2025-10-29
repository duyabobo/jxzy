package common

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

// EmbeddingService 向量化服务
type EmbeddingService struct {
	logger logx.Logger
	apiKey string
}

// NewEmbeddingService 创建向量化服务实例
func NewEmbeddingService(apiKey string) *EmbeddingService {
	return &EmbeddingService{
		logger: logx.WithContext(nil),
		apiKey: apiKey,
	}
}

// GenerateEmbedding 生成文本的向量表示
func (e *EmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	return e.callBailianEmbeddingAPI(text)
}

// callBailianEmbeddingAPI 调用阿里云百炼 Embedding API
func (e *EmbeddingService) callBailianEmbeddingAPI(text string) ([]float32, error) {
	// 构建请求体 - 根据百炼API文档修正格式
	requestBody := map[string]interface{}{
		"model": "text-embedding-v4",
		"input": map[string]interface{}{
			"texts": []string{text},
		},
		"parameters": map[string]interface{}{
			"dimensions": consts.DashVectorDefaultDimension,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	apiKey := e.getBailianAPIKey()
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-DashScope-SSE", "disable")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
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

	// 返回第一个embedding向量
	embedding := response.Output.Embeddings[0].Embedding
	e.logger.Debugf("Generated embedding vector for text '%s': length=%d", text, len(embedding))
	return embedding, nil
}

// getBailianAPIKey 获取百炼 API Key
func (e *EmbeddingService) getBailianAPIKey() string {
	// 优先使用配置中的API Key
	if e.apiKey != "" {
		return e.apiKey
	}

	// 如果配置中没有，尝试从环境变量获取
	if apiKey := os.Getenv("BAILIAN_API_KEY"); apiKey != "" {
		return apiKey
	}

	// 如果都没有，返回默认值
	return "your-bailian-api-key"
}
