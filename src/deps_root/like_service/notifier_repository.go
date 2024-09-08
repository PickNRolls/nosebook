package rootlikeservice

import (
	presenteruser "nosebook/src/application/presenters/user"
	"nosebook/src/application/services/like"
	"nosebook/src/application/services/socket"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type notifierRepository struct {
	hub           *socket.Hub
	db            *sqlx.DB
	userPresenter *presenteruser.Presenter
}

func newNotifierRepository(hub *socket.Hub, db *sqlx.DB) *notifierRepository {
	return &notifierRepository{
		hub:           hub,
		db:            db,
		userPresenter: presenteruser.New(db),
	}
}

func (this *notifierRepository) FindByUserId(id uuid.UUID) like.Notifier {
	if this.hub.UserClients(id) == nil {
		return nil
	}

	return &socketNotifier{
		db:            this.db,
		hub:           this.hub,
		userId:        id,
		userPresenter: this.userPresenter,
	}
}
