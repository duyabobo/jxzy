package logic

import (
	"context"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/provider/types"
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
	documents := make([]types.Document, len(in.Documents))
	insertedIds := make([]string, 0, len(in.Documents))

	for i, doc := range in.Documents {
		// 转换向量
		vector := make([]float32, len(doc.Vector))
		for j, v := range doc.Vector {
			vector[j] = float32(v)
		}

		documents[i] = types.Document{
			ID:       doc.Id,
			Vector:   vector,
			Metadata: doc.Metadata,
			Content:  doc.Content,
		}
		insertedIds = append(insertedIds, doc.Id)
	}

	// 执行插入
	err := l.svcCtx.VectorProvider.Insert(l.ctx, in.CollectionName, documents)
	if err != nil {
		return nil, err
	}

	return &bs_rag.VectorInsertResponse{
		InsertedCount: int32(len(documents)),
		InsertedIds:   insertedIds,
	}, nil
}
