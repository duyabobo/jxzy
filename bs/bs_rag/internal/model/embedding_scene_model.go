package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EmbeddingSceneModel = (*customEmbeddingSceneModel)(nil)

type (
	// EmbeddingSceneModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEmbeddingSceneModel.
	EmbeddingSceneModel interface {
		embeddingSceneModel
	}

	customEmbeddingSceneModel struct {
		*defaultEmbeddingSceneModel
	}
)

// NewEmbeddingSceneModel returns a model for the database table.
func NewEmbeddingSceneModel(conn sqlx.SqlConn) EmbeddingSceneModel {
	return &customEmbeddingSceneModel{
		defaultEmbeddingSceneModel: newEmbeddingSceneModel(conn),
	}
}
