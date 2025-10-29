package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ KnowledgeSegmentModel = (*customKnowledgeSegmentModel)(nil)

type (
	// KnowledgeSegmentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customKnowledgeSegmentModel.
	KnowledgeSegmentModel interface {
		knowledgeSegmentModel
	}

	customKnowledgeSegmentModel struct {
		*defaultKnowledgeSegmentModel
	}
)

// NewKnowledgeSegmentModel returns a model for the database table.
func NewKnowledgeSegmentModel(conn sqlx.SqlConn) KnowledgeSegmentModel {
	return &customKnowledgeSegmentModel{
		defaultKnowledgeSegmentModel: newKnowledgeSegmentModel(conn),
	}
}
