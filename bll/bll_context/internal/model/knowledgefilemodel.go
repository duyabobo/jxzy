package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ KnowledgeFileModel = (*customKnowledgeFileModel)(nil)

type (
	// KnowledgeFileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customKnowledgeFileModel.
	KnowledgeFileModel interface {
		knowledgeFileModel
	}

	customKnowledgeFileModel struct {
		*defaultKnowledgeFileModel
	}
)

// NewKnowledgeFileModel returns a model for the database table.
func NewKnowledgeFileModel(conn sqlx.SqlConn) KnowledgeFileModel {
	return &customKnowledgeFileModel{
		defaultKnowledgeFileModel: newKnowledgeFileModel(conn),
	}
}
