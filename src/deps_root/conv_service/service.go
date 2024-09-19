package rootconvservice

import (
	"nosebook/src/application/services/conversation"
	"nosebook/src/infra/rabbitmq"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, rmqConn *rabbitmq.Connection) *conversation.Service {
	return conversation.New(
		newChatRepository(db),
		newNotifier(db, rmqConn),
		newUserRepository(db),
	)
}
