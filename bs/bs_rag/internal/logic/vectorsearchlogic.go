package logic

import (
	"context"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/svc"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type VectorSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVectorSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VectorSearchLogic {
	serviceLogger := logger.NewServiceLogger("bs-rag").WithContext(ctx)

	return &VectorSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// 向量相似度搜索
func (l *VectorSearchLogic) VectorSearch(in *bs_rag.VectorSearchRequest) (*bs_rag.VectorSearchResponse, error) {
	// 参数验证
	if in.QueryText == "" {
		l.Logger.Errorf("QueryText is required")
		return &bs_rag.VectorSearchResponse{
			Results:      []*bs_rag.VectorSearchResult{},
			TotalCount:   0,
			SearchTimeMs: 0,
		}, nil
	}

	// 自动生成查询向量
	l.Logger.Infof("Generating query vector from text, length: %d", len(in.QueryText))
	queryVector, err := l.svcCtx.EmbeddingService.GenerateEmbedding(in.QueryText)
	if err != nil {
		l.Logger.Errorf("Failed to generate embedding for query text: %v", err)
		return &bs_rag.VectorSearchResponse{
			Results:      []*bs_rag.VectorSearchResult{},
			TotalCount:   0,
			SearchTimeMs: 0,
		}, nil
	}
	l.Logger.Debugf("Generated query vector, length: %d", len(queryVector))

	if in.TopK <= 0 {
		in.TopK = 10
	}

	if in.CollectionName == "" {
		in.CollectionName = consts.DefaultCollectionName
	}

	// 执行搜索
	results, err := l.svcCtx.VectorProvider.Search(l.ctx, in.CollectionName, queryVector, int(in.TopK), float32(in.MinScore))
	if err != nil {
		return nil, err
	}

	// 转换结果
	protoResults := make([]*bs_rag.VectorSearchResult, len(results))
	for i, result := range results {
		protoResults[i] = &bs_rag.VectorSearchResult{
			Id:       result.ID,
			Vector:   result.Vector,
			Score:    result.Score,
			Metadata: result.Metadata,
			Content:  result.Content,
		}
	}

	return &bs_rag.VectorSearchResponse{
		Results:      protoResults,
		TotalCount:   int32(len(protoResults)),
		SearchTimeMs: 0, // TODO: 添加实际搜索时间
	}, nil
}
