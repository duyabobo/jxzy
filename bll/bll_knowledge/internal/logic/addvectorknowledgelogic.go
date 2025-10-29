package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	knowledgepb "jxzy/bll/bll_knowledge/bll_knowledge"
	"jxzy/bll/bll_knowledge/internal/svc"
	"jxzy/bs/bs_rag/bs_rag"
	consts "jxzy/common/const"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddVectorKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddVectorKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVectorKnowledgeLogic {
	// 使用自定义的 ServiceLogger，在日志中显示服务名
	serviceLogger := logger.NewServiceLogger("bll-knowledge").WithContext(ctx)

	return &AddVectorKnowledgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: serviceLogger,
	}
}

// AddVectorKnowledge 添加知识库到向量数据库
func (l *AddVectorKnowledgeLogic) AddVectorKnowledge(in *knowledgepb.AddVectorKnowledgeRequest) (*knowledgepb.AddVectorKnowledgeResponse, error) {
	l.Logger.Infof("AddVectorKnowledge called with summary: %s, content: %s, user_id: %s",
		in.Summary, in.Content, in.UserId)

	// 1. 验证输入参数
	if err := l.validateInput(in); err != nil {
		l.Logger.Errorf("Input validation failed: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{
			Success: false,
			Message: fmt.Sprintf("输入参数验证失败: %v", err),
		}, nil
	}

	// 2. 生成向量ID（summary和content的MD5）
	vectorId := l.generateVectorId(in.Summary, in.Content)
	l.Logger.Infof("Generated vector ID: %s", vectorId)

	// 3. 构建向量文档（使用text字段，让bs_rag自动向量化）
	document := &bs_rag.VectorDocument{
		Id:   vectorId,
		Text: in.Summary, // 使用text字段，bs_rag会自动生成向量
		Metadata: map[string]string{
			"summary": in.Summary,
			"content": in.Content,
			"user_id": in.UserId,
		},
		Content: in.Content,
	}

	// 5. 检查RAG服务是否可用
	if l.svcCtx.RagRpc == nil {
		l.Logger.Error("RAG service is not available")
		return &knowledgepb.AddVectorKnowledgeResponse{
			VectorId: vectorId, // 即使RAG服务不可用，也要返回向量ID
			Success:  false,
			Message:  "RAG服务不可用",
		}, nil
	}

	// 6. 调用RAG服务插入向量
	ragReq := &bs_rag.VectorInsertRequest{
		CollectionName: consts.DefaultCollectionName,
		Documents:      []*bs_rag.VectorDocument{document},
		UserId:         in.UserId,
	}

	ragResp, err := l.svcCtx.RagRpc.VectorInsert(l.ctx, ragReq)
	if err != nil {
		l.Logger.Errorf("Failed to insert vector to RAG service: %v", err)
		return &knowledgepb.AddVectorKnowledgeResponse{
			VectorId: vectorId, // 即使插入失败，也要返回向量ID
			Success:  false,
			Message:  fmt.Sprintf("插入向量到RAG服务失败: %v", err),
		}, nil
	}

	l.Logger.Infof("Successfully inserted vector to RAG service, response: %v", ragResp)

	return &knowledgepb.AddVectorKnowledgeResponse{
		VectorId: vectorId,
		Success:  true,
		Message:  "知识库添加成功",
	}, nil
}

// validateInput 验证输入参数
func (l *AddVectorKnowledgeLogic) validateInput(in *knowledgepb.AddVectorKnowledgeRequest) error {
	if strings.TrimSpace(in.Summary) == "" {
		return fmt.Errorf("summary不能为空")
	}
	if strings.TrimSpace(in.Content) == "" {
		return fmt.Errorf("content不能为空")
	}
	if strings.TrimSpace(in.UserId) == "" {
		return fmt.Errorf("user_id不能为空")
	}
	return nil
}

// generateVectorId 生成向量ID（summary和content的MD5）
func (l *AddVectorKnowledgeLogic) generateVectorId(summary, content string) string {
	// 组合summary和content
	combined := summary + "|||" + content

	// 计算MD5
	hash := md5.Sum([]byte(combined))

	// 转换为十六进制字符串
	return hex.EncodeToString(hash[:])
}
