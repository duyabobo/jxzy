package test

import (
	"flag"
	"sync"
	"testing"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/config"
	"jxzy/bs/bs_rag/internal/server"
	"jxzy/bs/bs_rag/internal/svc"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	configOnce     sync.Once
	configInstance config.Config
)

// InitConfig 初始化测试配置
func InitConfig() config.Config {
	configOnce.Do(func() {
		var configFile = flag.String("test-config", "/Users/yb/GolandProjects/jxzy/bs/bs_rag/etc/bsrag.yaml", "the config file for test")
		conf.MustLoad(*configFile, &configInstance)
	})
	return configInstance
}

// SetupTestEnvironment 设置测试环境
func SetupTestEnvironment() *server.BsRagServiceServer {
	// 初始化logx
	logx.SetUp(logx.LogConf{
		ServiceName: "bs-rag-test",
		Mode:        "console",
		Level:       "info",
		Encoding:    "json",
	})

	// 创建服务配置
	cfg := InitConfig()

	// 创建服务上下文
	svcCtx := svc.NewServiceContext(cfg)

	// 创建RAG服务器
	return server.NewBsRagServiceServer(svcCtx)
}

// CreateTestVectorDocument 创建测试向量文档
func CreateTestVectorDocument(id string, dimension int) *bs_rag.VectorDocument {
	// 生成标准化的向量数据 - 范围在[0,1]之间
	vector := make([]float32, dimension)
	for i := 0; i < dimension; i++ {
		vector[i] = float32(i+1) / float32(dimension) // 生成 [0.2, 0.4, 0.6, 0.8, 1.0] 类似的值
	}

	return &bs_rag.VectorDocument{
		Id:      id,
		Vector:  vector,
		Content: "Test content for document " + id,
		// 简化metadata，只保留必要的字符串字段
		Metadata: map[string]string{
			"source": "test",
		},
	}
}

// CreateTestQueryVector 创建测试查询向量
func CreateTestQueryVector(dimension int) []float32 {
	// 生成标准化的查询向量 - 范围在[0,1]之间
	vector := make([]float32, dimension)
	for i := 0; i < dimension; i++ {
		vector[i] = float32(i+1) / float32(dimension) // 与文档向量保持一致的生成方式
	}
	return vector
}

// AssertVectorSearchResponse 验证向量搜索响应
func AssertVectorSearchResponse(t *testing.T, resp *bs_rag.VectorSearchResponse, expectedCount int32, expectNoError bool) {
	assert.NotNil(t, resp)

	if expectNoError {
		assert.Equal(t, expectedCount, resp.TotalCount)
		assert.Len(t, resp.Results, int(expectedCount))

		// 验证结果的基本结构
		for _, result := range resp.Results {
			assert.NotEmpty(t, result.Id)
			assert.GreaterOrEqual(t, result.Score, float32(0))
		}
	}
}

// AssertVectorInsertResponse 验证向量插入响应
func AssertVectorInsertResponse(t *testing.T, resp *bs_rag.VectorInsertResponse, expectedCount int32, expectNoError bool) {
	assert.NotNil(t, resp)

	if expectNoError {
		assert.Equal(t, expectedCount, resp.InsertedCount)
		assert.Len(t, resp.InsertedIds, int(expectedCount))
		assert.Empty(t, resp.ErrorMessage)
	} else {
		assert.NotEmpty(t, resp.ErrorMessage)
	}
}

// AssertVectorDeleteResponse 验证向量删除响应
func AssertVectorDeleteResponse(t *testing.T, resp *bs_rag.VectorDeleteResponse, expectedCount int32, expectNoError bool) {
	assert.NotNil(t, resp)

	if expectNoError {
		assert.Equal(t, expectedCount, resp.DeletedCount)
		assert.Len(t, resp.DeletedIds, int(expectedCount))
		assert.Empty(t, resp.ErrorMessage)
	} else {
		assert.NotEmpty(t, resp.ErrorMessage)
	}
}

// AssertCollectionInfoResponse 验证集合信息响应
func AssertCollectionInfoResponse(t *testing.T, resp *bs_rag.CollectionInfoResponse, expectedExists bool, expectNoError bool) {
	assert.NotNil(t, resp)

	if expectNoError {
		assert.Equal(t, expectedExists, resp.Exists)
		if expectedExists {
			assert.NotEmpty(t, resp.CollectionName)
			assert.GreaterOrEqual(t, resp.VectorDimension, int32(0))
		}
		assert.Empty(t, resp.ErrorMessage)
	} else {
		assert.NotEmpty(t, resp.ErrorMessage)
	}
}

// LogTestStart 记录测试开始
func LogTestStart(testName string, params map[string]interface{}) {
	logx.Infof("Starting test: %s with params: %+v", testName, params)
}

// LogTestResult 记录测试结果
func LogTestResult(testName string, success bool, message string) {
	if success {
		logx.Infof("Test %s passed: %s", testName, message)
	} else {
		logx.Errorf("Test %s failed: %s", testName, message)
	}
}
