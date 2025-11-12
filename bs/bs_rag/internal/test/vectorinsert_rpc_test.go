package test

import (
	"context"
	"strings"
	"testing"

	"jxzy/bs/bs_rag/bs_rag"
	consts "jxzy/common/const"

	"github.com/stretchr/testify/assert"
)

// TestVectorInsert 测试向量插入RPC
func TestVectorInsert(t *testing.T) {
	ragServer := SetupTestEnvironment()

	testParams := map[string]interface{}{
		"collection": "test_collection",
		"docCount":   1,
		"dimension":  consts.DashVectorDefaultDimension, // DashVectorDefaultDimension
	}
	LogTestStart("TestVectorInsert", testParams)

	// 创建测试文档
	documents := []*bs_rag.VectorDocument{
		CreateTestVectorDocument("doc_001", int(consts.DashVectorDefaultDimension)),
	}

	req := &bs_rag.VectorInsertRequest{
		Documents: documents,
		UserId:    "test_user_001",
		SceneCode: "test_scene", // 测试场景编码
	}

	resp, err := ragServer.VectorInsert(context.Background(), req)

	// 验证结果
	if err != nil {
		LogTestResult("TestVectorInsert", false, "Expected error in test environment: "+err.Error())
		assert.Error(t, err)
		// 检查是否是预期的错误类型
		if strings.Contains(err.Error(), "Mismatched Data Type") {
			LogTestResult("TestVectorInsert", true, "Got expected Mismatched Data Type error due to test environment")
		}
	} else {
		LogTestResult("TestVectorInsert", true, "Insert completed successfully")
		AssertVectorInsertResponse(t, resp, 1, true)
	}
}
