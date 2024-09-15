package presenterpostlike

import (
	presenterdto "nosebook/src/application/presenters/dto"
	presenterlike "nosebook/src/application/presenters/like"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

type Presenter struct {
	db            *sqlx.DB
	likePresenter *presenterlike.Presenter
}

func New(db *sqlx.DB, userPresenter presenterlike.UserPresenter) *Presenter {
	return &Presenter{
		db:            db,
		likePresenter: presenterlike.New(db, userPresenter),
	}
}

func (this *Presenter) WithTracer(tracer trace.Tracer) *Presenter {
  this.likePresenter.WithTracer(tracer)

  return this
}

func (this *Presenter) FindByPostIds(
  parent context.Context,
	ids uuid.UUIDs,
	auth *auth.Auth,
) (map[uuid.UUID]*presenterdto.Likes, *errors.Error) {
	return this.likePresenter.FindByResourceIds(parent, &postResource{}, ids, auth)
}
