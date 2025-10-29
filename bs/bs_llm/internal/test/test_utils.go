package test

import (
	"context"
	"flag"
	"sync"

	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/config"

	"github.com/zeromicro/go-zero/core/conf"
	"google.golang.org/grpc/metadata"
)

// MockStreamServer 模拟gRPC流服务器
type MockStreamServer struct {
	responses []*bs_llm.StreamLLMResponse
	ctx       context.Context
}

func NewMockStreamServer() *MockStreamServer {
	return &MockStreamServer{
		responses: make([]*bs_llm.StreamLLMResponse, 0),
		ctx:       context.Background(),
	}
}

func (m *MockStreamServer) Send(response *bs_llm.StreamLLMResponse) error {
	m.responses = append(m.responses, response)
	return nil
}

func (m *MockStreamServer) Context() context.Context {
	return m.ctx
}

func (m *MockStreamServer) SendMsg(msg interface{}) error {
	return nil
}

func (m *MockStreamServer) RecvMsg(msg interface{}) error {
	return nil
}

func (m *MockStreamServer) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockStreamServer) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockStreamServer) SetTrailer(md metadata.MD) {
}

var testConfig config.Config
var configOnce sync.Once

func GetTestConfig() config.Config {
	configOnce.Do(func() {
		var configFile = flag.String("f", "/Users/yb/GolandProjects/jxzy/bs/bs_llm/etc/bsllm.yaml", "the config file")
		flag.Parse()
		conf.MustLoad(*configFile, &testConfig)
	})
	return testConfig
}
