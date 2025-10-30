package test

import (
	"context"
	"testing"

	knowledgepb "jxzy/bll/bll_knowledge/bll_knowledge"
	"jxzy/bll/bll_knowledge/internal/config"
	"jxzy/bll/bll_knowledge/internal/server"
	"jxzy/bll/bll_knowledge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var cfg config.Config

func init() {
	conf.MustLoad("/Users/yb/GolandProjects/jxzy/bll/bll_knowledge/etc/bllknowledge.yaml", &cfg)
}

func TestAddVectorKnowledge_RPC(t *testing.T) {
	svcCtx := svc.NewServiceContext(cfg)
	s := server.NewBllKnowledgeServiceServer(svcCtx)

	req := &knowledgepb.AddVectorKnowledgeRequest{
		Summary: "测试总结",
		Content: "测试内容",
		UserId:  "test-user",
	}

	_, err := s.AddVectorKnowledge(context.Background(), req)
	if err != nil {
		t.Logf("AddVectorKnowledge error: %v", err)
	}
}

func TestDeleteVectorKnowledge_RPC(t *testing.T) {
	svcCtx := svc.NewServiceContext(cfg)
	s := server.NewBllKnowledgeServiceServer(svcCtx)

	req := &knowledgepb.DeleteVectorKnowledgeRequest{
		VectorId: "test-vector-id",
		UserId:   "test-user",
	}

	_, err := s.DeleteVectorKnowledge(context.Background(), req)
	if err != nil {
		t.Logf("DeleteVectorKnowledge error: %v", err)
	}
}
