package knowledge

import (
	"context"

	"jxzy/apis/api_knowledge/internal/svc"
	"jxzy/apis/api_knowledge/internal/types"
	"jxzy/bll/bll_context/bll_context"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddKnowledgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddKnowledgeLogic {
	serviceLogger := logger.NewServiceLogger("api-knowledge").WithContext(ctx)

	return &AddKnowledgeLogic{
		Logger: serviceLogger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddKnowledgeLogic) AddKnowledge(req *types.AddKnowledgeRequest) (resp *types.AddKnowledgeResponse, err error) {
	l.Logger.Infof("AddKnowledge called with summary: %s, content: %s, user_id: %s",
		req.Summary, req.Content, req.UserId)

	// 调用bll_context的AddVectorKnowledge RPC
	rpcReq := &bll_context.AddVectorKnowledgeRequest{
		Summary: req.Summary,
		Content: req.Content,
		UserId:  req.UserId,
	}

	rpcResp, err := l.svcCtx.BllContextRpc.AddVectorKnowledge(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Failed to call AddVectorKnowledge RPC: %v", err)
		return &types.AddKnowledgeResponse{
			VectorId: "",
			Success:  false,
			Message:  "调用知识库服务失败: " + err.Error(),
		}, nil
	}

	return &types.AddKnowledgeResponse{
		VectorId: rpcResp.VectorId,
		Success:  rpcResp.Success,
		Message:  rpcResp.Message,
	}, nil
}
