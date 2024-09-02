package rootcommentservice

import (
	"github.com/jmoiron/sqlx"
	"nosebook/src/application/services/commenting"
	rootpostservice "nosebook/src/deps_root/post_service"
)

func New(db *sqlx.DB) *commenting.Service {
	commentService := commenting.New(
		newCommentRepository(db, newPermissions(db)),
		rootpostservice.NewRepository(db),
	)

	return commentService
}
