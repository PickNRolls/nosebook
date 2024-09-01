package rootcommentservice

import (
	"nosebook/src/application/services/commenting"
	rootpostservice "nosebook/src/deps_root/post_service"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *commenting.Service {
	commentService := commenting.New(newCommentRepository(db), rootpostservice.NewRepository(db))

	return commentService
}
