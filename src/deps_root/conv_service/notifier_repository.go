package rootconvservice

import (
	presentermessage "nosebook/src/application/presenters/message"
	presenteruser "nosebook/src/application/presenters/user"
	"nosebook/src/application/services/conversation"
	"nosebook/src/application/services/socket"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type notifierRepository struct {
	hub       *socket.Hub
	db        *sqlx.DB
	presenter *presentermessage.Presenter
}

func newNotifierRepository(hub *socket.Hub, db *sqlx.DB) *notifierRepository {
	return &notifierRepository{
		hub:       hub,
		presenter: presentermessage.New(db, presenteruser.New(db)),
	}
}

func (this *notifierRepository) FindByRecipientId(id uuid.UUID) conversation.Notifier {
	client := this.hub.UserClient(id)
	if client == nil {
		return nil
	}

	return &socketNotifier{
		client:    client,
		presenter: this.presenter,
	}
}
