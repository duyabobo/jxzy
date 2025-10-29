package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LlmSceneModel = (*customLlmSceneModel)(nil)

type (
	// LlmSceneModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLlmSceneModel.
	LlmSceneModel interface {
		llmSceneModel
	}

	customLlmSceneModel struct {
		*defaultLlmSceneModel
	}
)

// NewLlmSceneModel returns a model for the database table.
func NewLlmSceneModel(conn sqlx.SqlConn) LlmSceneModel {
	return &customLlmSceneModel{
		defaultLlmSceneModel: newLlmSceneModel(conn),
	}
}
