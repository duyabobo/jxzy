package test

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"testing"

	contextpb "jxzy/bll/bll_context/bll_context"
	"jxzy/bll/bll_context/internal/config"
	"jxzy/bll/bll_context/internal/server"
	"jxzy/bll/bll_context/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"google.golang.org/grpc/metadata"
)

// MockStreamServer 模拟gRPC流服务器
type MockStreamServer struct {
	responses []*contextpb.StreamChatResponse
	ctx       context.Context
}

func NewMockStreamServer() *MockStreamServer {
	return &MockStreamServer{
		responses: make([]*contextpb.StreamChatResponse, 0),
		ctx:       context.Background(),
	}
}

func (m *MockStreamServer) Send(response *contextpb.StreamChatResponse) error {
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

var initOnce sync.Once
var globalConfig config.Config

func InitConfig() config.Config {
	initOnce.Do(func() {
		var configFile = flag.String("f", "/Users/yb/GolandProjects/jxzy/bll/bll_context/etc/bllcontext.yaml", "the config file")
		flag.Parse()
		conf.MustLoad(*configFile, &globalConfig)
	})
	return globalConfig
}

// TestStreamChat_RPC 测试RPC级别的StreamChat调用
func TestStreamChat_RPC(t *testing.T) {
	// 创建服务配置
	cfg := InitConfig()

	// 创建服务上下文
	svcCtx := svc.NewServiceContext(cfg)

	// 创建Context服务器
	contextServer := server.NewBllContextServiceServer(svcCtx)

	mockStream := NewMockStreamServer()

	req := &contextpb.ChatRequest{
		Message:   "hello",
		SceneCode: "chat_general",
		UserId:    "test-user",
		SessionId: "",
	}

	err := contextServer.StreamChat(req, mockStream)
	fmt.Println(err)
}
