package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ChatSessionQasModel = (*customChatSessionQasModel)(nil)

type (
	// ChatSessionQasModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatSessionQasModel.
	ChatSessionQasModel interface {
		chatSessionQasModel
	}

	customChatSessionQasModel struct {
		*defaultChatSessionQasModel
	}
)

// NewChatSessionQasModel returns a model for the database table.
func NewChatSessionQasModel(conn sqlx.SqlConn) ChatSessionQasModel {
	return &customChatSessionQasModel{
		defaultChatSessionQasModel: newChatSessionQasModel(conn),
	}
}
