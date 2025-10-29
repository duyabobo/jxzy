package main

import (
	"flag"

	"jxzy/apis/api_chat/internal/config"
	"jxzy/apis/api_chat/internal/handler"
	"jxzy/apis/api_chat/internal/svc"
	"jxzy/common/logger"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/chat-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化统一日志系统
	if err := logger.InitUnifiedLogger("chat-api"); err != nil {
		logx.Errorf("Failed to initialize logger: %v", err)
		return
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
	server.Start()
}
