package test

import (
	"context"
	"testing"

	contextpb "jxzy/bll/bll_context/bll_context"
	"jxzy/bll/bll_context/internal/logic"
	"jxzy/bll/bll_context/internal/svc"

	"github.com/stretchr/testify/assert"
)

// TestAddVectorKnowledge 测试添加知识库到向量数据库RPC
func TestAddVectorKnowledge(t *testing.T) {
	// 创建测试配置
	cfg := InitConfig()

	// 创建ServiceContext
	svcCtx := svc.NewServiceContext(cfg)

	// 创建logic实例
	ctx := context.Background()
	logic := logic.NewAddVectorKnowledgeLogic(ctx, svcCtx)

	// 创建测试请求
	req := &contextpb.AddVectorKnowledgeRequest{
		Summary: "这是一个测试总结",
		Content: "这是测试内容，用于验证知识库添加功能",
		UserId:  "test_user_001",
	}

	// 执行测试
	resp, err := logic.AddVectorKnowledge(req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.VectorId)
}
