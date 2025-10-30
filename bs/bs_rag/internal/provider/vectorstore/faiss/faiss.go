package faiss

import (
    "context"
    "fmt"
    "os"
    "sync"

    "jxzy/bs/bs_rag/internal/config"
    ptypes "jxzy/bs/bs_rag/internal/provider/vectorstore/types"
    consts "jxzy/common/const"
)

type FaissProvider struct {
    config      config.FaissConfig
    collections map[string]*FaissCollection
    mutex       sync.RWMutex
}

type FaissCollection struct {
    Name          string
    Dimension     int
    IndexType     string
    DocumentCount int
}

func NewFaissProvider(config config.FaissConfig) *FaissProvider {
    if err := os.MkdirAll(config.IndexPath, 0755); err != nil { panic(fmt.Sprintf("Failed to create index directory: %v", err)) }
    return &FaissProvider{config: config, collections: make(map[string]*FaissCollection)}
}

func (p *FaissProvider) GetCollection(name string) (*FaissCollection, error) {
    p.mutex.Lock(); defer p.mutex.Unlock()
    if c, ok := p.collections[name]; ok { return c, nil }
    c := &FaissCollection{Name: name, Dimension: consts.FaissDefaultDimension, IndexType: p.config.IndexType, DocumentCount: 0}
    p.collections[name] = c
    return c, nil
}

func (p *FaissProvider) Search(ctx context.Context, collectionName string, queryVector []float32, topK int, minScore float32) ([]ptypes.SearchResult, error) {
    if _, err := p.GetCollection(collectionName); err != nil { return nil, fmt.Errorf("failed to get collection %s: %v", collectionName, err) }
    results := []ptypes.SearchResult{{ID: "doc1", Score: 0.95, Vector: queryVector, Metadata: map[string]string{"source": "test"}, Content: "This is a test document"}}
    return results, nil
}

func (p *FaissProvider) Insert(ctx context.Context, collectionName string, documents []ptypes.Document) error {
    c, err := p.GetCollection(collectionName); if err != nil { return fmt.Errorf("failed to get collection %s: %v", collectionName, err) }
    c.DocumentCount += len(documents)
    return nil
}

func (p *FaissProvider) Delete(ctx context.Context, collectionName string, documentIDs []string) error {
    c, err := p.GetCollection(collectionName); if err != nil { return fmt.Errorf("failed to get collection %s: %v", collectionName, err) }
    c.DocumentCount -= len(documentIDs); if c.DocumentCount < 0 { c.DocumentCount = 0 }
    return nil
}

func (p *FaissProvider) GetCollectionInfo(ctx context.Context, collectionName string) (*ptypes.CollectionInfo, error) {
    c, err := p.GetCollection(collectionName); if err != nil { return nil, fmt.Errorf("failed to get collection %s: %v", collectionName, err) }
    return &ptypes.CollectionInfo{Name: c.Name, Dimension: c.Dimension, IndexType: c.IndexType, DocumentCount: c.DocumentCount, Metadata: map[string]string{"provider": "faiss"}}, nil
}

func (p *FaissProvider) CreateCollection(ctx context.Context, collectionName string, dimension int, indexType string) error {
    c := &FaissCollection{Name: collectionName, Dimension: dimension, IndexType: indexType, DocumentCount: 0}
    p.mutex.Lock(); defer p.mutex.Unlock(); p.collections[collectionName] = c
    return nil
}

func (p *FaissProvider) DeleteCollection(ctx context.Context, collectionName string) error {
    p.mutex.Lock(); defer p.mutex.Unlock()
    if _, ok := p.collections[collectionName]; !ok { return ptypes.ErrCollectionNotFound }
    delete(p.collections, collectionName)
    return nil
}

func (p *FaissProvider) ListCollections(ctx context.Context) ([]string, error) {
    p.mutex.RLock(); defer p.mutex.RUnlock()
    names := make([]string, 0, len(p.collections))
    for n := range p.collections { names = append(names, n) }
    return names, nil
}

func (p *FaissProvider) Close() error { return nil }



