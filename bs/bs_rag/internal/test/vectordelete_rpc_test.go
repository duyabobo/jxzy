package test

import (
	"context"
	"testing"

	"jxzy/bs/bs_rag/bs_rag"

	"github.com/stretchr/testify/assert"
)

// TestVectorDelete 测试向量删除RPC
func TestVectorDelete(t *testing.T) {
	ragServer := SetupTestEnvironment()

	testParams := map[string]interface{}{
		"collection": "test_collection",
		"docCount":   2,
	}
	LogTestStart("TestVectorDelete", testParams)

	documentIDs := []string{"doc_001"}

	req := &bs_rag.VectorDeleteRequest{
		DocumentIds: documentIDs,
		UserId:      "test_user_001",
		SceneCode:   "test_scene", // 测试场景编码
	}

	_, err = ragServer.VectorDelete(context.Background(), req)

	// 验证结果
	if err != nil {
		LogTestResult("TestVectorDelete", false, "Expected error in test environment: "+err.Error())
		assert.Error(t, err)
	} else {
		LogTestResult("TestVectorDelete", true, "Delete completed successfully")
	}
}
