package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LlmCompletionModel = (*customLlmCompletionModel)(nil)

type (
	// LlmCompletionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLlmCompletionModel.
	LlmCompletionModel interface {
		llmCompletionModel
	}

	customLlmCompletionModel struct {
		*defaultLlmCompletionModel
	}
)

// NewLlmCompletionModel returns a model for the database table.
func NewLlmCompletionModel(conn sqlx.SqlConn) LlmCompletionModel {
	return &customLlmCompletionModel{
		defaultLlmCompletionModel: newLlmCompletionModel(conn),
	}
}
