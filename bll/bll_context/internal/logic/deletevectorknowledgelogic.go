package logic

import (
	"context"
	"fmt"
	"strings"

	contextpb "jxzy/bll/bll_context/bll_context"
	"jxzy/bll/bll_context/internal/svc"
	"jxzy/bs/bs_rag/bs_rag"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteVectorKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteVectorKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVectorKnowledgeLogic {
	// 使用自定义的 ServiceLogger，在日志中显示服务名
	serviceLogger := logger.NewServiceLogger("bll-context").WithContext(ctx)

	return &DeleteVectorKnowledgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// DeleteVectorKnowledge 从向量数据库删除知识库
func (l *DeleteVectorKnowledgeLogic) DeleteVectorKnowledge(in *contextpb.DeleteVectorKnowledgeRequest) (*contextpb.DeleteVectorKnowledgeResponse, error) {
	l.Logger.Infof("DeleteVectorKnowledge called with vector_id: %s, user_id: %s", in.VectorId, in.UserId)

	// 1. 验证输入参数
	if err := l.validateInput(in); err != nil {
		l.Logger.Errorf("Input validation failed: %v", err)
		return &contextpb.DeleteVectorKnowledgeResponse{
			Success: false,
			Message: fmt.Sprintf("输入参数验证失败: %v", err),
		}, nil
	}

	// 2. 检查RAG服务是否可用
	if l.svcCtx.RagRpc == nil {
		l.Logger.Error("RAG service is not available")
		return &contextpb.DeleteVectorKnowledgeResponse{
			Success: false,
			Message: "RAG服务不可用",
		}, nil
	}

	// 3. 调用RAG服务删除向量
	ragReq := &bs_rag.VectorDeleteRequest{
		CollectionName: consts.DefaultCollectionName,
		DocumentIds:    []string{in.VectorId},
		UserId:         in.UserId,
	}

	ragResp, err := l.svcCtx.RagRpc.VectorDelete(l.ctx, ragReq)
	if err != nil {
		l.Logger.Errorf("Failed to delete vector from RAG service: %v", err)
		return &contextpb.DeleteVectorKnowledgeResponse{
			Success: false,
			Message: fmt.Sprintf("从RAG服务删除向量失败: %v", err),
		}, nil
	}

	l.Logger.Infof("Successfully deleted vector from RAG service, response: %v", ragResp)

	// 检查RAG响应中是否有错误信息
	if ragResp != nil && ragResp.ErrorMessage != "" {
		l.Logger.Errorf("RAG service returned error: %s", ragResp.ErrorMessage)
		return &contextpb.DeleteVectorKnowledgeResponse{
			Success: false,
			Message: fmt.Sprintf("从RAG服务删除向量失败: %s", ragResp.ErrorMessage),
		}, nil
	}

	return &contextpb.DeleteVectorKnowledgeResponse{
		Success: true,
		Message: "知识库删除成功",
	}, nil
}

// validateInput 验证输入参数
func (l *DeleteVectorKnowledgeLogic) validateInput(in *contextpb.DeleteVectorKnowledgeRequest) error {
	if strings.TrimSpace(in.VectorId) == "" {
		return fmt.Errorf("vector_id不能为空")
	}
	if strings.TrimSpace(in.UserId) == "" {
		return fmt.Errorf("user_id不能为空")
	}
	return nil
}
