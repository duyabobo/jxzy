package logic

import (
	"context"

	"jxzy/bll/bll_prompt/bll_prompt"
	"jxzy/bll/bll_prompt/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 占位方法，后续添加具体实现
func (l *PingLogic) Ping(in *bll_prompt.Empty) (*bll_prompt.Pong, error) {
	// todo: add your logic here and delete this line

	return &bll_prompt.Pong{}, nil
}
