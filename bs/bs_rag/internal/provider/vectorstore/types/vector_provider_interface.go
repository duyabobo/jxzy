package types

import (
    "context"
)

// VectorProvider 向量数据库提供者接口
type VectorProvider interface {
    Search(ctx context.Context, collectionName string, queryVector []float32, topK int, minScore float32) ([]SearchResult, error)
    Insert(ctx context.Context, collectionName string, documents []Document) error
    Delete(ctx context.Context, collectionName string, documentIDs []string) error
    GetCollectionInfo(ctx context.Context, collectionName string) (*CollectionInfo, error)
    CreateCollection(ctx context.Context, collectionName string, dimension int, indexType string) error
    DeleteCollection(ctx context.Context, collectionName string) error
    ListCollections(ctx context.Context) ([]string, error)
    Close() error
}

type VectorProviderType string

const (
    VectorProviderTypeFaiss      VectorProviderType = "faiss"
    VectorProviderTypeMilvus     VectorProviderType = "milvus"
    VectorProviderTypePinecone   VectorProviderType = "pinecone"
    VectorProviderTypeWeaviate   VectorProviderType = "weaviate"
    VectorProviderTypeDashVector VectorProviderType = "dashvector"
    VectorProviderTypeMock       VectorProviderType = "mock"
)

var (
    ErrInvalidConfig       = &VectorProviderError{Message: "invalid configuration"}
    ErrUnsupportedProvider = &VectorProviderError{Message: "unsupported vector provider"}
    ErrCollectionNotFound  = &VectorProviderError{Message: "collection not found"}
    ErrDocumentNotFound    = &VectorProviderError{Message: "document not found"}
)

type VectorProviderError struct{
    Message string
}

func (e *VectorProviderError) Error() string { return e.Message }

type SearchResult struct {
    ID       string            `json:"id"`
    Score    float32           `json:"score"`
    Vector   []float32         `json:"vector,omitempty"`
    Metadata map[string]string `json:"metadata,omitempty"`
    Content  string            `json:"content,omitempty"`
}

type Document struct {
    ID       string            `json:"id"`
    Vector   []float32         `json:"vector"`
    Metadata map[string]string `json:"metadata,omitempty"`
    Content  string            `json:"content,omitempty"`
}

type CollectionInfo struct {
    Name          string            `json:"name"`
    Dimension     int               `json:"dimension"`
    IndexType     string            `json:"index_type"`
    DocumentCount int               `json:"document_count"`
    Metadata      map[string]string `json:"metadata,omitempty"`
    Exists        bool              `json:"exists"`
}


