package rootlikeservice

import (
	"nosebook/src/application/services/like"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *like.Service {
	likeService := like.New(newRepository(db))

	return likeService
}
