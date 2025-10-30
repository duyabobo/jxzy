package test

import (
    "context"
    "testing"

    promptpb "jxzy/bll/bll_prompt/bll_prompt"
    "jxzy/bll/bll_prompt/internal/config"
    "jxzy/bll/bll_prompt/internal/server"
    "jxzy/bll/bll_prompt/internal/svc"

    "github.com/zeromicro/go-zero/core/conf"
)

var cfg config.Config

func init() {
    conf.MustLoad("/Users/yb/GolandProjects/jxzy/bll/bll_prompt/etc/bllprompt.yaml", &cfg)
}

func TestPing_RPC(t *testing.T) {
    svcCtx := svc.NewServiceContext(cfg)
    s := server.NewBllPromptServiceServer(svcCtx)

    // 如果未来去掉 Ping，可将此测试替换成具体方法
    _, err := s.Ping(context.Background(), &promptpb.Empty{})
    if err != nil {
        t.Logf("Ping error: %v", err)
    }
}


