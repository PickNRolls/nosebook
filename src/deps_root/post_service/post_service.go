package rootpostservice

import (
	"nosebook/src/services/posting"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *posting.Service {
	postService := posting.New(NewRepository(db))

	return postService
}
