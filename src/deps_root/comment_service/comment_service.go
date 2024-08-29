package rootcommentservice

import (
	rootpostservice "nosebook/src/deps_root/post_service"
	"nosebook/src/services/commenting"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *commenting.Service {
	commentService := commenting.New(newCommentRepository(db), rootpostservice.NewRepository(db))

	return commentService
}
