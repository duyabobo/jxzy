package test

import (
	"context"
	"fmt"
	"testing"

	"jxzy/bs/bs_rag/bs_rag"

	"github.com/stretchr/testify/assert"
)

// TestGetCollectionInfo 测试获取集合信息RPC
func TestGetCollectionInfo(t *testing.T) {
	ragServer := SetupTestEnvironment()

	testParams := map[string]interface{}{
		"collection": "test_collection",
	}
	LogTestStart("TestGetCollectionInfo", testParams)

	req := &bs_rag.CollectionInfoRequest{
		CollectionName: "test_collection",
		UserId:         "test_user_001",
	}

	resp, err := ragServer.GetCollectionInfo(context.Background(), req)

	// 验证响应结构和内容
	if err != nil {
		LogTestResult("TestGetCollectionInfo", false, "Get collection info failed: "+err.Error())
		assert.Error(t, err)
	} else {
		LogTestResult("TestGetCollectionInfo", true, "Get collection info completed successfully")
		assert.NotNil(t, resp)

		// 如果集合存在，验证集合信息
		if resp.Exists {
			assert.NotEmpty(t, resp.CollectionName)
			assert.Greater(t, resp.VectorDimension, int32(0))
			LogTestResult("TestGetCollectionInfo", true,
				fmt.Sprintf("Collection exists - name: %s, dimension: %d", resp.CollectionName, resp.VectorDimension))
		} else {
			LogTestResult("TestGetCollectionInfo", true, "Collection does not exist, which is also valid")
		}

		// 验证错误信息为空
		assert.Empty(t, resp.ErrorMessage)
	}
}
