package knowledge

import (
	"context"

	"jxzy/apis/api_knowledge/internal/svc"
	"jxzy/apis/api_knowledge/internal/types"
	"jxzy/bll/bll_context/bll_context"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteKnowledgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKnowledgeLogic {
	serviceLogger := logger.NewServiceLogger("api-knowledge").WithContext(ctx)

	return &DeleteKnowledgeLogic{
		Logger: serviceLogger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteKnowledgeLogic) DeleteKnowledge(req *types.DeleteKnowledgeRequest) (resp *types.DeleteKnowledgeResponse, err error) {
	l.Logger.Infof("DeleteKnowledge called with vector_id: %s, user_id: %s", req.VectorId, req.UserId)

	// 调用bll_context的DeleteVectorKnowledge RPC
	rpcReq := &bll_context.DeleteVectorKnowledgeRequest{
		VectorId: req.VectorId,
		UserId:   req.UserId,
	}

	rpcResp, err := l.svcCtx.BllContextRpc.DeleteVectorKnowledge(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to call DeleteVectorKnowledge RPC: %v", err)
		return &types.DeleteKnowledgeResponse{
			Success: false,
			Message: "调用知识库服务失败: " + err.Error(),
		}, nil
	}

	return &types.DeleteKnowledgeResponse{
		Success: rpcResp.Success,
		Message: rpcResp.Message,
	}, nil
}
