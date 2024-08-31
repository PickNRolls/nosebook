package presenterpostlike

import (
	"nosebook/src/errors"
	presenterdto "nosebook/src/presenters/dto"
	presenterlike "nosebook/src/presenters/like"
	"nosebook/src/services/auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db            *sqlx.DB
	likePresenter *presenterlike.Presenter
}

func New(db *sqlx.DB, userPresenter presenterlike.UserPresenter) *Presenter {
	return &Presenter{
		db:            db,
		likePresenter: presenterlike.New(db, userPresenter, &postResource{}),
	}
}

func (this *Presenter) FindByPostIds(
	ids uuid.UUIDs,
	auth *auth.Auth,
) (map[uuid.UUID]*presenterdto.Likes, *errors.Error) {
	return this.likePresenter.FindByResourceIds(ids, auth)
}
