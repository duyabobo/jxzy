package logic

import (
	"context"

	"jxzy/bll/bll_knowledge/bll_knowledge"
	"jxzy/bll/bll_knowledge/internal/svc"

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
func (l *PingLogic) Ping(in *bll_knowledge.Empty) (*bll_knowledge.Pong, error) {
	// todo: add your logic here and delete this line

	return &bll_knowledge.Pong{}, nil
}
