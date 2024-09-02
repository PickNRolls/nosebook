package rootfriendshippresenter

import (
	presenterfriendship "nosebook/src/application/presenters/friendship"
	presenteruser "nosebook/src/application/presenters/user"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *presenterfriendship.Presenter {
	userPresenter := presenteruser.New(db)
	presenter := presenterfriendship.New(
		db,
		userPresenter,
	)

	return presenter
}
