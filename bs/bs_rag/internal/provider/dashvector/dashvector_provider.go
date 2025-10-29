package dashvector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"jxzy/bs/bs_rag/internal/provider/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// DashVectorConfig 阿里云 DashVector 配置
type DashVectorConfig struct {
	Endpoint string            `json:"endpoint"` // DashVector 服务端点
	APIKey   string            `json:"api_key"`  // API 密钥
	Region   string            `json:"region"`   // 地域
	Timeout  int               `json:"timeout"`  // 请求超时时间（秒）
	Headers  map[string]string `json:"headers"`  // 自定义请求头
}

// DashVectorProvider 阿里云 DashVector 向量数据库提供者
type DashVectorProvider struct {
	config     DashVectorConfig
	httpClient *http.Client
	baseURL    string
}

// API 请求和响应结构体
type apiResponse struct {
	RequestID string          `json:"request_id,omitempty"`
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"data,omitempty"`
	Output    json.RawMessage `json:"output,omitempty"`
}

// 错误响应中的操作结果
type operationResult struct {
	DocOp   string `json:"doc_op"`
	ID      string `json:"id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Collection 相关结构体
type collectionInfo struct {
	Name                string                  `json:"name"`
	Dimension           int                     `json:"dimension"`
	Dtype               string                  `json:"dtype"`
	Metric              string                  `json:"metric"`
	FieldsSchema        map[string]string       `json:"fields_schema"`
	Status              string                  `json:"status"`
	Partitions          map[string]string       `json:"partitions"`
	VectorsSchema       map[string]vectorSchema `json:"vectors_schema"`
	SparseVectorsSchema map[string]interface{}  `json:"sparse_vectors_schema"`
}

type vectorSchema struct {
	Dimension    int    `json:"dimension"`
	Dtype        string `json:"dtype"`
	Metric       string `json:"metric"`
	QuantizeType string `json:"quantize_type"`
}

type createCollectionRequest struct {
	Name        string `json:"name"`
	Dimension   int    `json:"dimension"`
	Metric      string `json:"metric"`
	Dtype       string `json:"dtype"`
	Description string `json:"description,omitempty"`
}

// Document 相关结构体
type docRequest struct {
	ID     string                 `json:"id"`
	Vector []float32              `json:"vector"`
	Fields map[string]interface{} `json:"fields"`
}

type insertDocsRequest struct {
	Docs []docRequest `json:"docs"`
}

// DashVector 插入操作直接返回操作结果数组，不需要单独的响应结构体

// Search 相关结构体
type searchRequest struct {
	Vector        []float32              `json:"vector"`
	TopK          int                    `json:"topk"`
	Filter        map[string]interface{} `json:"filter,omitempty"`
	IncludeVector bool                   `json:"include_vector,omitempty"`
	IncludeFields bool                   `json:"include_fields,omitempty"`
}

// DashVector 搜索操作直接返回搜索结果数组，不需要单独的响应结构体

type searchResult struct {
	ID     string                 `json:"id"`
	Vector []float32              `json:"vector,omitempty"`
	Score  float32                `json:"score"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

// Delete 相关结构体
type deleteDocsRequest struct {
	IDs []string `json:"ids"`
}

// DashVector 删除操作直接返回操作结果数组，不需要单独的响应结构体

// NewDashVectorProvider 创建新的 DashVector 提供者
func NewDashVectorProvider(config DashVectorConfig) *DashVectorProvider {
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	// 使用常量定义，不再从配置读取

	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	// 确保baseURL不以斜杠结尾
	baseURL := strings.TrimSuffix(config.Endpoint, "/")

	return &DashVectorProvider{
		config:     config,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// Search 执行向量搜索
func (p *DashVectorProvider) Search(ctx context.Context, collectionName string, queryVector []float32, topK int, minScore float32) ([]types.SearchResult, error) {
	logger := logx.WithContext(ctx)

	logger.Infof("Starting vector search - collection: %s, vector_dim: %d, topK: %d, minScore: %.3f",
		collectionName, len(queryVector), topK, minScore)

	reqBody := searchRequest{
		Vector:        queryVector,
		TopK:          topK,
		IncludeVector: false, // 不返回向量数据，减少响应大小
		IncludeFields: true,  // 包含字段数据
	}

	url := fmt.Sprintf("%s/v1/collections/%s/query", p.baseURL, collectionName)

	var searchResults []searchResult
	if err := p.makeRequest(ctx, "POST", url, reqBody, &searchResults); err != nil {
		logger.Errorf("Vector search failed - collection: %s, error: %v", collectionName, err)
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	logger.Infof("Vector search API response received - collection: %s, results_count: %d",
		collectionName, len(searchResults))

	results := make([]types.SearchResult, 0, len(searchResults))
	filteredCount := 0

	for _, result := range searchResults {
		// 过滤低于最小分数的结果
		if result.Score < minScore {
			filteredCount++
			continue
		}

		// 转换元数据
		metadata := make(map[string]string)
		var content string

		if result.Fields != nil {
			for k, v := range result.Fields {
				if k == "content" {
					if str, ok := v.(string); ok {
						content = str
					}
				} else {
					if str, ok := v.(string); ok {
						metadata[k] = str
					} else {
						metadata[k] = fmt.Sprintf("%v", v)
					}
				}
			}
		}

		results = append(results, types.SearchResult{
			ID:       result.ID,
			Vector:   result.Vector,
			Score:    result.Score,
			Metadata: metadata,
			Content:  content,
		})
	}

	logger.Infof("Vector search completed successfully - collection: %s, total_results: %d, filtered_results: %d, final_results: %d",
		collectionName, len(searchResults), filteredCount, len(results))

	return results, nil
}

// Insert 插入向量文档
func (p *DashVectorProvider) Insert(ctx context.Context, collectionName string, documents []types.Document) error {
	logger := logx.WithContext(ctx)

	if len(documents) == 0 {
		logger.Infof("Insert operation skipped - collection: %s, reason: no documents provided", collectionName)
		return nil
	}

	logger.Infof("Starting document insertion - collection: %s, document_count: %d",
		collectionName, len(documents))

	// 转换文档格式
	docs := make([]docRequest, len(documents))
	vectorDimensions := make(map[int]int) // 统计向量维度分布

	for i, doc := range documents {
		// 创建fields对象，包含metadata和content
		fields := make(map[string]interface{})

		// 添加metadata
		for k, v := range doc.Metadata {
			fields[k] = v
		}

		// 添加content（如果存在）
		if doc.Content != "" {
			fields["content"] = doc.Content
		}

		docs[i] = docRequest{
			ID:     doc.ID,
			Vector: doc.Vector,
			Fields: fields,
		}

		// 统计向量维度
		dim := len(doc.Vector)
		vectorDimensions[dim]++
	}

	// 记录向量维度统计
	for dim, count := range vectorDimensions {
		logger.Infof("Document vector dimensions - collection: %s, dimension: %d, count: %d",
			collectionName, dim, count)
	}

	reqBody := insertDocsRequest{
		Docs: docs,
	}

	url := fmt.Sprintf("%s/v1/collections/%s/docs", p.baseURL, collectionName)

	var insertResults []operationResult
	if err := p.makeRequest(ctx, "POST", url, reqBody, &insertResults); err != nil {
		logger.Errorf("Document insertion failed - collection: %s, document_count: %d, error: %v",
			collectionName, len(documents), err)
		return fmt.Errorf("insert request failed: %w", err)
	}

	// 统计插入结果
	successCount := 0
	var successIds []string
	var failedOps []string

	for _, result := range insertResults {
		if result.Code == 0 {
			successCount++
			successIds = append(successIds, result.ID)
		} else {
			failedOps = append(failedOps, fmt.Sprintf("id=%s, error=%s", result.ID, result.Message))
		}
	}

	if len(failedOps) > 0 {
		logger.Errorf("Document insertion partially failed - collection: %s, requested_count: %d, success_count: %d, failed_operations: %v",
			collectionName, len(documents), successCount, failedOps)
	}

	logger.Infof("Document insertion completed - collection: %s, requested_count: %d, success_count: %d, success_ids: %v",
		collectionName, len(documents), successCount, successIds)

	return nil
}

// Delete 删除向量文档
func (p *DashVectorProvider) Delete(ctx context.Context, collectionName string, documentIDs []string) error {
	logger := logx.WithContext(ctx)

	if len(documentIDs) == 0 {
		logger.Infof("Delete operation skipped - collection: %s, reason: no document IDs provided", collectionName)
		return nil
	}

	logger.Infof("Starting document deletion - collection: %s, document_count: %d, ids: %v",
		collectionName, len(documentIDs), documentIDs)

	reqBody := deleteDocsRequest{
		IDs: documentIDs,
	}

	url := fmt.Sprintf("%s/v1/collections/%s/docs", p.baseURL, collectionName)

	var deleteResults []operationResult
	if err := p.makeRequest(ctx, "DELETE", url, reqBody, &deleteResults); err != nil {
		logger.Errorf("Document deletion failed - collection: %s, document_count: %d, ids: %v, error: %v",
			collectionName, len(documentIDs), documentIDs, err)
		return fmt.Errorf("delete request failed: %w", err)
	}

	// 统计删除结果
	successCount := 0
	var successIds []string
	var failedOps []string

	for _, result := range deleteResults {
		if result.Code == 0 {
			successCount++
			successIds = append(successIds, result.ID)
		} else {
			failedOps = append(failedOps, fmt.Sprintf("id=%s, error=%s", result.ID, result.Message))
		}
	}

	if len(failedOps) > 0 {
		logger.Errorf("Document deletion partially failed - collection: %s, requested_count: %d, success_count: %d, failed_operations: %v",
			collectionName, len(documentIDs), successCount, failedOps)
	}

	logger.Infof("Document deletion completed - collection: %s, requested_count: %d, success_count: %d, success_ids: %v",
		collectionName, len(documentIDs), successCount, successIds)

	return nil
}

// GetCollectionInfo 获取集合信息
func (p *DashVectorProvider) GetCollectionInfo(ctx context.Context, collectionName string) (*types.CollectionInfo, error) {
	logger := logx.WithContext(ctx)

	logger.Infof("Starting get collection info - collection: %s", collectionName)

	url := fmt.Sprintf("%s/v1/collections/%s", p.baseURL, collectionName)

	var collection collectionInfo
	if err := p.makeRequest(ctx, "GET", url, nil, &collection); err != nil {
		// 如果集合不存在，返回不存在的状态
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			logger.Infof("Collection does not exist - collection: %s", collectionName)
			return &types.CollectionInfo{
				Name:   collectionName,
				Exists: false,
			}, nil
		}
		logger.Errorf("Get collection info failed - collection: %s, error: %v", collectionName, err)
		return nil, fmt.Errorf("describe collection request failed: %w", err)
	}

	// 计算文档数量（从partitions状态推断）
	docCount := 0
	// DashVector API中没有直接的doc_count字段，这里设置为0或者从其他地方获取

	collectionInfo := &types.CollectionInfo{
		Name:          collection.Name,
		Dimension:     collection.Dimension,
		IndexType:     collection.Metric,
		DocumentCount: docCount,
		Metadata: map[string]string{
			"provider": "dashvector",
			"dtype":    collection.Dtype,
			"metric":   collection.Metric,
			"status":   collection.Status,
		},
		Exists: collection.Status == "SERVING",
	}

	logger.Infof("Get collection info completed successfully - collection: %s, dimension: %d, status: %s, dtype: %s, metric: %s, partitions: %v",
		collectionName, collection.Dimension, collection.Status, collection.Dtype, collection.Metric, collection.Partitions)

	return collectionInfo, nil
}

// CreateCollection 创建集合
func (p *DashVectorProvider) CreateCollection(ctx context.Context, collectionName string, dimension int, indexType string) error {
	logger := logx.WithContext(ctx)

	logger.Infof("Starting create collection - collection: %s, dimension: %d, index_type: %s",
		collectionName, dimension, indexType)

	reqBody := createCollectionRequest{
		Name:        collectionName,
		Dimension:   dimension,
		Metric:      indexType,
		Dtype:       "FLOAT", // 默认使用FLOAT类型
		Description: fmt.Sprintf("Collection created by bs_rag service at %s", time.Now().Format(time.RFC3339)),
	}

	url := fmt.Sprintf("%s/v1/collections", p.baseURL)

	if err := p.makeRequest(ctx, "POST", url, reqBody, nil); err != nil {
		logger.Errorf("Create collection failed - collection: %s, dimension: %d, index_type: %s, error: %v",
			collectionName, dimension, indexType, err)
		return fmt.Errorf("create collection request failed: %w", err)
	}

	logger.Infof("Create collection completed successfully - collection: %s, dimension: %d, index_type: %s",
		collectionName, dimension, indexType)

	return nil
}

// DeleteCollection 删除集合
func (p *DashVectorProvider) DeleteCollection(ctx context.Context, collectionName string) error {
	logger := logx.WithContext(ctx)

	logger.Infof("Starting delete collection - collection: %s", collectionName)

	url := fmt.Sprintf("%s/v1/collections/%s", p.baseURL, collectionName)

	if err := p.makeRequest(ctx, "DELETE", url, nil, nil); err != nil {
		logger.Errorf("Delete collection failed - collection: %s, error: %v", collectionName, err)
		return fmt.Errorf("delete collection request failed: %w", err)
	}

	logger.Infof("Delete collection completed successfully - collection: %s", collectionName)

	return nil
}

// ListCollections 列出所有集合
func (p *DashVectorProvider) ListCollections(ctx context.Context) ([]string, error) {
	logger := logx.WithContext(ctx)

	logger.Infof("Starting list collections")

	url := fmt.Sprintf("%s/v1/collections", p.baseURL)

	var collections []collectionInfo
	if err := p.makeRequest(ctx, "GET", url, nil, &collections); err != nil {
		logger.Errorf("List collections failed - error: %v", err)
		return nil, fmt.Errorf("list collections request failed: %w", err)
	}

	logger.Infof("List collections API response received - total_collections: %d", len(collections))

	result := make([]string, 0, len(collections))
	servingCount := 0

	for _, collection := range collections {
		if collection.Status == "SERVING" {
			result = append(result, collection.Name)
			servingCount++
		}
		logger.Infof("Collection found - name: %s, status: %s, dimension: %d, dtype: %s, metric: %s",
			collection.Name, collection.Status, collection.Dimension, collection.Dtype, collection.Metric)
	}

	logger.Infof("List collections completed successfully - total_collections: %d, serving_collections: %d",
		len(collections), servingCount)

	return result, nil
}

// Close 关闭连接
func (p *DashVectorProvider) Close() error {
	// HTTP 客户端不需要显式关闭
	return nil
}

// makeRequest 发送 HTTP 请求到 DashVector API
func (p *DashVectorProvider) makeRequest(ctx context.Context, method, url string, reqBody interface{}, respBody interface{}) error {
	logger := logx.WithContext(ctx)

	startTime := time.Now()

	// 记录请求开始信息
	logger.Infof("Starting HTTP request - method: %s, url: %s", method, url)

	var req *http.Request
	var err error
	var reqBodySize int

	// 创建请求体
	if reqBody != nil && (method == "POST" || method == "PUT" || method == "DELETE") {
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			logger.Errorf("Failed to marshal request body - method: %s, url: %s, error: %v", method, url, err)
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBodySize = len(jsonData)
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Errorf("Failed to create request - method: %s, url: %s, error: %v", method, url, err)
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		logger.Infof("Request body prepared - method: %s, url: %s, body_size: %d bytes", method, url, reqBodySize)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			logger.Errorf("Failed to create request - method: %s, url: %s, error: %v", method, url, err)
			return fmt.Errorf("failed to create request: %w", err)
		}
		logger.Infof("Request prepared without body - method: %s, url: %s", method, url)
	}

	// 设置认证头 - 尝试多种可能的认证方式
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	req.Header.Set("dashvector-auth-token", p.config.APIKey)
	req.Header.Set("X-API-Key", p.config.APIKey)

	// 设置地域头
	if p.config.Region != "" {
		req.Header.Set("X-Region", p.config.Region)
	}

	// 设置自定义请求头
	for k, v := range p.config.Headers {
		req.Header.Set(k, v)
	}

	logger.Infof("Request headers configured - method: %s, url: %s, region: %s", method, url, p.config.Region)

	// 发送请求
	resp, err := p.httpClient.Do(req)
	if err != nil {
		duration := time.Since(startTime)
		logger.Errorf("HTTP request failed - method: %s, url: %s, duration: %v, error: %v", method, url, duration, err)
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)
	logger.Infof("HTTP response received - method: %s, url: %s, status: %d, duration: %v",
		method, url, resp.StatusCode, duration)

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(resp.Body)
		// 如果响应是HTML，可能是错误页面或重定向
		errorBodyStr := string(errorBody)
		if len(errorBodyStr) > 500 {
			errorBodyStr = errorBodyStr[:500] + "..."
		}
		logger.Errorf("HTTP request failed with error status - method: %s, url: %s, status: %d, response: %s",
			method, url, resp.StatusCode, errorBodyStr)
		return fmt.Errorf("request failed with status %d, URL: %s, response: %s", resp.StatusCode, url, errorBodyStr)
	}

	// 解析响应体
	if respBody != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("Failed to read response body - method: %s, url: %s, error: %v", method, url, err)
			return fmt.Errorf("failed to read response body: %w", err)
		}

		bodyStr := string(body)
		respBodySize := len(body)

		logger.Infof("Response body received - method: %s, url: %s, body_size: %d bytes", method, url, respBodySize)

		// 检查是否返回了HTML而不是JSON
		if strings.HasPrefix(strings.TrimSpace(bodyStr), "<") {
			// 返回了HTML，可能是错误页面
			if len(bodyStr) > 500 {
				bodyStr = bodyStr[:500] + "..."
			}
			logger.Errorf("Received HTML response instead of JSON - method: %s, url: %s, response: %s", method, url, bodyStr)
			return fmt.Errorf("received HTML response instead of JSON, URL: %s, response: %s", url, bodyStr)
		}

		// 先尝试解析为通用API响应格式
		var apiResp apiResponse
		if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Code != 0 {
			// 如果有详细的操作结果，解析并显示
			if apiResp.Output != nil {
				var operations []operationResult
				if err := json.Unmarshal(apiResp.Output, &operations); err == nil && len(operations) > 0 {
					op := operations[0]
					logger.Errorf("API returned error with operation details - method: %s, url: %s, code: %d, message: %s, operation: %s, document: %s, detail: %s",
						method, url, apiResp.Code, apiResp.Message, op.DocOp, op.ID, op.Message)
					return fmt.Errorf("API error: %s (code: %d), operation: %s, document: %s, detail: %s",
						apiResp.Message, apiResp.Code, op.DocOp, op.ID, op.Message)
				}
			}
			logger.Errorf("API returned error - method: %s, url: %s, code: %d, message: %s",
				method, url, apiResp.Code, apiResp.Message)
			return fmt.Errorf("API error: %s (code: %d)", apiResp.Message, apiResp.Code)
		}

		// 优先尝试解析Output字段（DashVector标准响应格式）
		if apiResp.Output != nil {
			if err := json.Unmarshal(apiResp.Output, respBody); err != nil {
				logger.Errorf("Failed to parse response output field - method: %s, url: %s, error: %v", method, url, err)
				return fmt.Errorf("failed to parse response output: %w, response: %s", err, bodyStr)
			}
			logger.Infof("Response output field parsed successfully - method: %s, url: %s", method, url)
		} else if apiResp.Data != nil {
			// 如果有Data字段，解析Data部分
			if err := json.Unmarshal(apiResp.Data, respBody); err != nil {
				logger.Errorf("Failed to parse response data field - method: %s, url: %s, error: %v", method, url, err)
				return fmt.Errorf("failed to parse response data: %w, response: %s", err, bodyStr)
			}
			logger.Infof("Response data field parsed successfully - method: %s, url: %s", method, url)
		} else {
			// 直接解析整个响应体
			if err := json.Unmarshal(body, respBody); err != nil {
				logger.Errorf("Failed to parse response body - method: %s, url: %s, error: %v", method, url, err)
				return fmt.Errorf("failed to parse response: %w, response: %s", err, bodyStr)
			}
			logger.Infof("Response body parsed successfully - method: %s, url: %s", method, url)
		}
	}

	logger.Infof("HTTP request completed successfully - method: %s, url: %s, total_duration: %v",
		method, url, time.Since(startTime))

	return nil
}
