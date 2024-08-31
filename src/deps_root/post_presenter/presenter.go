package rootpostpresenter

import (
	rootcommentpresenter "nosebook/src/deps_root/comment_presenter"
	presenterpost "nosebook/src/presenters/post"
	"nosebook/src/presenters/post_like"
	presenteruser "nosebook/src/presenters/user"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *presenterpost.Presenter {
	userPresenter := presenteruser.New(db)
	likePresenter := presenterpostlike.New(db, userPresenter)
	presenter := presenterpost.New(
		db,
		userPresenter,
		rootcommentpresenter.New(db),
		likePresenter,
	)

	return presenter
}
