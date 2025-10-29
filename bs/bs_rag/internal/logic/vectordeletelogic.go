package logic

import (
	"context"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/svc"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type VectorDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVectorDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VectorDeleteLogic {
	serviceLogger := logger.NewServiceLogger("bs-rag").WithContext(ctx)

	return &VectorDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// 删除向量文档
func (l *VectorDeleteLogic) VectorDelete(in *bs_rag.VectorDeleteRequest) (*bs_rag.VectorDeleteResponse, error) {
	// 参数验证
	if in.CollectionName == "" {
		in.CollectionName = consts.DefaultCollectionName
	}

	if len(in.DocumentIds) == 0 {
		return &bs_rag.VectorDeleteResponse{
			DeletedCount: 0,
			DeletedIds:   []string{},
		}, nil
	}

	// 执行删除
	err := l.svcCtx.VectorProvider.Delete(l.ctx, in.CollectionName, in.DocumentIds)
	if err != nil {
		return &bs_rag.VectorDeleteResponse{
			DeletedCount: 0,
			DeletedIds:   []string{},
			ErrorMessage: err.Error(),
		}, nil
	}

	return &bs_rag.VectorDeleteResponse{
		DeletedCount: int32(len(in.DocumentIds)),
		DeletedIds:   in.DocumentIds,
	}, nil
}
