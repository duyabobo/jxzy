package logic

import (
	"context"
	"fmt"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/svc"
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
	if len(in.DocumentIds) == 0 {
		return &bs_rag.VectorDeleteResponse{
			DeletedCount: 0,
			DeletedIds:   []string{},
		}, nil
	}

	if in.SceneCode == "" {
		return &bs_rag.VectorDeleteResponse{
			DeletedCount: 0,
			DeletedIds:   []string{},
			ErrorMessage: "scene_code is required",
		}, nil
	}

	// 根据scene_code获取collection_name
	collectionName, err := l.svcCtx.GetCollectionName(l.ctx, in.SceneCode)
	if err != nil {
		l.Logger.Errorf("Failed to get collection_name for scene_code %s: %v", in.SceneCode, err)
		return &bs_rag.VectorDeleteResponse{
			DeletedCount: 0,
			DeletedIds:   []string{},
			ErrorMessage: fmt.Sprintf("获取collection_name失败: %v", err),
		}, nil
	}

	// 执行删除
	err = l.svcCtx.VectorProvider.Delete(l.ctx, collectionName, in.DocumentIds)
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
