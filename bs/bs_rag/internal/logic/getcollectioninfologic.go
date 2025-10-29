package logic

import (
	"context"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/svc"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCollectionInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCollectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCollectionInfoLogic {
	serviceLogger := logger.NewServiceLogger("bs-rag").WithContext(ctx)

	return &GetCollectionInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// 获取集合信息
func (l *GetCollectionInfoLogic) GetCollectionInfo(in *bs_rag.CollectionInfoRequest) (*bs_rag.CollectionInfoResponse, error) {
	// 参数验证
	if in.CollectionName == "" {
		in.CollectionName = consts.DefaultCollectionName
	}

	// 获取集合信息
	info, err := l.svcCtx.VectorProvider.GetCollectionInfo(l.ctx, in.CollectionName)
	if err != nil {
		return &bs_rag.CollectionInfoResponse{
			CollectionName: in.CollectionName,
			Exists:         false,
			ErrorMessage:   err.Error(),
		}, nil
	}

	return &bs_rag.CollectionInfoResponse{
		CollectionName:  info.Name,
		DocumentCount:   int32(info.DocumentCount),
		VectorDimension: int32(info.Dimension),
		IndexType:       info.IndexType,
		Exists:          info.Exists,
	}, nil
}
