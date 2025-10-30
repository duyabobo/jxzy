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

    "jxzy/bs/bs_rag/internal/provider/vectorstore/types"

    "github.com/zeromicro/go-zero/core/logx"
)

type DashVectorConfig struct {
    Endpoint string            `json:"endpoint"`
    APIKey   string            `json:"api_key"`
    Region   string            `json:"region"`
    Timeout  int               `json:"timeout"`
    Headers  map[string]string `json:"headers"`
}

type DashVectorProvider struct {
    config     DashVectorConfig
    httpClient *http.Client
    baseURL    string
}

type apiResponse struct {
    RequestID string          `json:"request_id,omitempty"`
    Code      int             `json:"code"`
    Message   string          `json:"message"`
    Data      json.RawMessage `json:"data,omitempty"`
    Output    json.RawMessage `json:"output,omitempty"`
}

type operationResult struct {
    DocOp   string `json:"doc_op"`
    ID      string `json:"id"`
    Code    int    `json:"code"`
    Message string `json:"message"`
}

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

type docRequest struct {
    ID     string                 `json:"id"`
    Vector []float32              `json:"vector"`
    Fields map[string]interface{} `json:"fields"`
}

type insertDocsRequest struct {
    Docs []docRequest `json:"docs"`
}

type searchRequest struct {
    Vector        []float32              `json:"vector"`
    TopK          int                    `json:"topk"`
    Filter        map[string]interface{} `json:"filter,omitempty"`
    IncludeVector bool                   `json:"include_vector,omitempty"`
    IncludeFields bool                   `json:"include_fields,omitempty"`
}

type searchResult struct {
    ID     string                 `json:"id"`
    Vector []float32              `json:"vector,omitempty"`
    Score  float32                `json:"score"`
    Fields map[string]interface{} `json:"fields,omitempty"`
}

type deleteDocsRequest struct {
    IDs []string `json:"ids"`
}

func NewDashVectorProvider(config DashVectorConfig) *DashVectorProvider {
    if config.Timeout == 0 {
        config.Timeout = 30
    }
    httpClient := &http.Client{Timeout: time.Duration(config.Timeout) * time.Second}
    baseURL := strings.TrimSuffix(config.Endpoint, "/")
    return &DashVectorProvider{config: config, httpClient: httpClient, baseURL: baseURL}
}

func (p *DashVectorProvider) Search(ctx context.Context, collectionName string, queryVector []float32, topK int, minScore float32) ([]types.SearchResult, error) {
    logger := logx.WithContext(ctx)
    logger.Infof("Starting vector search - collection: %s, vector_dim: %d, topK: %d, minScore: %.3f", collectionName, len(queryVector), topK, minScore)

    reqBody := searchRequest{Vector: queryVector, TopK: topK, IncludeVector: false, IncludeFields: true}
    url := fmt.Sprintf("%s/v1/collections/%s/query", p.baseURL, collectionName)
    var searchResults []searchResult
    if err := p.makeRequest(ctx, "POST", url, reqBody, &searchResults); err != nil {
        logger.Errorf("Vector search failed - collection: %s, error: %v", collectionName, err)
        return nil, fmt.Errorf("search request failed: %w", err)
    }

    results := make([]types.SearchResult, 0, len(searchResults))
    for _, r := range searchResults {
        if r.Score < minScore { continue }
        metadata := make(map[string]string)
        var content string
        if r.Fields != nil {
            for k, v := range r.Fields {
                if k == "content" { if s, ok := v.(string); ok { content = s; continue } }
                if s, ok := v.(string); ok { metadata[k] = s } else { metadata[k] = fmt.Sprintf("%v", v) }
            }
        }
        results = append(results, types.SearchResult{ID: r.ID, Vector: r.Vector, Score: r.Score, Metadata: metadata, Content: content})
    }
    return results, nil
}

func (p *DashVectorProvider) Insert(ctx context.Context, collectionName string, documents []types.Document) error {
    if len(documents) == 0 { return nil }
    docs := make([]docRequest, len(documents))
    for i, d := range documents {
        fields := make(map[string]interface{})
        for k, v := range d.Metadata { fields[k] = v }
        if d.Content != "" { fields["content"] = d.Content }
        docs[i] = docRequest{ID: d.ID, Vector: d.Vector, Fields: fields}
    }
    reqBody := insertDocsRequest{Docs: docs}
    url := fmt.Sprintf("%s/v1/collections/%s/docs", p.baseURL, collectionName)
    var insertResults []operationResult
    if err := p.makeRequest(ctx, "POST", url, reqBody, &insertResults); err != nil { return fmt.Errorf("insert request failed: %w", err) }
    return nil
}

func (p *DashVectorProvider) Delete(ctx context.Context, collectionName string, documentIDs []string) error {
    if len(documentIDs) == 0 { return nil }
    reqBody := deleteDocsRequest{IDs: documentIDs}
    url := fmt.Sprintf("%s/v1/collections/%s/docs", p.baseURL, collectionName)
    var deleteResults []operationResult
    if err := p.makeRequest(ctx, "DELETE", url, reqBody, &deleteResults); err != nil { return fmt.Errorf("delete request failed: %w", err) }
    return nil
}

func (p *DashVectorProvider) GetCollectionInfo(ctx context.Context, collectionName string) (*types.CollectionInfo, error) {
    logger := logx.WithContext(ctx)
    url := fmt.Sprintf("%s/v1/collections/%s", p.baseURL, collectionName)
    var c collectionInfo
    if err := p.makeRequest(ctx, "GET", url, nil, &c); err != nil {
        if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
            return &types.CollectionInfo{Name: collectionName, Exists: false}, nil
        }
        logger.Errorf("Get collection info failed - collection: %s, error: %v", collectionName, err)
        return nil, fmt.Errorf("describe collection request failed: %w", err)
    }
    info := &types.CollectionInfo{Name: c.Name, Dimension: c.Dimension, IndexType: c.Metric, DocumentCount: 0, Metadata: map[string]string{"provider": "dashvector", "dtype": c.Dtype, "metric": c.Metric, "status": c.Status}, Exists: c.Status == "SERVING"}
    return info, nil
}

func (p *DashVectorProvider) CreateCollection(ctx context.Context, collectionName string, dimension int, indexType string) error {
    reqBody := createCollectionRequest{Name: collectionName, Dimension: dimension, Metric: indexType, Dtype: "FLOAT", Description: fmt.Sprintf("Collection created by bs_rag service at %s", time.Now().Format(time.RFC3339))}
    url := fmt.Sprintf("%s/v1/collections", p.baseURL)
    if err := p.makeRequest(ctx, "POST", url, reqBody, nil); err != nil { return fmt.Errorf("create collection request failed: %w", err) }
    return nil
}

func (p *DashVectorProvider) DeleteCollection(ctx context.Context, collectionName string) error {
    url := fmt.Sprintf("%s/v1/collections/%s", p.baseURL, collectionName)
    if err := p.makeRequest(ctx, "DELETE", url, nil, nil); err != nil { return fmt.Errorf("delete collection request failed: %w", err) }
    return nil
}

func (p *DashVectorProvider) ListCollections(ctx context.Context) ([]string, error) {
    url := fmt.Sprintf("%s/v1/collections", p.baseURL)
    var collections []collectionInfo
    if err := p.makeRequest(ctx, "GET", url, nil, &collections); err != nil { return nil, fmt.Errorf("list collections request failed: %w", err) }
    result := make([]string, 0, len(collections))
    for _, c := range collections { if c.Status == "SERVING" { result = append(result, c.Name) } }
    return result, nil
}

func (p *DashVectorProvider) Close() error { return nil }

func (p *DashVectorProvider) makeRequest(ctx context.Context, method, url string, reqBody interface{}, respBody interface{}) error {
    logger := logx.WithContext(ctx)
    var req *http.Request
    var err error
    if reqBody != nil && (method == "POST" || method == "PUT" || method == "DELETE") {
        jsonData, err := json.Marshal(reqBody)
        if err != nil { return fmt.Errorf("failed to marshal request body: %w", err) }
        req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
        if err != nil { return fmt.Errorf("failed to create request: %w", err) }
        req.Header.Set("Content-Type", "application/json")
    } else {
        req, err = http.NewRequestWithContext(ctx, method, url, nil)
        if err != nil { return fmt.Errorf("failed to create request: %w", err) }
    }
    req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
    req.Header.Set("dashvector-auth-token", p.config.APIKey)
    req.Header.Set("X-API-Key", p.config.APIKey)
    if p.config.Region != "" { req.Header.Set("X-Region", p.config.Region) }
    for k, v := range p.config.Headers { req.Header.Set(k, v) }
    resp, err := p.httpClient.Do(req)
    if err != nil { return fmt.Errorf("request failed: %w", err) }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        b, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("request failed with status %d, URL: %s, response: %s", resp.StatusCode, url, string(b))
    }
    if respBody != nil {
        body, err := io.ReadAll(resp.Body)
        if err != nil { return fmt.Errorf("failed to read response body: %w", err) }
        var apiResp apiResponse
        if err := json.Unmarshal(body, &apiResp); err == nil {
            if apiResp.Code != 0 { return fmt.Errorf("API error: %s (code: %d)", apiResp.Message, apiResp.Code) }
            if apiResp.Output != nil { return json.Unmarshal(apiResp.Output, respBody) }
            if apiResp.Data != nil { return json.Unmarshal(apiResp.Data, respBody) }
        }
        if err := json.Unmarshal(body, respBody); err != nil { return fmt.Errorf("failed to parse response: %w", err) }
    }
    logger.Infof("HTTP request completed successfully - method: %s, url: %s", method, url)
    return nil
}


