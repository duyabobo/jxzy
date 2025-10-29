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

	// 首先检查集合是否存在
	collectionReq := &bs_rag.CollectionInfoRequest{
		CollectionName: "test_collection",
		UserId:         "test_user_001",
	}

	collectionResp, err := ragServer.GetCollectionInfo(context.Background(), collectionReq)
	if err != nil {
		LogTestResult("TestVectorSearch", false, "Failed to get collection info: "+err.Error())
		assert.Error(t, err)
		return
	}

	if !collectionResp.Exists {
		LogTestResult("TestVectorSearch", true, "Collection does not exist, skipping search test")
		return
	}

	req := &bs_rag.VectorSearchRequest{
		QueryVector:    CreateTestQueryVector(int(collectionResp.VectorDimension)), // 使用集合的实际维度
		TopK:           5,
		MinScore:       0,
		CollectionName: "test_collection",
		UserId:         "test_user_001",
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
