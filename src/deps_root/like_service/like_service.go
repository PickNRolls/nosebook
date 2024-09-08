package rootlikeservice

import (
	"nosebook/src/application/services/like"
	"nosebook/src/application/services/socket"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, hub *socket.Hub) *like.Service {
	likeService := like.New(newRepository(db), newNotifierRepository(hub, db))

	return likeService
}
