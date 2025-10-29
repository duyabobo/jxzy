package logic

import (
	"context"
	"fmt"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/svc"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type VectorizeTextLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVectorizeTextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VectorizeTextLogic {
	serviceLogger := logger.NewServiceLogger("bs-rag").WithContext(ctx)

	return &VectorizeTextLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// VectorizeText 向量化文本
func (l *VectorizeTextLogic) VectorizeText(in *bs_rag.VectorizeTextRequest) (*bs_rag.VectorizeTextResponse, error) {
	// 参数验证
	if in.Text == "" {
		return &bs_rag.VectorizeTextResponse{
			Vector:       []float32{},
			ErrorMessage: "文本不能为空",
		}, nil
	}

	l.Logger.Infof("Vectorizing text, length: %d", len(in.Text))

	// 生成向量
	vector, err := l.svcCtx.EmbeddingService.GenerateEmbedding(in.Text)
	if err != nil {
		l.Logger.Errorf("Failed to generate embedding: %v", err)
		return &bs_rag.VectorizeTextResponse{
			Vector:       []float32{},
			ErrorMessage: fmt.Sprintf("生成向量失败: %v", err),
		}, nil
	}

	l.Logger.Debugf("Generated vector, length: %d", len(vector))

	// 直接返回 float32 向量（proto 中使用 fixed32 类型）
	return &bs_rag.VectorizeTextResponse{
		Vector: vector,
	}, nil
}
