package rootpostpresenter

import (
	presenterpost "nosebook/src/application/presenters/post"
	"nosebook/src/application/presenters/post_like"
	presenteruser "nosebook/src/application/presenters/user"
	rootcommentpresenter "nosebook/src/deps_root/comment_presenter"

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
		&permissions{},
	)

	return presenter
}
