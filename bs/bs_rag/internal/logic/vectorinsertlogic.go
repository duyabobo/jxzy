package logic

import (
	"context"
	"fmt"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/provider/vectorstore/types"
	"jxzy/bs/bs_rag/internal/svc"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type VectorInsertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVectorInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VectorInsertLogic {
	serviceLogger := logger.NewServiceLogger("bs-rag").WithContext(ctx)

	return &VectorInsertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// 插入向量文档
func (l *VectorInsertLogic) VectorInsert(in *bs_rag.VectorInsertRequest) (*bs_rag.VectorInsertResponse, error) {
	// 参数验证
	if in.CollectionName == "" {
		in.CollectionName = consts.DefaultCollectionName
	}

	if len(in.Documents) == 0 {
		return &bs_rag.VectorInsertResponse{
			InsertedCount: 0,
			InsertedIds:   []string{},
		}, nil
	}

	// 转换文档
	documents := make([]types.Document, 0, len(in.Documents))
	insertedIds := make([]string, 0, len(in.Documents))

	for _, doc := range in.Documents {
		// 参数验证
		if doc.Text == "" {
			l.Logger.Errorf("Document id %s has no text, skipping", doc.Id)
			continue
		}

		// 自动生成向量
		l.Logger.Infof("Auto-generating vector for document id: %s, text length: %d", doc.Id, len(doc.Text))
		vector, err := l.svcCtx.EmbeddingService.GenerateEmbedding(doc.Text)
		if err != nil {
			l.Logger.Errorf("Failed to generate embedding for document id %s: %v", doc.Id, err)
			return &bs_rag.VectorInsertResponse{
				InsertedCount: 0,
				InsertedIds:   []string{},
				ErrorMessage:  fmt.Sprintf("生成向量失败 (文档ID: %s): %v", doc.Id, err),
			}, nil
		}
		l.Logger.Debugf("Generated vector for document id %s: length=%d", doc.Id, len(vector))

		documents = append(documents, types.Document{
			ID:       doc.Id,
			Vector:   vector,
			Metadata: doc.Metadata,
			Content:  doc.Content,
		})
		insertedIds = append(insertedIds, doc.Id)
	}

	if len(documents) == 0 {
		return &bs_rag.VectorInsertResponse{
			InsertedCount: 0,
			InsertedIds:   []string{},
			ErrorMessage:  "没有有效的文档可插入",
		}, nil
	}

	// 执行插入
	err := l.svcCtx.VectorProvider.Insert(l.ctx, in.CollectionName, documents)
	if err != nil {
		return &bs_rag.VectorInsertResponse{
			InsertedCount: 0,
			InsertedIds:   []string{},
			ErrorMessage:  fmt.Sprintf("插入向量失败: %v", err),
		}, nil
	}

	l.Logger.Infof("Successfully inserted %d documents to collection: %s", len(documents), in.CollectionName)

	return &bs_rag.VectorInsertResponse{
		InsertedCount: int32(len(documents)),
		InsertedIds:   insertedIds,
	}, nil
}
