package main

import (
	"flag"
	"fmt"

	"jxzy/bs/bs_rag/bs_rag"
	"jxzy/bs/bs_rag/internal/config"
	"jxzy/bs/bs_rag/internal/server"
	"jxzy/bs/bs_rag/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/bsrag.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewBsRagServiceServer(ctx)

	s, err := zrpc.NewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		bs_rag.RegisterBsRagServiceServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	if err != nil {
		fmt.Printf("Failed to create rpc server: %v\n", err)
		return
	}
	defer s.Stop()

	fmt.Printf("Starting bs_rag rpc server at %s...\n", c.ListenOn)
	s.Start()
}
