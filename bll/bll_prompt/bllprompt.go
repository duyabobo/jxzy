package main

import (
	"flag"

	"jxzy/bll/bll_prompt/bll_prompt"
	"jxzy/bll/bll_prompt/internal/config"
	"jxzy/bll/bll_prompt/internal/server"
	"jxzy/bll/bll_prompt/internal/svc"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/bllprompt.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化统一日志系统
	if err := logger.InitUnifiedLogger("bll-prompt"); err != nil {
		logx.Errorf("Failed to initialize logger: %v", err)
		return
	}

	ctx := svc.NewServiceContext(c)

	s, err := zrpc.NewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		bll_prompt.RegisterBllPromptServiceServer(grpcServer, server.NewBllPromptServiceServer(ctx))

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
