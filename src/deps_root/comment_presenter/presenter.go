package rootcommentpresenter

import (
	presentercomment "nosebook/src/presenters/comment"
	presentercommentlike "nosebook/src/presenters/comment_like"
	presenteruser "nosebook/src/presenters/user"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *presentercomment.Presenter {
	userPresenter := presenteruser.New(db)
	likePresenter := presentercommentlike.New(db, userPresenter)
	presenter := presentercomment.New(
		db,
		likePresenter,
		userPresenter,
	)

	return presenter
}
