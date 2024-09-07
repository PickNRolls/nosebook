package rootconvservice

import (
	"nosebook/src/application/services/conversation"
	"nosebook/src/application/services/socket"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB, hub *socket.Hub) *conversation.Service {
	service := conversation.New(
		newChatRepository(db),
		newNotifierRepository(hub, db),
		newUserRepository(db),
	)

	return service
}
