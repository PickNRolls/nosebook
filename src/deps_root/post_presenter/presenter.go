package rootpostpresenter

import (
	presentercomment "nosebook/src/presenters/comment"
	presenterlike "nosebook/src/presenters/like"
	presenterpost "nosebook/src/presenters/post"
	presenteruser "nosebook/src/presenters/user"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *presenterpost.Presenter {
	userPresenter := presenteruser.New(db)
	presenter := presenterpost.New(
		db,
		userPresenter,
		presentercomment.New(db),
		presenterlike.New(db, userPresenter),
	)

	return presenter
}
