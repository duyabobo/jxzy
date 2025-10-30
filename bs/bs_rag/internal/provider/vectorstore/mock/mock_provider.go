package mock

import (
    "context"
    "fmt"
    "sync"

    "jxzy/bs/bs_rag/internal/provider/vectorstore/types"
)

type MockProvider struct {
    collections map[string]*MockCollection
    mutex       sync.RWMutex
}

type MockCollection struct {
    Name          string
    Dimension     int
    IndexType     string
    DocumentCount int
    Documents     map[string]*MockDocument
}

type MockDocument struct {
    ID       string
    Vector   []float32
    Metadata map[string]string
    Content  string
}

func NewMockProvider() *MockProvider {
    return &MockProvider{collections: make(map[string]*MockCollection)}
}

func (p *MockProvider) Search(ctx context.Context, collectionName string, queryVector []float32, topK int, minScore float32) ([]types.SearchResult, error) {
    p.mutex.RLock(); defer p.mutex.RUnlock()
    c, ok := p.collections[collectionName]
    if !ok { return nil, types.ErrCollectionNotFound }
    results := make([]types.SearchResult, 0)
    for _, d := range c.Documents {
        if len(results) >= topK { break }
        score := float32(0.8)
        if score >= minScore { results = append(results, types.SearchResult{ID: d.ID, Score: score, Vector: d.Vector, Metadata: d.Metadata, Content: d.Content}) }
    }
    return results, nil
}

func (p *MockProvider) Insert(ctx context.Context, collectionName string, documents []types.Document) error {
    p.mutex.Lock(); defer p.mutex.Unlock()
    c, ok := p.collections[collectionName]
    if !ok { return types.ErrCollectionNotFound }
    for _, doc := range documents {
        c.Documents[doc.ID] = &MockDocument{ID: doc.ID, Vector: doc.Vector, Metadata: doc.Metadata, Content: doc.Content}
        c.DocumentCount++
    }
    return nil
}

func (p *MockProvider) Delete(ctx context.Context, collectionName string, documentIDs []string) error {
    p.mutex.Lock(); defer p.mutex.Unlock()
    c, ok := p.collections[collectionName]
    if !ok { return types.ErrCollectionNotFound }
    for _, id := range documentIDs {
        if _, ok := c.Documents[id]; ok { delete(c.Documents, id); c.DocumentCount-- }
    }
    return nil
}

func (p *MockProvider) GetCollectionInfo(ctx context.Context, collectionName string) (*types.CollectionInfo, error) {
    p.mutex.RLock(); defer p.mutex.RUnlock()
    c, ok := p.collections[collectionName]
    if !ok { return &types.CollectionInfo{Name: collectionName, Exists: false}, nil }
    return &types.CollectionInfo{Name: c.Name, Dimension: c.Dimension, IndexType: c.IndexType, DocumentCount: c.DocumentCount, Metadata: map[string]string{"provider": "mock"}}, nil
}

func (p *MockProvider) CreateCollection(ctx context.Context, collectionName string, dimension int, indexType string) error {
    p.mutex.Lock(); defer p.mutex.Unlock()
    if _, ok := p.collections[collectionName]; ok { return fmt.Errorf("collection %s already exists", collectionName) }
    p.collections[collectionName] = &MockCollection{Name: collectionName, Dimension: dimension, IndexType: indexType, DocumentCount: 0, Documents: make(map[string]*MockDocument)}
    return nil
}

func (p *MockProvider) DeleteCollection(ctx context.Context, collectionName string) error {
    p.mutex.Lock(); defer p.mutex.Unlock()
    if _, ok := p.collections[collectionName]; !ok { return types.ErrCollectionNotFound }
    delete(p.collections, collectionName)
    return nil
}

func (p *MockProvider) ListCollections(ctx context.Context) ([]string, error) {
    p.mutex.RLock(); defer p.mutex.RUnlock()
    names := make([]string, 0, len(p.collections))
    for n := range p.collections { names = append(names, n) }
    return names, nil
}

func (p *MockProvider) Close() error { return nil }


