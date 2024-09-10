package rootconvservice

import (
	"nosebook/src/application/services/conversation"
	"nosebook/src/infra/rabbitmq"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, rmqCh *rabbitmq.Channel) *conversation.Service {
	service := conversation.New(
		newChatRepository(db),
		newNotifier(db, rmqCh),
		newUserRepository(db),
	)

	return service
}
