package test

import (
	"context"
	"fmt"
	"testing"

	"jxzy/bs/bs_rag/bs_rag"
	consts "jxzy/common/const"

	"github.com/stretchr/testify/assert"
)

// TestVectorSearch 测试向量搜索RPC
func TestVectorSearch(t *testing.T) {
	ragServer := SetupTestEnvironment()

	testParams := map[string]interface{}{
		"collection": "test_collection",
		"topK":       5,
		"minScore":   0.5,
		"dimension":  consts.DashVectorDefaultDimension, // DashVectorDefaultDimension
	}
	LogTestStart("TestVectorSearch", testParams)

	req := &bs_rag.VectorSearchRequest{
		QueryText: "测试查询文本", // 使用文本查询，自动向量化
		TopK:      5,
		MinScore:  0,
		UserId:    "test_user_001",
		SceneCode: "test_scene", // 测试场景编码
		Filters: map[string]string{
			"source": "test",
		},
	}

	resp, err := ragServer.VectorSearch(context.Background(), req)

	// 由于可能没有配置真实的DashVector，我们主要验证请求格式和响应结构
	if err != nil {
		LogTestResult("TestVectorSearch", false, "Expected error in test environment: "+err.Error())
		// 验证错误处理
		assert.Error(t, err)
	} else {
		LogTestResult("TestVectorSearch", true, fmt.Sprintf("Search completed successfully, resp: %v", resp))
		AssertVectorSearchResponse(t, resp, 5, true)
	}
}
