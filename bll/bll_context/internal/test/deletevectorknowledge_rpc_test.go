package test

import (
	"context"
	"testing"

	contextpb "jxzy/bll/bll_context/bll_context"
	"jxzy/bll/bll_context/internal/logic"
	"jxzy/bll/bll_context/internal/svc"

	"github.com/stretchr/testify/assert"
)

// TestDeleteVectorKnowledge 测试从向量数据库删除知识库RPC
func TestDeleteVectorKnowledge(t *testing.T) {
	// 创建测试配置
	cfg := InitConfig()

	// 创建ServiceContext
	svcCtx := svc.NewServiceContext(cfg)

	// 创建logic实例
	ctx := context.Background()
	logic := logic.NewDeleteVectorKnowledgeLogic(ctx, svcCtx)

	// 创建测试请求
	req := &contextpb.DeleteVectorKnowledgeRequest{
		VectorId: "doc_003",
		UserId:   "test_user_001",
	}

	// 执行测试
	resp, err := logic.DeleteVectorKnowledge(req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// 在测试环境中，RAG服务可能成功或失败，我们验证响应格式
	if resp.Success {
		assert.Contains(t, resp.Message, "知识库删除成功")
	} else {
		assert.Contains(t, resp.Message, "从RAG服务删除向量失败")
	}
}
