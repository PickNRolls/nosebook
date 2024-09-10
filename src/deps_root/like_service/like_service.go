package rootlikeservice

import (
	presenteruser "nosebook/src/application/presenters/user"
	"nosebook/src/application/services/like"
	"nosebook/src/infra/rabbitmq"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, rmqCh *rabbitmq.Channel) *like.Service {
	likeService := like.New(newRepository(db), newNotifier(rmqCh, db, presenteruser.New(db)))

	return likeService
}
