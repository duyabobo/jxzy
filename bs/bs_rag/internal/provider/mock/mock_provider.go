package mock

import (
	"context"
	"fmt"
	"sync"

	"jxzy/bs/bs_rag/internal/provider/types"
)

// MockProvider 模拟向量数据库提供者
type MockProvider struct {
	collections map[string]*MockCollection
	mutex       sync.RWMutex
}

// MockCollection 模拟集合
type MockCollection struct {
	Name          string
	Dimension     int
	IndexType     string
	DocumentCount int
	Documents     map[string]*MockDocument
}

// MockDocument 模拟文档
type MockDocument struct {
	ID       string
	Vector   []float32
	Metadata map[string]string
	Content  string
}

// NewMockProvider 创建新的模拟提供者
func NewMockProvider() *MockProvider {
	return &MockProvider{
		collections: make(map[string]*MockCollection),
	}
}

// Search 执行向量搜索
func (p *MockProvider) Search(ctx context.Context, collectionName string, queryVector []float32, topK int, minScore float32) ([]types.SearchResult, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	collection, exists := p.collections[collectionName]
	if !exists {
		return nil, types.ErrCollectionNotFound
	}

	// 模拟搜索结果
	results := make([]types.SearchResult, 0)
	for _, doc := range collection.Documents {
		if len(results) >= topK {
			break
		}
		// 简单的相似度计算（这里只是示例）
		score := float32(0.8) // 模拟分数
		if score >= minScore {
			results = append(results, types.SearchResult{
				ID:       doc.ID,
				Score:    score,
				Vector:   doc.Vector,
				Metadata: doc.Metadata,
				Content:  doc.Content,
			})
		}
	}

	return results, nil
}

// Insert 插入向量文档
func (p *MockProvider) Insert(ctx context.Context, collectionName string, documents []types.Document) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	collection, exists := p.collections[collectionName]
	if !exists {
		return types.ErrCollectionNotFound
	}

	for _, doc := range documents {
		collection.Documents[doc.ID] = &MockDocument{
			ID:       doc.ID,
			Vector:   doc.Vector,
			Metadata: doc.Metadata,
			Content:  doc.Content,
		}
		collection.DocumentCount++
	}

	return nil
}

// Delete 删除向量文档
func (p *MockProvider) Delete(ctx context.Context, collectionName string, documentIDs []string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	collection, exists := p.collections[collectionName]
	if !exists {
		return types.ErrCollectionNotFound
	}

	for _, docID := range documentIDs {
		if _, exists := collection.Documents[docID]; exists {
			delete(collection.Documents, docID)
			collection.DocumentCount--
		}
	}

	return nil
}

// GetCollectionInfo 获取集合信息
func (p *MockProvider) GetCollectionInfo(ctx context.Context, collectionName string) (*types.CollectionInfo, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	collection, exists := p.collections[collectionName]
	if !exists {
		return &types.CollectionInfo{
			Name:   collectionName,
			Exists: false,
		}, nil
	}

	return &types.CollectionInfo{
		Name:          collection.Name,
		Dimension:     collection.Dimension,
		IndexType:     collection.IndexType,
		DocumentCount: collection.DocumentCount,
		Metadata:      map[string]string{"provider": "mock"},
	}, nil
}

// CreateCollection 创建集合
func (p *MockProvider) CreateCollection(ctx context.Context, collectionName string, dimension int, indexType string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.collections[collectionName]; exists {
		return fmt.Errorf("collection %s already exists", collectionName)
	}

	p.collections[collectionName] = &MockCollection{
		Name:          collectionName,
		Dimension:     dimension,
		IndexType:     indexType,
		DocumentCount: 0,
		Documents:     make(map[string]*MockDocument),
	}

	return nil
}

// DeleteCollection 删除集合
func (p *MockProvider) DeleteCollection(ctx context.Context, collectionName string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.collections[collectionName]; !exists {
		return types.ErrCollectionNotFound
	}

	delete(p.collections, collectionName)
	return nil
}

// ListCollections 列出所有集合
func (p *MockProvider) ListCollections(ctx context.Context) ([]string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	collections := make([]string, 0, len(p.collections))
	for name := range p.collections {
		collections = append(collections, name)
	}

	return collections, nil
}

// Close 关闭连接
func (p *MockProvider) Close() error {
	// 模拟关闭操作
	return nil
}
