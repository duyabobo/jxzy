package faiss

import (
	"context"
	"fmt"
	"os"
	"sync"

	"jxzy/bs/bs_rag/internal/config"
	"jxzy/bs/bs_rag/internal/provider/types"
	consts "jxzy/common/const"
)

// FaissProvider 提供 Faiss 向量数据库操作
type FaissProvider struct {
	config      config.FaissConfig
	collections map[string]*FaissCollection
	mutex       sync.RWMutex
}

// FaissCollection 表示一个 Faiss 集合
type FaissCollection struct {
	Name          string
	Dimension     int
	IndexType     string
	DocumentCount int
	// 这里将来会添加实际的 Faiss 索引对象
	// index *faiss.Index
}

// NewFaissProvider 创建新的 Faiss 提供者
func NewFaissProvider(config config.FaissConfig) *FaissProvider {
	// 确保索引目录存在
	if err := os.MkdirAll(config.IndexPath, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create index directory: %v", err))
	}

	return &FaissProvider{
		config:      config,
		collections: make(map[string]*FaissCollection),
	}
}

// GetCollection 获取或创建集合
func (p *FaissProvider) GetCollection(name string) (*FaissCollection, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if collection, exists := p.collections[name]; exists {
		return collection, nil
	}

	// 创建新集合
	collection := &FaissCollection{
		Name:          name,
		Dimension:     consts.FaissDefaultDimension, // FaissDefaultDimension
		IndexType:     p.config.IndexType,
		DocumentCount: 0,
	}

	p.collections[name] = collection
	return collection, nil
}

// Search 执行向量搜索
func (p *FaissProvider) Search(ctx context.Context, collectionName string, queryVector []float32, topK int, minScore float32) ([]types.SearchResult, error) {
	_, err := p.GetCollection(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection %s: %v", collectionName, err)
	}

	// TODO: 实现实际的 Faiss 搜索
	// 这里返回模拟结果
	results := []types.SearchResult{
		{
			ID:       "doc1",
			Score:    0.95,
			Vector:   queryVector,
			Metadata: map[string]string{"source": "test"},
			Content:  "This is a test document",
		},
	}

	return results, nil
}

// Insert 插入向量文档
func (p *FaissProvider) Insert(ctx context.Context, collectionName string, documents []types.Document) error {
	collection, err := p.GetCollection(collectionName)
	if err != nil {
		return fmt.Errorf("failed to get collection %s: %v", collectionName, err)
	}

	// TODO: 实现实际的 Faiss 插入
	collection.DocumentCount += len(documents)
	return nil
}

// Delete 删除向量文档
func (p *FaissProvider) Delete(ctx context.Context, collectionName string, documentIDs []string) error {
	collection, err := p.GetCollection(collectionName)
	if err != nil {
		return fmt.Errorf("failed to get collection %s: %v", collectionName, err)
	}

	// TODO: 实现实际的 Faiss 删除
	collection.DocumentCount -= len(documentIDs)
	if collection.DocumentCount < 0 {
		collection.DocumentCount = 0
	}
	return nil
}

// GetCollectionInfo 获取集合信息
func (p *FaissProvider) GetCollectionInfo(ctx context.Context, collectionName string) (*types.CollectionInfo, error) {
	collection, err := p.GetCollection(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection %s: %v", collectionName, err)
	}

	return &types.CollectionInfo{
		Name:          collection.Name,
		Dimension:     collection.Dimension,
		IndexType:     collection.IndexType,
		DocumentCount: collection.DocumentCount,
		Metadata:      map[string]string{"provider": "faiss"},
	}, nil
}

// CreateCollection 创建集合
func (p *FaissProvider) CreateCollection(ctx context.Context, collectionName string, dimension int, indexType string) error {
	// TODO: 实现实际的 Faiss 集合创建
	collection := &FaissCollection{
		Name:          collectionName,
		Dimension:     dimension,
		IndexType:     indexType,
		DocumentCount: 0,
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.collections[collectionName] = collection

	return nil
}

// DeleteCollection 删除集合
func (p *FaissProvider) DeleteCollection(ctx context.Context, collectionName string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.collections[collectionName]; !exists {
		return types.ErrCollectionNotFound
	}

	delete(p.collections, collectionName)
	return nil
}

// ListCollections 列出所有集合
func (p *FaissProvider) ListCollections(ctx context.Context) ([]string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	collections := make([]string, 0, len(p.collections))
	for name := range p.collections {
		collections = append(collections, name)
	}

	return collections, nil
}

// Close 关闭连接
func (p *FaissProvider) Close() error {
	// TODO: 实现实际的 Faiss 连接关闭
	return nil
}
