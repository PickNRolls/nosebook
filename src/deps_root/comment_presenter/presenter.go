package rootcommentpresenter

import (
	presentercomment "nosebook/src/application/presenters/comment"
	presentercommentlike "nosebook/src/application/presenters/comment_like"
	presenteruser "nosebook/src/application/presenters/user"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *presentercomment.Presenter {
	userPresenter := presenteruser.New(db)
	likePresenter := presentercommentlike.New(db, userPresenter)

	presenter := presentercomment.New(
		db,
		likePresenter,
		userPresenter,
		&permissions{},
	)

	return presenter
}
