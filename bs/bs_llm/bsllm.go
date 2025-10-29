package main

import (
	"flag"

	"jxzy/bs/bs_llm/bs_llm"
	"jxzy/bs/bs_llm/internal/config"
	"jxzy/bs/bs_llm/internal/server"
	"jxzy/bs/bs_llm/internal/svc"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/bsllm.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化统一日志系统
	if err := logger.InitUnifiedLogger("bs-llm"); err != nil {
		logx.Errorf("Failed to initialize logger: %v", err)
		return
	}

	ctx := svc.NewServiceContext(c)

	s, err := zrpc.NewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		bs_llm.RegisterBsLlmServiceServer(grpcServer, server.NewBsLlmServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	if err != nil {
		logx.Errorf("Failed to create rpc server: %v", err)
		return
	}
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}
