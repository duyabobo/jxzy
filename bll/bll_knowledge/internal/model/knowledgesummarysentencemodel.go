package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ KnowledgeSummarySentenceModel = (*customKnowledgeSummarySentenceModel)(nil)

type (
	// KnowledgeSummarySentenceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customKnowledgeSummarySentenceModel.
	KnowledgeSummarySentenceModel interface {
		knowledgeSummarySentenceModel
	}

	customKnowledgeSummarySentenceModel struct {
		*defaultKnowledgeSummarySentenceModel
	}
)

// NewKnowledgeSummarySentenceModel returns a model for the database table.
func NewKnowledgeSummarySentenceModel(conn sqlx.SqlConn, c cache.CacheConf) KnowledgeSummarySentenceModel {
	return &customKnowledgeSummarySentenceModel{
		defaultKnowledgeSummarySentenceModel: newKnowledgeSummarySentenceModel(conn, c),
	}
}
