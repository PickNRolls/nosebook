package rootlikeservice

import (
	presenteruser "nosebook/src/application/presenters/user"
	"nosebook/src/application/services/like"
	"nosebook/src/infra/rabbitmq"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, rmqConn *rabbitmq.Connection) *like.Service {
	likeService := like.New(newRepository(db), newNotifier(rmqConn, db, presenteruser.New(db)))

	return likeService
}
