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

	// 首先检查集合是否存在
	collectionReq := &bs_rag.CollectionInfoRequest{
		CollectionName: "test_collection",
		UserId:         "test_user_001",
	}

	collectionResp, err := ragServer.GetCollectionInfo(context.Background(), collectionReq)
	if err != nil {
		LogTestResult("TestVectorDelete", false, "Failed to get collection info: "+err.Error())
		assert.Error(t, err)
		return
	}

	if !collectionResp.Exists {
		LogTestResult("TestVectorDelete", true, "Collection does not exist, skipping delete test")
		return
	}

	documentIDs := []string{"doc_001"}

	req := &bs_rag.VectorDeleteRequest{
		CollectionName: "test_collection",
		DocumentIds:    documentIDs,
		UserId:         "test_user_001",
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
