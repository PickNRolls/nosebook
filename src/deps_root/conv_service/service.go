package rootconvservice

import (
	presentermessage "nosebook/src/application/presenters/message"
	"nosebook/src/application/services/conversation"
	"nosebook/src/infra/rabbitmq"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

func New(db *sqlx.DB, rmqConn *rabbitmq.Connection, presenter *presentermessage.Presenter, tracer trace.Tracer) *conversation.Service {
  chatRepository := newChatRepository(db)
  go chatRepository.Run()
  
	return conversation.New(
		chatRepository,
		newNotifier(rmqConn, presenter),
		newUserRepository(db),
    tracer,
	)
}
