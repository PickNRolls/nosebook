package rootconvservice

import (
	"nosebook/src/application/services/conversation"
	"nosebook/src/infra/rabbitmq"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

func New(db *sqlx.DB, rmqConn *rabbitmq.Connection, tracer trace.Tracer) *conversation.Service {
	return conversation.New(
		newChatRepository(db),
		newNotifier(db, rmqConn),
		newUserRepository(db),
    tracer,
	)
}
