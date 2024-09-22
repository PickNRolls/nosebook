package presentercommentlike

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	presenterlike "nosebook/src/application/presenters/like"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type Presenter struct {
	db            *sqlx.DB
	likePresenter *presenterlike.Presenter
}

func New(db *sqlx.DB, userPresenter presenterlike.UserPresenter, tracer trace.Tracer) *Presenter {
	return &Presenter{
		db:            db,
		likePresenter: presenterlike.New(db, userPresenter).WithTracer(tracer),
	}
}

func (this *Presenter) FindByCommentIds(
  parent context.Context,
	ids uuid.UUIDs,
	auth *auth.Auth,
) (map[uuid.UUID]*presenterdto.Likes, *errors.Error) {
	return this.likePresenter.FindByResourceIds(parent, &commentResource{}, ids, auth)
}
